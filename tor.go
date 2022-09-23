package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/tor"
)

// createNewHiddenService automatically sets up a v3 onion service in
// order to listen for inbound connections over Tor.
func createNewHiddenService(loadedConfig *config, torController *tor.Controller) error {
	// If tor is active make a tor controller for creating onion REST service interface

	// Determine the different ports the server is listening on. The onion
	// service's virtual port will map to these ports and one will be picked
	// at random when the onion service is being accessed.
	restListenPorts := make([]int, 0, len(loadedConfig.RawRESTListeners))
	for _, RawRESTListener := range loadedConfig.RawRESTListeners {

		// Addresses can either be in network://address:port format,
		// network:address:port, address:port, or just port. We want to support
		// all possible types.
		var parsedPort string
		if strings.Contains(RawRESTListener, ":") {
			parts := strings.Split(RawRESTListener, ":")
			parsedPort = parts[len(parts)-1]
		} else {
			parsedPort = RawRESTListener
		}
		port, err := strconv.Atoi(parsedPort)
		if err != nil {
			return err
		}
		restListenPorts = append(restListenPorts, port)
	}

	// Once the port mapping has been set, we can go ahead and automatically
	// create our onion service. The service's private key will be saved to
	// disk in order to regain access to this service when restarting `lnd`.

	onionCfg := tor.AddOnionConfig{
		VirtualPort: defaultRESTPort,
		TargetPorts: restListenPorts,
		Store:       tor.NewOnionFile(loadedConfig.Tor.RESTKeyPath, 0o600, false, MockEncrypter{}),
		Type:        tor.V3,
	}

	onionAddr, err := torController.AddOnion(onionCfg)
	if err != nil {
		return err
	}
	addr := onionAddr.String()
	log.Println("Onion service created for LND REST interface")
	log.Printf("serviceID: %v", addr)

	err = updateHostAddrConfig(addr, loadedConfig)

	return err
}
