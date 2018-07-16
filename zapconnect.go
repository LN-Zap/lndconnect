package main

import (
	"fmt"
	"net"
	"io/ioutil"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"encoding/json"
	b64 "encoding/base64"
	"net/http"
	"os"
	"strings"
	"encoding/pem"
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

func getPublicIP() string {
	resp, err := http.Get("http://ipv4.myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
	    bodyBytes, _ := ioutil.ReadAll(resp.Body)
	    return strings.TrimSpace(string(bodyBytes))
	}

	return ""
}

func main() {
	loadedConfig, err := loadConfig()
	if err != nil {
		fmt.Println()
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

	certificate := b64.StdEncoding.EncodeToString([]byte(block.Bytes))

	macBytes, err := ioutil.ReadFile(loadedConfig.AdminMacPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	macaroonB64 := b64.StdEncoding.EncodeToString([]byte(macBytes))

	ipString := ""
	if loadedConfig.LocalIp {
		ipString = getLocalIP() + ":10009"
	} else {
		ipString = getPublicIP() + ":10009"
	}

	cert := &certificates{
		Cert:     certificate,
		Macaroon: macaroonB64,
		Ip:       ipString}
	certB, _ := json.Marshal(cert)


	if loadedConfig.Json {
		fmt.Println(string(certB))
	} else {
		obj := qrcodeTerminal.New()
		obj.Get(string(certB)).Print()
	}
}
