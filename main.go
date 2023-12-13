package main

import (
	"log"
	"net"

	"github.com/pion/turn/v2"
)

func main() {
	// Set up the listener for the server. This will listen on all interfaces on port 3478.
	listener, err := net.ListenPacket("udp", "0.0.0.0:3478")
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Define a simple auth handler. In a real application, you would use a more robust authentication mechanism.
	authHandler := func(username, realm string, srcAddr net.Addr) ([]byte, bool) {
		log.Printf("Auth request from %v for user %s", srcAddr, username)
		if username == "user" {
			// The password must be a MD5 hash of "username:realm:password".
			// In this case, the password is "secret".
			return turn.GenerateAuthKey(username, realm, "secret"), true
		}
		return nil, false
	}

	// Create a TURN server instance.
	_, err = turn.NewServer(turn.ServerConfig{
		Realm:       "tribes.xyz",
		AuthHandler: authHandler,
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: listener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP("0.0.0.0"),
					Address:      "0.0.0.0",
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("Failed to create TURN server: %v", err)
	}

	log.Println("STUN/TURN server is up and running on 0.0.0.0:3478")
	select {} // Block forever
}
