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

	externalip "github.com/glendc/go-external-ip"
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

func getURI(loadedConfig *config) (string, error) {
	var err error

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
			return "", err
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
		return "", err
	}

	macaroonB64 := b64.RawURLEncoding.EncodeToString([]byte(macBytes))

	q.Add("macaroon", macaroonB64)

	// custom query
	for _, s := range loadedConfig.LndConnect.Query {
		queryParts := strings.Split(s, "=")

		if len(queryParts) != 2 {
			return "", fmt.Errorf("Invalid Query Argument: %s", s)
		}

		q.Add(queryParts[0], queryParts[1])
	}

	u.RawQuery = q.Encode()

	fmt.Println("\nURI generated successfully.")
	return u.String(), nil
}
