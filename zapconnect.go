package main

import (
// "flag"
"fmt"
"net"
"io/ioutil"
"github.com/Baozisoftware/qrcode-terminal-go"
"encoding/json"
b64 "encoding/base64"
"net/http"
"os"
"strings"
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
		return
	}

	certBytes, err := ioutil.ReadFile(loadedConfig.TLSCertPath)
	if err != nil {
		fmt.Print(err)
		return
	}

	macBytes, err := ioutil.ReadFile(loadedConfig.AdminMacPath)
	if err != nil {
		fmt.Print(err)
		return
	}

	sEnc := b64.StdEncoding.EncodeToString([]byte(macBytes))

	ipString := ""
	if loadedConfig.LocalIp {
		ipString = getLocalIP() + ":10009"
	} else {
		ipString = getPublicIP() + ":10009"
	}

	cert := &certificates{
		Cert:	  string(certBytes),
		Macaroon: sEnc,
		Ip:		  ipString}
	certB, _ := json.Marshal(cert)


	if loadedConfig.Json {
		fmt.Println(string(certB))
	} else {
		obj := qrcodeTerminal.New()
		obj.Get(string(certB)).Print()
	}
}
