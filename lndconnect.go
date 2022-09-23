package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/lightningnetwork/lnd/tor"
)

func main() {
	loadedConfig, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// If createonion option is selected, tor is active and v3 onion services have been
	// specified, make a tor controller and pass it into the REST controller server
	var torController *tor.Controller
	if loadedConfig.LndConnect.CreateOnion && loadedConfig.Tor.Active && loadedConfig.Tor.V3 {
		var targetIPAddress string
		if net.ParseIP(loadedConfig.Tor.TargetIPAddress) == nil {
			addrs, err := loadedConfig.net.LookupHost(loadedConfig.Tor.TargetIPAddress)
			if err != nil {
				log.Fatalln(err)
			}
			targetIPAddress = addrs[0]
			log.Printf(
				"`tor.targetipaddress` doesn't define an IP address, hostname %s was resolved to %s",
				loadedConfig.Tor.TargetIPAddress,
				targetIPAddress,
			)
		} else {
			targetIPAddress = loadedConfig.Tor.TargetIPAddress
		}
		torController = tor.NewController(
			loadedConfig.Tor.Control,
			targetIPAddress,
			loadedConfig.Tor.Password,
		)

		// Start the tor controller before giving it to any other
		// subsystems.
		if err := torController.Start(); err != nil {
			log.Fatalf("error starting tor controller: %v", err)
		}
		defer func() {
			if err := torController.Stop(); err != nil {
				log.Printf("error stopping tor controller: %v", err)
			}
		}()

		if err := createNewHiddenService(loadedConfig, torController); err != nil {
			log.Fatal(err)
		}
	}

	// Generate URI
	uri, err := getURI(loadedConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Print URI or QR Code to selected output
	if loadedConfig.LndConnect.URL {
		log.Println(uri)
	} else {
		err = getQR(uri, loadedConfig.LndConnect.Image)
		if err != nil {
			log.Println(err)
		}
	}
	if torController != nil {
		cancelChan := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		// catch SIGINT, SIGQUIT or SIGETRM
		signal.Notify(cancelChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		go func() {
			sig := <-cancelChan
			log.Printf("Caught %v signal", sig)
			done <- true
		}()
		<-done
		log.Println("lndconnect is shutting down now...")
	}
}
