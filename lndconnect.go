package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	loadedConfig, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if loadedConfig.LndConnect.CreateOnion {
		addr, err := createNewHiddenService(loadedConfig)
		if err != nil {
			log.Fatal(err)
		}
		parsedAddr := strings.Split(addr, ":")
		loadedConfig.LndConnect.Host = parsedAddr[0]
		parsedPort, err := strconv.ParseUint(parsedAddr[1], 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		loadedConfig.LndConnect.Port = uint16(parsedPort)
		loadedConfig.LndConnect.NoCert = true
	}

	// Generate URI
	uri, err := getURI(loadedConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Print URI or QR Code to selected output
	if loadedConfig.LndConnect.Url {
		fmt.Println(uri)

	} else {
		getQR(uri, loadedConfig.LndConnect.Image)
	}
}
