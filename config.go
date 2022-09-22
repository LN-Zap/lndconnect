// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (C) 2015-2017 The Lightning Network Developers

package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd/tor"
)

const (
	defaultConfigFilename     = "lnd.conf"
	defaultDataDirname        = "data"
	defaultChainSubDirname    = "chain"
	defaultTLSCertFilename    = "tls.cert"
	defaultAdminMacFilename   = "admin.macaroon"
	defaultReadMacFilename    = "readonly.macaroon"
	defaultInvoiceMacFilename = "invoice.macaroon"
	defaultRPCPort            = 10009
	defaultRESTPort           = 8080
	defaultRESTKeyFileName    = "rest_onion_private_key"
	defaultQRFileName         = "lndconnect-qr.png"
)

var (
	defaultLndDir      = btcutil.AppDataDir("lnd", false)
	defaultConfigFile  = filepath.Join(defaultLndDir, defaultConfigFilename)
	defaultDataDir     = filepath.Join(defaultLndDir, defaultDataDirname)
	defaultTLSCertPath = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultRESTKeyPath = filepath.Join(defaultLndDir, defaultRESTKeyFileName)
	defaultQRFilePath  = filepath.Join(defaultLndDir, defaultQRFileName)
)

type chainConfig struct {
	Active   bool `long:"active" description:"If the chain should be active or not"`
	MainNet  bool `long:"mainnet" description:"Use the main network"`
	TestNet3 bool `long:"testnet" description:"Use the test network"`
	SimNet   bool `long:"simnet" description:"Use the simulation test network"`
	RegTest  bool `long:"regtest" description:"Use the regression test network"`
}

type torConfig struct {
	Active          bool   `long:"active" description:"Allow outbound and inbound connections to be routed through Tor"`
	V3              bool   `long:"v3" description:"Automatically set up a v3 onion service to listen for inbound connections"`
	Control         string `long:"control" description:"The host:port that Tor is listening on for Tor control connections"`
	TargetIPAddress string `long:"targetipaddress" description:"IP address that Tor should use as the target of the hidden service"`
	Password        string `long:"password" description:"If provided, the HASHEDPASSWORD authentication method will be used instead of the SAFECOOKIE one."`
	RESTKeyPath     string `short:"r" long:"restkeypath" description:"The path to the private key of the onion service being created if provided."`
}

type arrayFlags []string

type lndConnectConfig struct {
	LocalIp     bool       `short:"i" long:"localip" description:"Include local ip in QRCode"`
	Localhost   bool       `short:"l" long:"localhost" description:"Use 127.0.0.1 for ip"`
	Host        string     `long:"host" description:"Use specific host name"`
	NoCert      bool       `long:"nocert" description:"Don't include the certificate"`
	Port        uint16     `short:"p" long:"port" description:"Use this port"`
	Url         bool       `short:"j" long:"url" description:"Display url instead of a QRCode"`
	Image       bool       `short:"o" long:"image" description:"Output QRCode to file"`
	Invoice     bool       `long:"invoice" description:"Use invoice macaroon"`
	Readonly    bool       `long:"readonly" description:"Use readonly macaroon"`
	Query       arrayFlags `short:"q" long:"query" description:"Add additional url query parameters"`
	CreateOnion bool       `short:"c" long:"createonion" description:"Create onion v3 hidden service to access REST interface."`
}

// config defines the configuration options for lndconnect.
//
// See loadConfig for further details regarding the configuration
// loading+parsing process.
type config struct {
	LndConnect *lndConnectConfig `group:"LndConnect"`

	LndDir           string   `long:"lnddir" description:"The base directory that contains lnd's data, logs, configuration file, etc."`
	ConfigFile       string   `long:"C" long:"configfile" description:"Path to configuration file"`
	DataDir          string   `short:"b" long:"datadir" description:"The directory to find lnd's data within"`
	TLSCertPath      string   `long:"tlscertpath" description:"Path to read the TLS certificate from"`
	AdminMacPath     string   `long:"adminmacaroonpath" description:"Path to read the admin macaroon from"`
	ReadMacPath      string   `long:"readonlymacaroonpath" description:"Path to read the read-only macaroon from"`
	InvoiceMacPath   string   `long:"invoicemacaroonpath" description:"Path to read the invoice-only macaroon from"`
	RawRESTListeners []string `long:"restlisten" description:"Interface/Port/Socket listening for REST connections"`

	Bitcoin  *chainConfig `group:"Bitcoin" namespace:"bitcoin"`
	Litecoin *chainConfig `group:"Litecoin" namespace:"litecoin"`

	Tor *torConfig `group:"Tor" namespace:"tor"`

	// The following lines we only need to be able to parse the
	// configuration INI file without errors. The content will be ignored.
	BtcdMode      *chainConfig `hidden:"true" group:"btcd" namespace:"btcd"`
	BitcoindMode  *chainConfig `hidden:"true" group:"bitcoind" namespace:"bitcoind"`
	NeutrinoMode  *chainConfig `hidden:"true" group:"neutrino" namespace:"neutrino"`
	LtcdMode      *chainConfig `hidden:"true" group:"ltcd" namespace:"ltcd"`
	LitecoindMode *chainConfig `hidden:"true" group:"litecoind" namespace:"litecoind"`
	Autopilot     *chainConfig `hidden:"true" group:"Autopilot" namespace:"autopilot"`
	Hodl          *chainConfig `hidden:"true" group:"hodl" namespace:"hodl"`

	net tor.Net
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
//  1. Start with a default config with sane settings
//  2. Pre-parse the command line to check for an alternative config file
//  3. Load configuration file overwriting defaults with any specified options
//  4. Parse CLI options and overwrite/add any specified options
func loadConfig() (*config, error) {
	defaultCfg := config{
		LndConnect: &lndConnectConfig{
			Port: defaultRPCPort,
		},
		LndDir:      defaultLndDir,
		ConfigFile:  defaultConfigFile,
		DataDir:     defaultDataDir,
		TLSCertPath: defaultTLSCertPath,
		Tor: &torConfig{
			RESTKeyPath: defaultRESTKeyPath,
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
		preCfg.Tor.RESTKeyPath = filepath.Join(lndDir, defaultRESTKeyFileName)
	}

	// Next, load any additional configuration options from the file.
	var configFileError error
	cfg := preCfg

	// We don't have a full representation of all LND options in lndconnect
	// so while parsing the config file, we only take what we need, ignoring
	// all the unknown (to us) options.
	p := flags.NewParser(&cfg, flags.IgnoreUnknown)
	if err := flags.NewIniParser(p).ParseFile(cfg.ConfigFile); err != nil {
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
		numNets := 0
		if cfg.Bitcoin.MainNet {
			numNets++
			networkName = "mainnet"
		}
		if cfg.Bitcoin.TestNet3 {
			numNets++
			networkName = "testnet"
		}
		if cfg.Bitcoin.RegTest {
			numNets++
			networkName = "regtest"
		}
		if cfg.Bitcoin.SimNet {
			numNets++
			networkName = "simnet"
		}
		if numNets > 1 {
			str := "The mainnet, testnet, regtest, and " +
				"simnet params can't be used together -- " +
				"choose one of the four"
			err := fmt.Errorf(str)
			return nil, err
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

		u, err := user.Current()
		if err == nil {
			homeDir = u.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but the variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}
