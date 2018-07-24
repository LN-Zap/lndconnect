package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/roasbeef/btcutil"
)

const (
	defaultConfigFilename     = "lnd.conf"
	defaultDataDirname        = "data"
	defaultTLSCertFilename    = "tls.cert"
	defaultAdminMacFilename   = "admin.macaroon"
)

var (
	defaultLndDir     = btcutil.AppDataDir("lnd", false)
	defaultConfigFile = filepath.Join(defaultLndDir, defaultConfigFilename)
	defaultDataDir    = filepath.Join(defaultLndDir, defaultDataDirname)
	defaultTLSCertPath = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultAdminMacPath   = filepath.Join(defaultLndDir, defaultAdminMacFilename)
)

// config defines the configuration options for lnd.
//
// See loadConfig for further details regarding the configuration
// loading+parsing process.
type config struct {
	LocalIp        bool     `short:"i" long:"localip" description:"Include local ip in QRCode."`
	Localhost      bool     `short:"l" long:"localhost" description:"Use 127.0.0.1 for ip."`
	Json           bool     `short:"j" long:"json" description:"Generate json instead of a QRCode."`
	LndDir         string   `long:"lnddir" description:"The base directory that contains lnd's data, logs, configuration file, etc."`
	ConfigFile     string   `long:"C" long:"configfile" description:"Path to configuration file"`
	DataDir        string   `short:"b" long:"datadir" description:"The directory to store lnd's data within"`
	TLSCertPath    string   `long:"tlscertpath" description:"Path to write the TLS certificate for lnd's RPC and REST services"`
	AdminMacPath   string   `long:"adminmacaroonpath" description:"Path to write the admin macaroon for lnd's RPC and REST services if it doesn't exist"`
}

func loadConfig() (*config, error) {
	defaultCfg := config{
		LndDir:         defaultLndDir,
		ConfigFile:     defaultConfigFile,
		DataDir:        defaultDataDir,
		TLSCertPath:    defaultTLSCertPath,
		AdminMacPath:   defaultAdminMacPath,
	}

	// Pre-parse the command line options to pick up an alternative config
	// file.
	preCfg := defaultCfg
	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	cfg := preCfg
	configFile := cleanAndExpandPath(defaultCfg.ConfigFile)
	flags.IniParse(configFile, &cfg)

	cfg.TLSCertPath = cleanAndExpandPath(cfg.TLSCertPath)
	cfg.AdminMacPath = cleanAndExpandPath(cfg.AdminMacPath)

	// If a custom macaroon directory wasn't specified and the data
	// directory has changed from the default path, then we'll also update
	// the path for the macaroons to be generated.
	if cfg.DataDir != defaultDataDir && cfg.AdminMacPath == defaultAdminMacPath {
		cfg.AdminMacPath = filepath.Join(
			cfg.DataDir, defaultAdminMacFilename,
		)
	}

	return &cfg, nil
}

// cleanAndExpandPath expands environment variables and leading ~ in the
// passed path, cleans the result, and returns it.
// This function is taken from https://github.com/btcsuite/btcd
func cleanAndExpandPath(path string) string {
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
