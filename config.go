// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (C) 2015-2017 The Lightning Network Developers

package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil"
	flags "github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd/htlcswitch/hodl"
	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/tor"
)

const (
	defaultConfigFilename      = "lnd.conf"
	defaultDataDirname         = "data"
	defaultChainSubDirname     = "chain"
	defaultGraphSubDirname     = "graph"
	defaultTLSCertFilename     = "tls.cert"
	defaultTLSKeyFilename      = "tls.key"
	defaultAdminMacFilename    = "admin.macaroon"
	defaultReadMacFilename     = "readonly.macaroon"
	defaultInvoiceMacFilename  = "invoice.macaroon"
	defaultLogLevel            = "info"
	defaultLogDirname          = "logs"
	defaultLogFilename         = "lnd.log"
	defaultRPCPort             = 10009
	defaultRESTPort            = 8080
	defaultPeerPort            = 9735
	defaultRPCHost             = "localhost"
	defaultMaxPendingChannels  = 1
	defaultNoSeedBackup        = false
	defaultTrickleDelay        = 30 * 1000
	defaultInactiveChanTimeout = 20 * time.Minute
	defaultMaxLogFiles         = 3
	defaultMaxLogFileSize      = 10

	defaultTorSOCKSPort            = 9050
	defaultTorDNSHost              = "soa.nodes.lightning.directory"
	defaultTorDNSPort              = 53
	defaultTorControlPort          = 9051
	defaultTorV2PrivateKeyFilename = "v2_onion_private_key"
	defaultTorV3PrivateKeyFilename = "v3_onion_private_key"

	defaultBroadcastDelta = 10

	// minTimeLockDelta is the minimum timelock we require for incoming
	// HTLCs on our channels.
	minTimeLockDelta = 4

	defaultAlias = ""
	defaultColor = "#3399FF"
)

var (
	defaultLndDir     = btcutil.AppDataDir("lnd", false)
	defaultConfigFile = filepath.Join(defaultLndDir, defaultConfigFilename)
	defaultDataDir    = filepath.Join(defaultLndDir, defaultDataDirname)
	defaultLogDir     = filepath.Join(defaultLndDir, defaultLogDirname)

	defaultTLSCertPath = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultTLSKeyPath  = filepath.Join(defaultLndDir, defaultTLSKeyFilename)

	defaultBtcdDir         = btcutil.AppDataDir("btcd", false)
	defaultBtcdRPCCertFile = filepath.Join(defaultBtcdDir, "rpc.cert")

	defaultLtcdDir         = btcutil.AppDataDir("ltcd", false)
	defaultLtcdRPCCertFile = filepath.Join(defaultLtcdDir, "rpc.cert")

	defaultBitcoindDir  = btcutil.AppDataDir("bitcoin", false)
	defaultLitecoindDir = btcutil.AppDataDir("litecoin", false)

	defaultTorSOCKS   = net.JoinHostPort("localhost", strconv.Itoa(defaultTorSOCKSPort))
	defaultTorDNS     = net.JoinHostPort(defaultTorDNSHost, strconv.Itoa(defaultTorDNSPort))
	defaultTorControl = net.JoinHostPort("localhost", strconv.Itoa(defaultTorControlPort))
)

type chainConfig struct {
	Active   bool   `long:"active" description:"If the chain should be active or not."`
	ChainDir string `long:"chaindir" description:"The directory to store the chain's data within."`

	Node string `long:"node" description:"The blockchain interface to use." choice:"btcd" choice:"bitcoind" choice:"neutrino" choice:"ltcd" choice:"litecoind"`

	MainNet  bool `long:"mainnet" description:"Use the main network"`
	TestNet3 bool `long:"testnet" description:"Use the test network"`
	SimNet   bool `long:"simnet" description:"Use the simulation test network"`
	RegTest  bool `long:"regtest" description:"Use the regression test network"`

	DefaultNumChanConfs int                 `long:"defaultchanconfs" description:"The default number of confirmations a channel must have before it's considered open. If this is not set, we will scale the value according to the channel size."`
	DefaultRemoteDelay  int                 `long:"defaultremotedelay" description:"The default number of blocks we will require our channel counterparty to wait before accessing its funds in case of unilateral close. If this is not set, we will scale the value according to the channel size."`
	MinHTLC             lnwire.MilliSatoshi `long:"minhtlc" description:"The smallest HTLC we are willing to forward on our channels, in millisatoshi"`
	BaseFee             lnwire.MilliSatoshi `long:"basefee" description:"The base fee in millisatoshi we will charge for forwarding payments on our channels"`
	FeeRate             lnwire.MilliSatoshi `long:"feerate" description:"The fee rate used when forwarding payments on our channels. The total fee charged is basefee + (amount * feerate / 1000000), where amount is the forwarded amount."`
	TimeLockDelta       uint32              `long:"timelockdelta" description:"The CLTV delta we will subtract from a forwarded HTLC's timelock value"`
}

type neutrinoConfig struct {
	AddPeers     []string      `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
	ConnectPeers []string      `long:"connect" description:"Connect only to the specified peers at startup"`
	MaxPeers     int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	BanDuration  time.Duration `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold uint32        `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`
}

type btcdConfig struct {
	Dir        string `long:"dir" description:"The base directory that contains the node's data, logs, configuration file, etc."`
	RPCHost    string `long:"rpchost" description:"The daemon's rpc listening address. If a port is omitted, then the default port for the selected chain parameters will be used."`
	RPCUser    string `long:"rpcuser" description:"Username for RPC connections"`
	RPCPass    string `long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	RPCCert    string `long:"rpccert" description:"File containing the daemon's certificate file"`
	RawRPCCert string `long:"rawrpccert" description:"The raw bytes of the daemon's PEM-encoded certificate chain which will be used to authenticate the RPC connection."`
}

type bitcoindConfig struct {
	Dir            string `long:"dir" description:"The base directory that contains the node's data, logs, configuration file, etc."`
	RPCHost        string `long:"rpchost" description:"The daemon's rpc listening address. If a port is omitted, then the default port for the selected chain parameters will be used."`
	RPCUser        string `long:"rpcuser" description:"Username for RPC connections"`
	RPCPass        string `long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	ZMQPubRawBlock string `long:"zmqpubrawblock" description:"The address listening for ZMQ connections to deliver raw block notifications"`
	ZMQPubRawTx    string `long:"zmqpubrawtx" description:"The address listening for ZMQ connections to deliver raw transaction notifications"`
}

type autoPilotConfig struct {
	Active         bool    `long:"active" description:"If the autopilot agent should be active or not."`
	MaxChannels    int     `long:"maxchannels" description:"The maximum number of channels that should be created"`
	Allocation     float64 `long:"allocation" description:"The percentage of total funds that should be committed to automatic channel establishment"`
	MinChannelSize int64   `long:"minchansize" description:"The smallest channel that the autopilot agent should create"`
	MaxChannelSize int64   `long:"maxchansize" description:"The largest channel that the autopilot agent should create"`
	Private        bool    `long:"private" description:"Whether the channels created by the autopilot agent should be private or not. Private channels won't be announced to the network."`
	MinConfs       int32   `long:"minconfs" description:"The minimum number of confirmations each of your inputs in funding transactions created by the autopilot agent must have."`
}

type torConfig struct {
	Active          bool   `long:"active" description:"Allow outbound and inbound connections to be routed through Tor"`
	SOCKS           string `long:"socks" description:"The host:port that Tor's exposed SOCKS5 proxy is listening on"`
	DNS             string `long:"dns" description:"The DNS server as host:port that Tor will use for SRV queries - NOTE must have TCP resolution enabled"`
	StreamIsolation bool   `long:"streamisolation" description:"Enable Tor stream isolation by randomizing user credentials for each connection."`
	Control         string `long:"control" description:"The host:port that Tor is listening on for Tor control connections"`
	V2              bool   `long:"v2" description:"Automatically set up a v2 onion service to listen for inbound connections"`
	V3              bool   `long:"v3" description:"Automatically set up a v3 onion service to listen for inbound connections"`
	PrivateKeyPath  string `long:"privatekeypath" description:"The path to the private key of the onion service being created"`
}

type zapConnectConfig struct {
	LocalIp   bool `short:"i" long:"localip" description:"Include local ip in QRCode."`
	Localhost bool `short:"l" long:"localhost" description:"Use 127.0.0.1 for ip."`
	Json      bool `short:"j" long:"json" description:"Generate json instead of a QRCode."`
	Image     bool `short:"o" long:"image" description:"Output QRCode to file."`
	Invoice   bool `long:"invoice" description:"use invoice macaroon"`
	Readonly  bool `long:"readonly" description:"use readonly macaroon"`
}

// config defines the configuration options for lnd.
//
// See loadConfig for further details regarding the configuration
// loading+parsing process.
type config struct {
	ZapConnect *zapConnectConfig `group:"ZapConnect"`

	ShowVersion bool `short:"V" long:"version" description:"Display version information and exit"`

	LndDir         string `long:"lnddir" description:"The base directory that contains lnd's data, logs, configuration file, etc."`
	ConfigFile     string `long:"C" long:"configfile" description:"Path to configuration file"`
	DataDir        string `short:"b" long:"datadir" description:"The directory to store lnd's data within"`
	TLSCertPath    string `long:"tlscertpath" description:"Path to write the TLS certificate for lnd's RPC and REST services"`
	TLSKeyPath     string `long:"tlskeypath" description:"Path to write the TLS private key for lnd's RPC and REST services"`
	TLSExtraIP     string `long:"tlsextraip" description:"Adds an extra ip to the generated certificate"`
	TLSExtraDomain string `long:"tlsextradomain" description:"Adds an extra domain to the generated certificate"`
	NoMacaroons    bool   `long:"no-macaroons" description:"Disable macaroon authentication"`
	AdminMacPath   string `long:"adminmacaroonpath" description:"Path to write the admin macaroon for lnd's RPC and REST services if it doesn't exist"`
	ReadMacPath    string `long:"readonlymacaroonpath" description:"Path to write the read-only macaroon for lnd's RPC and REST services if it doesn't exist"`
	InvoiceMacPath string `long:"invoicemacaroonpath" description:"Path to the invoice-only macaroon for lnd's RPC and REST services if it doesn't exist"`
	LogDir         string `long:"logdir" description:"Directory to log output."`
	MaxLogFiles    int    `long:"maxlogfiles" description:"Maximum logfiles to keep (0 for no rotation)"`
	MaxLogFileSize int    `long:"maxlogfilesize" description:"Maximum logfile size in MB"`

	// We'll parse these 'raw' string arguments into real net.Addrs in the
	// loadConfig function. We need to expose the 'raw' strings so the
	// command line library can access them.
	// Only the parsed net.Addrs should be used!
	RawRPCListeners  []string `long:"rpclisten" description:"Add an interface/port/socket to listen for RPC connections"`
	RawRESTListeners []string `long:"restlisten" description:"Add an interface/port/socket to listen for REST connections"`
	RawListeners     []string `long:"listen" description:"Add an interface/port to listen for peer connections"`
	RawExternalIPs   []string `long:"externalip" description:"Add an ip:port to the list of local addresses we claim to listen on to peers. If a port is not specified, the default (9735) will be used regardless of other parameters"`
	RPCListeners     []net.Addr
	RESTListeners    []net.Addr
	Listeners        []net.Addr
	ExternalIPs      []net.Addr
	DisableListen    bool `long:"nolisten" description:"Disable listening for incoming peer connections"`
	NAT              bool `long:"nat" description:"Toggle NAT traversal support (using either UPnP or NAT-PMP) to automatically advertise your external IP address to the network -- NOTE this does not support devices behind multiple NATs"`

	DebugLevel string `short:"d" long:"debuglevel" description:"Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available subsystems"`

	CPUProfile string `long:"cpuprofile" description:"Write CPU profile to the specified file"`

	Profile string `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65535"`

	DebugHTLC          bool `long:"debughtlc" description:"Activate the debug htlc mode. With the debug HTLC mode, all payments sent use a pre-determined R-Hash. Additionally, all HTLCs sent to a node with the debug HTLC R-Hash are immediately settled in the next available state transition."`
	UnsafeDisconnect   bool `long:"unsafe-disconnect" description:"Allows the rpcserver to intentionally disconnect from peers with open channels. USED FOR TESTING ONLY."`
	UnsafeReplay       bool `long:"unsafe-replay" description:"Causes a link to replay the adds on its commitment txn after starting up, this enables testing of the sphinx replay logic."`
	MaxPendingChannels int  `long:"maxpendingchannels" description:"The maximum number of incoming pending channels permitted per peer."`

	Bitcoin      *chainConfig    `group:"Bitcoin" namespace:"bitcoin"`
	BtcdMode     *btcdConfig     `group:"btcd" namespace:"btcd"`
	BitcoindMode *bitcoindConfig `group:"bitcoind" namespace:"bitcoind"`
	NeutrinoMode *neutrinoConfig `group:"neutrino" namespace:"neutrino"`

	Litecoin      *chainConfig    `group:"Litecoin" namespace:"litecoin"`
	LtcdMode      *btcdConfig     `group:"ltcd" namespace:"ltcd"`
	LitecoindMode *bitcoindConfig `group:"litecoind" namespace:"litecoind"`

	Autopilot *autoPilotConfig `group:"Autopilot" namespace:"autopilot"`

	Tor *torConfig `group:"Tor" namespace:"tor"`

	Hodl *hodl.Config `group:"hodl" namespace:"hodl"`

	NoNetBootstrap bool `long:"nobootstrap" description:"If true, then automatic network bootstrapping will not be attempted."`

	NoSeedBackup bool `long:"noseedbackup" description:"If true, NO SEED WILL BE EXPOSED AND THE WALLET WILL BE ENCRYPTED USING THE DEFAULT PASSPHRASE -- EVER. THIS FLAG IS ONLY FOR TESTING AND IS BEING DEPRECATED."`

	TrickleDelay        int           `long:"trickledelay" description:"Time in milliseconds between each release of announcements to the network"`
	InactiveChanTimeout time.Duration `long:"inactivechantimeout" description:"If a channel has been inactive for the set time, send a ChannelUpdate disabling it."`

	Alias       string `long:"alias" description:"The node alias. Used as a moniker by peers and intelligence services"`
	Color       string `long:"color" description:"The color of the node in hex format (i.e. '#3399FF'). Used to customize node appearance in intelligence services"`
	MinChanSize int64  `long:"minchansize" description:"The smallest channel size (in satoshis) that we should accept. Incoming channels smaller than this will be rejected"`

	NoChanUpdates bool `long:"nochanupdates" description:"If specified, lnd will not request real-time channel updates from connected peers. This option should be used by routing nodes to save bandwidth."`

	RejectPush bool `long:"rejectpush" description:"If true, lnd will not accept channel opening requests with non-zero push amounts. This should prevent accidental pushes to merchant nodes."`

	net tor.Net

	// Routing *routing.Conf `group:"routing" namespace:"routing"`
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
func loadConfig() (*config, error) {
	defaultCfg := config{
		LndDir:         defaultLndDir,
		ConfigFile:     defaultConfigFile,
		DataDir:        defaultDataDir,
		DebugLevel:     defaultLogLevel,
		TLSCertPath:    defaultTLSCertPath,
		TLSKeyPath:     defaultTLSKeyPath,
		LogDir:         defaultLogDir,
		MaxLogFiles:    defaultMaxLogFiles,
		MaxLogFileSize: defaultMaxLogFileSize,
		Bitcoin: &chainConfig{
			MinHTLC:       0,
			BaseFee:       0,
			FeeRate:       0,
			TimeLockDelta: 0,
			Node:          "btcd",
		},
		BtcdMode: &btcdConfig{
			Dir:     defaultBtcdDir,
			RPCHost: defaultRPCHost,
			RPCCert: defaultBtcdRPCCertFile,
		},
		BitcoindMode: &bitcoindConfig{
			Dir:     defaultBitcoindDir,
			RPCHost: defaultRPCHost,
		},
		Litecoin: &chainConfig{
			MinHTLC:       0,
			BaseFee:       0,
			FeeRate:       0,
			TimeLockDelta: 0,
			Node:          "ltcd",
		},
		LtcdMode: &btcdConfig{
			Dir:     defaultLtcdDir,
			RPCHost: defaultRPCHost,
			RPCCert: defaultLtcdRPCCertFile,
		},
		LitecoindMode: &bitcoindConfig{
			Dir:     defaultLitecoindDir,
			RPCHost: defaultRPCHost,
		},
		MaxPendingChannels: defaultMaxPendingChannels,
		NoSeedBackup:       defaultNoSeedBackup,
		Autopilot: &autoPilotConfig{
			MaxChannels:    5,
			Allocation:     0.6,
			MinChannelSize: 0,
			MaxChannelSize: 0,
		},
		TrickleDelay:        defaultTrickleDelay,
		InactiveChanTimeout: defaultInactiveChanTimeout,
		Alias:               defaultAlias,
		Color:               defaultColor,
		MinChanSize:         0,
		Tor: &torConfig{
			SOCKS:   defaultTorSOCKS,
			DNS:     defaultTorDNS,
			Control: defaultTorControl,
		},
		net: &tor.ClearNet{},
	}

	// Pre-parse the command line options to pick up an alternative config
	// file.
	preCfg := defaultCfg
	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	// If the provided lnd directory is not the default, we'll modify the
	// path to all of the files and directories that will live within it.
	lndDir := cleanAndExpandPath(preCfg.LndDir)
	configFilePath := cleanAndExpandPath(preCfg.ConfigFile)
	if lndDir != defaultLndDir {
		// If the config file path has not been modified by the user,
		// then we'll use the default config file path. However, if the
		// user has modified their lnddir, then we should assume they
		// intend to use the config file within it.
		if configFilePath == defaultConfigFile {
			preCfg.ConfigFile = filepath.Join(lndDir, defaultConfigFilename)
		}
		preCfg.DataDir = filepath.Join(lndDir, defaultDataDirname)
		preCfg.TLSCertPath = filepath.Join(lndDir, defaultTLSCertFilename)
	}

	// Next, load any additional configuration options from the file.
	var configFileError error
	cfg := preCfg
	if err := flags.IniParse(cfg.ConfigFile, &cfg); err != nil {
		configFileError = err
	}

	// Finally, parse the remaining command line options again to ensure
	// they take precedence.
	if _, err := flags.Parse(&cfg); err != nil {
		return nil, err
	}

	primaryChain := "bitcoin"
	networkName := "mainnet"

	switch {
	case cfg.Litecoin.Active:
		primaryChain = "litecoin"
		networkName = "mainnet"
	case cfg.Bitcoin.Active:
		if cfg.Bitcoin.MainNet {
			networkName = "mainnet"
		}
		if cfg.Bitcoin.TestNet3 {
			networkName = "testnet"
		}
		if cfg.Bitcoin.RegTest {
			networkName = "regtest"
		}
		if cfg.Bitcoin.SimNet {
			networkName = "simnet"
		}

		primaryChain = "bitcoin"
	}

	// As soon as we're done parsing configuration options, ensure all paths
	// to directories and files are cleaned and expanded before attempting
	// to use them later on.
	cfg.DataDir = cleanAndExpandPath(cfg.DataDir)
	cfg.TLSCertPath = cleanAndExpandPath(cfg.TLSCertPath)
	cfg.AdminMacPath = cleanAndExpandPath(cfg.AdminMacPath)
	cfg.ReadMacPath = cleanAndExpandPath(cfg.ReadMacPath)
	cfg.InvoiceMacPath = cleanAndExpandPath(cfg.InvoiceMacPath)

	networkDir := filepath.Join(
		cfg.DataDir, defaultChainSubDirname,
		primaryChain,
		networkName,
	)

	// If a custom macaroon directory wasn't specified and the data
	// directory has changed from the default path, then we'll also update
	// the path for the macaroons to be generated.
	if cfg.AdminMacPath == "" {
		cfg.AdminMacPath = filepath.Join(
			networkDir, defaultAdminMacFilename,
		)
	}
	if cfg.ReadMacPath == "" {
		cfg.ReadMacPath = filepath.Join(
			networkDir, defaultReadMacFilename,
		)
	}
	if cfg.InvoiceMacPath == "" {
		cfg.InvoiceMacPath = filepath.Join(
			networkDir, defaultInvoiceMacFilename,
		)
	}

	// At least one RPCListener is required. So listen on localhost per
	// default.
	if len(cfg.RawRPCListeners) == 0 {
		addr := fmt.Sprintf("localhost:%d", defaultRPCPort)
		cfg.RawRPCListeners = append(cfg.RawRPCListeners, addr)
	}

	var err error
	// Add default port to all RPC listener addresses if needed and remove
	// duplicate addresses.
	cfg.RPCListeners, err = lncfg.NormalizeAddresses(
		cfg.RawRPCListeners, strconv.Itoa(defaultRPCPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}

	// Warn about missing config file only after all other configuration is
	// done.  This prevents the warning on help messages and invalid
	// options.  Note this should go directly before the return.
	if configFileError != nil {
		fmt.Println(configFileError)
	}

	return &cfg, nil
}

// cleanAndExpandPath expands environment variables and leading ~ in the
// passed path, cleans the result, and returns it.
// This function is taken from https://github.com/btcsuite/btcd
func cleanAndExpandPath(path string) string {
	if path == "" {
		return ""
	}

	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		var homeDir string

		user, err := user.Current()
		if err == nil {
			homeDir = user.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but the variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}
