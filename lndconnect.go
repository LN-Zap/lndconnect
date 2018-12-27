package main

import (
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/glendc/go-external-ip"
	"github.com/skip2/go-qrcode"
)

type certificates struct {
	Cert     string `json:"c"`
	Macaroon string `json:"m"`
	Ip       string `json:"ip,omitempty"`
}

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

	ipString = net.JoinHostPort(
		ipString, fmt.Sprint(loadedConfig.LndConnect.Port),
	)

	urlString := fmt.Sprintf("lndconnect://%s?cert=%s&macaroon=%s", ipString, certificate, macaroonB64)

	if loadedConfig.LndConnect.Json {
		fmt.Println(urlString)
	} else if loadedConfig.LndConnect.Image {
		qrcode.WriteFile(urlString, qrcode.Medium, 512, "lndconnect-qr.png")
		fmt.Println("Wrote QR Code to file \"lndconnect-qr.png\"")
	} else {
		obj := qrcodeTerminal.New()
		obj.Get(urlString).Print()
		fmt.Println("\n⚠️  Press \"cmd + -\" a few times to see the full QR Code!\nIf that doesn't work run \"lndconnect -j\" to get a code you can copy paste into the app.")
	}
}
