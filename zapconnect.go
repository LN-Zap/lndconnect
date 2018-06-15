package main

import (
"flag"
"fmt"
"net"
"io/ioutil"
"github.com/Baozisoftware/qrcode-terminal-go"
"encoding/json"
b64 "encoding/base64"
)

type certificates struct {
	Cert		string `json:"c"`
	Macaroon	string `json:"m"`
	Ip 			string `json:"ip,omitempty"`
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

func main() {
	ipPtr := flag.Bool("i", false, "Include local ip in QRCode.")
	jsonPtr := flag.Bool("j", false, "Generate json instead of a QRCode.")
	flag.Parse()

	loadedConfig, _ := loadConfig()

	certBytes, err := ioutil.ReadFile(loadedConfig.TLSCertPath)
	if err != nil {
		fmt.Print(err)
	}

	macBytes, err := ioutil.ReadFile(loadedConfig.AdminMacPath)
	if err != nil {
		fmt.Print(err)
	}

	sEnc := b64.StdEncoding.EncodeToString([]byte(macBytes))

	ipString := ""
	if *ipPtr {
		ipString = getLocalIP() + ":10009"
	}

	cert := &certificates{
		Cert:	  string(certBytes),
		Macaroon: sEnc,
		Ip:		  ipString}
	certB, _ := json.Marshal(cert)


	if *jsonPtr {
		fmt.Println(string(certB))
	} else {
		obj := qrcodeTerminal.New()
		obj.Get(string(certB)).Print()
	}
}
