package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/skip2/go-qrcode"
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

	} else if loadedConfig.LndConnect.Image {
		BrightGreen := color.RGBA{95, 191, 95, 255}
		qrcode.WriteColorFile(uri, qrcode.Low, 512, BrightGreen, color.Black, defaultQRFilePath)
		fmt.Printf("\nWrote QR Code to file \"%s\"", defaultQRFilePath)

	} else {
		obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Low)
		obj.Get(uri).Print()
		fmt.Println("\n⚠️  Press \"cmd + -\" a few times to see the full QR Code!\nIf that doesn't work run \"lndconnect -j\" to get a code you can copy paste into the app.")
	}

}
