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
	defaultTLSCertFilename    = "tls.cert"
	defaultAdminMacFilename   = "admin.macaroon"
)

var (
	defaultLndDir     = btcutil.AppDataDir("lnd", false)
	defaultConfigFile = filepath.Join(defaultLndDir, defaultConfigFilename)
	defaultTLSCertPath = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultAdminMacPath   = filepath.Join(defaultLndDir, defaultAdminMacFilename)
)

// config defines the configuration options for lnd.
//
// See loadConfig for further details regarding the configuration
// loading+parsing process.
type config struct {
	LndDir         string   `long:"lnddir" description:"The base directory that contains lnd's data, logs, configuration file, etc."`
	ConfigFile     string   `long:"C" long:"configfile" description:"Path to configuration file"`
	TLSCertPath    string   `long:"tlscertpath" description:"Path to write the TLS certificate for lnd's RPC and REST services"`
	AdminMacPath   string   `long:"adminmacaroonpath" description:"Path to write the admin macaroon for lnd's RPC and REST services if it doesn't exist"`
}

func loadConfig() (*config, error) {
	defaultCfg := config{
		LndDir:         defaultLndDir,
		ConfigFile:     defaultConfigFile,
		TLSCertPath:    defaultTLSCertPath,
		AdminMacPath:   defaultAdminMacPath,
	}

	cfg := defaultCfg
	configFile := cleanAndExpandPath(defaultCfg.ConfigFile)
	flags.IniParse(configFile, &cfg)

	cfg.TLSCertPath = cleanAndExpandPath(cfg.TLSCertPath)
	cfg.AdminMacPath = cleanAndExpandPath(cfg.AdminMacPath)

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
