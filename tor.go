package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/tor"
)

// createNewHiddenService automatically sets up a v3 onion service in
// order to listen for inbound connections over Tor.
func createNewHiddenService(loadedConfig *config) (string, error) {

	// If tor is active make a tor controller for creating onion REST service interface
	var torController *tor.Controller
	if loadedConfig.Tor.Active && loadedConfig.Tor.V3 {
		torController = tor.NewController(
			loadedConfig.Tor.Control, loadedConfig.Tor.TargetIPAddress,
			loadedConfig.Tor.Password,
		)

		// Start the tor controller before giving it to any other
		// subsystems.
		if err := torController.Start(); err != nil {
			return "", err
		}
		defer func() {
			if err := torController.Stop(); err != nil {
				fmt.Printf("error stopping tor "+
					"controller: %v", err)
			}
		}()
	} else {
		return "", fmt.Errorf("Either tor is inactive or v3 onion service are disabled, check configuration.")
	}

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
			return "", err
		}
		restListenPorts = append(restListenPorts, port)
	}

	// Once the port mapping has been set, we can go ahead and automatically
	// create our onion service. The service's private key will be saved to
	// disk in order to regain access to this service when restarting `lnd`.

	onionCfg := tor.AddOnionConfig{
		VirtualPort: defaultRESTPort,
		TargetPorts: restListenPorts,
		Store:       tor.NewOnionFile(loadedConfig.Tor.RESTKeyPath, 0600, false, MockEncrypter{}),
		Type:        tor.V3,
	}

	addr, err := torController.AddOnion(onionCfg)
	if err != nil {
		return "", err
	}
	fmt.Println("\nOnion service created for LND REST interface, address reads:")
	fmt.Println(addr)

	return addr.String(), nil
}
