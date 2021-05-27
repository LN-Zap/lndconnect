package main

import (
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"
	"crypto/sha256"
	"strconv"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/glendc/go-external-ip"
	"github.com/skip2/go-qrcode"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getPublicIP() string {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ip.String()
}

func main() {
	loadedConfig, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	displayLink(loadedConfig)
}

func displayLink(loadedConfig *config) {
	var err error

	if loadedConfig.LndConnect.Version == 1 {
		// host
		ipString := ""
		if loadedConfig.LndConnect.Host != "" {
			ipString = loadedConfig.LndConnect.Host
		} else if loadedConfig.LndConnect.LocalIp {
			ipString = getLocalIP()
		} else if loadedConfig.LndConnect.Localhost {
			ipString = "127.0.0.1"
		} else {
			ipString = getPublicIP()
		}

		ipString = net.JoinHostPort(ipString, fmt.Sprint(loadedConfig.LndConnect.Port))

		u := url.URL{Scheme: "lndconnect", Host: ipString}
		q := u.Query()

		// cert
		if !loadedConfig.LndConnect.NoCert {
			certBytes, err := ioutil.ReadFile(loadedConfig.TLSCertPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			block, _ := pem.Decode(certBytes)
			if block == nil || block.Type != "CERTIFICATE" {
				fmt.Println("failed to decode PEM block containing certificate")
			}

			certificate := b64.RawURLEncoding.EncodeToString([]byte(block.Bytes))

			q.Add("cert", certificate)
		}

		// macaroon
		var macBytes []byte
		if loadedConfig.LndConnect.Invoice {
			macBytes, err = ioutil.ReadFile(loadedConfig.InvoiceMacPath)
		} else if loadedConfig.LndConnect.Readonly {
			macBytes, err = ioutil.ReadFile(loadedConfig.ReadMacPath)
		} else {
			macBytes, err = ioutil.ReadFile(loadedConfig.AdminMacPath)
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		macaroonB64 := b64.RawURLEncoding.EncodeToString([]byte(macBytes))

		q.Add("macaroon", macaroonB64)

		// custom query
		for _, s := range loadedConfig.LndConnect.Query {
			queryParts := strings.Split(s, "=")

			if len(queryParts) != 2 {
				fmt.Println("Invalid Query Argument:", s)
				return
			}

			q.Add(queryParts[0], queryParts[1])
		}

		u.RawQuery = q.Encode()

		// generate link / QR Code
		if loadedConfig.LndConnect.Url {
			fmt.Println(u.String())
		} else if loadedConfig.LndConnect.Image {
			qrcode.WriteFile(u.String(), qrcode.Low, 512, "lndconnect-qr.png")
			fmt.Println("Wrote QR Code to file \"lndconnect-qr.png\"")
		} else {
			obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite,qrcodeTerminal.QRCodeRecoveryLevels.Low)
			obj.Get(u.String()).Print()
			fmt.Println("\n⚠️  Press \"cmd + -\" a few times to see the full QR Code!\nIf that doesn't work run \"lndconnect -j\" to get a code you can copy paste into the app.")
		}
	} else { // Version 2 and higher...
		// host
		ipString := ""
		if loadedConfig.LndConnect.Host != "" {
			ipString = loadedConfig.LndConnect.Host
		} else if loadedConfig.LndConnect.LocalIp {
			ipString = getLocalIP()
		} else if loadedConfig.LndConnect.Localhost {
			ipString = "127.0.0.1"
		} else {
			ipString = getPublicIP()
		}

		ipString = net.JoinHostPort(ipString, fmt.Sprint(loadedConfig.LndConnect.Port))

		u := url.URL{Scheme: "lndconn", Host: ipString}
		q := u.Query()

		// version
		q.Add("v", strconv.Itoa(int(loadedConfig.LndConnect.Version)))

		// cert
		if !loadedConfig.LndConnect.NoCert {
			certBytes, err := ioutil.ReadFile(loadedConfig.TLSCertPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			block, _ := pem.Decode(certBytes)
			if block == nil || block.Type != "CERTIFICATE" {
				fmt.Println("failed to decode PEM block containing certificate")
			}

			hash := sha256.Sum256(block.Bytes)
			certificate := b64.RawURLEncoding.EncodeToString(hash[:])

			q.Add("c", certificate)
		}

		// macaroon
		var macBytes []byte
		if loadedConfig.LndConnect.Invoice {
			macBytes, err = ioutil.ReadFile(loadedConfig.InvoiceMacPath)
		} else if loadedConfig.LndConnect.Readonly {
			macBytes, err = ioutil.ReadFile(loadedConfig.ReadMacPath)
		} else {
			macBytes, err = ioutil.ReadFile(loadedConfig.AdminMacPath)
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		macaroonB64 := b64.RawURLEncoding.EncodeToString([]byte(macBytes))

		q.Add("m", macaroonB64)

		// custom query
		for _, s := range loadedConfig.LndConnect.Query {
			queryParts := strings.Split(s, "=")

			if len(queryParts) != 2 {
				fmt.Println("Invalid Query Argument:", s)
				return
			}

			q.Add(queryParts[0], queryParts[1])
		}

		u.RawQuery = q.Encode()

		// generate link / QR Code
		if loadedConfig.LndConnect.Url {
			fmt.Println(u.String())
		} else if loadedConfig.LndConnect.Image {
			qrcode.WriteFile(u.String(), qrcode.Low, 512, "lndconnect-qr.png")
			fmt.Println("Wrote QR Code to file \"lndconnect-qr.png\"")
		} else {
			obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite,qrcodeTerminal.QRCodeRecoveryLevels.Low)
			obj.Get(u.String()).Print()
			fmt.Println("\n⚠️  Press \"cmd + -\" a few times to see the full QR Code!\nIf that doesn't work run \"lndconnect -j\" to get a code you can copy paste into the app.")
		}
	}
}
