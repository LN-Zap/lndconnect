package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lightningnetwork/lnd/tor"
)

// var wg sync.WaitGroup

func main() {
	loadedConfig, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// If createonion option is selected, tor is active and v3 onion services have been
	// specified, make a tor controller and pass it into the REST controller server
	var torController *tor.Controller
	if loadedConfig.LndConnect.CreateOnion && loadedConfig.Tor.Active && loadedConfig.Tor.V3 {
		torController = tor.NewController(
			loadedConfig.Tor.Control,
			loadedConfig.Tor.TargetIPAddress,
			loadedConfig.Tor.Password,
		)

		// Start the tor controller before giving it to any other
		// subsystems.
		if err := torController.Start(); err != nil {
			log.Fatalf("error starting tor controller: %v", err)
		}
		// defer func() {
		// 	if err := torController.Stop(); err != nil {
		// 		log.Printf("error stopping tor controller: %v", err)
		// 	}
		// }()
	}

	if torController != nil {
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
		// catch SIGINT, SIGQUIT or SIGETRM
		signal.Notify(cancelChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		done := make(chan bool, 1)
		go func() {
			sig := <-cancelChan
			log.Printf("Caught SIGTERM %v", sig)
			done <- true
		}()
		<-done
		log.Println("exiting...")
		if err := torController.Stop(); err != nil {
			log.Printf("error stopping tor controller: %v", err)
		}

	}
}
