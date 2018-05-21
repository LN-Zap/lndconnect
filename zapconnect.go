package main

import (
"flag"
"fmt"
"io/ioutil"
"github.com/Baozisoftware/qrcode-terminal-go"
"encoding/json"
b64 "encoding/base64"
)

type certificates struct {
	Cert		string `json:"c"`
	Macaroon	string `json:"m"`
}

func main() {
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

	cert := &certificates{
		Cert:	  string(certBytes),
		Macaroon: sEnc}
	certB, _ := json.Marshal(cert)

	jsonPtr := flag.Bool("j", false, "Generate json instead of a QRCode.")
	flag.Parse()

    if *jsonPtr {
    	fmt.Print(string(certB))
    } else {
		obj := qrcodeTerminal.New()
		obj.Get(string(certB)).Print()
	}
}
