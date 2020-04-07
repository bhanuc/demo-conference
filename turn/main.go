package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pion/logging"
	"github.com/pion/turn"
)

func createAuthHandler() turn.AuthHandler {
	return func(username string,  realm string, srcAddr net.Addr) (key []byte, ok bool) {
		return []byte("password"), true
	}
}

// RelayAddressGenerator is used to generate a RelayAddress when creating an allocation.
// You can use one of the provided ones or provide your own.
type RelayAddressGenerator interface {
	// Validate confirms that the RelayAddressGenerator is properly initialized
	Validate() error

	// Allocate a PacketConn (UDP) RelayAddress
	AllocatePacketConn(network string, requestedPort int) (net.PacketConn, net.Addr, error)

	// Allocate a Conn (TCP) RelayAddress
	AllocateConn(network string, requestedPort int) (net.Conn, net.Addr, error)
}

type PacketConnConfig struct {
	PacketConn net.PacketConn

	// When an allocation is generated the RelayAddressGenerator
	// creates the net.PacketConn and returns the IP/Port it is available at
	RelayAddressGenerator RelayAddressGenerator
}

// ListenerConfig is a single net.Listener to accept connections on. This will be used for TCP, TLS and DTLS listeners
type ListenerConfig struct {
	Listener net.Listener

	// When an allocation is generated the RelayAddressGenerator
	// creates the net.PacketConn and returns the IP/Port it is available at
	RelayAddressGenerator RelayAddressGenerator
}

func (c *ListenerConfig) validate() error {
	if c.Listener == nil {
		return errListenerUnset
	}

	if c.RelayAddressGenerator == nil {
		return errRelayAddressGeneratorUnset
	}

	return c.RelayAddressGenerator.Validate()
}


func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	realm := os.Getenv("REALM")
	if realm == "" {
		log.Panic("REALM is a required environment variable")
	}


	serverIP := os.Getenv("SERVERIP")
	if serverIP == "" {
		serverIP = "127.0.0.1"
	}


	udpPortStr := os.Getenv("UDP_PORT")
	if udpPortStr == "" {
		udpPortStr = "3478"
	}
	udpPort, err := strconv.Atoi(udpPortStr)
	if err != nil {
		log.Panic(err)
	}

	var channelBindTimeout time.Duration
	channelBindTimeoutStr := os.Getenv("CHANNEL_BIND_TIMEOUT")
	if channelBindTimeoutStr != "" {
		channelBindTimeout, err = time.ParseDuration(channelBindTimeoutStr)
		if err != nil {
			log.Panicf("CHANNEL_BIND_TIMEOUT=%s is an invalid time Duration", channelBindTimeoutStr)
		}
	}

	s := turn.NewServer(turn.ServerConfig{
		Realm:              realm,
		AuthHandler:        createAuthHandler(),
		ChannelBindTimeout: channelBindTimeout,
		PacketConnConfigs: []PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(serverIP),
					Address:      "0.0.0.0",
				},
			},
		},
		LoggerFactory:      logging.NewDefaultLoggerFactory(),
	})

	err = s.Start()
	if err != nil {
		log.Panic(err)
	}

	<-sigs

	err = s.Close()
	if err != nil {
		log.Panic(err)
	}
}
