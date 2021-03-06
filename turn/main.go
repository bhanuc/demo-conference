package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
    "errors"
	"github.com/pion/logging"
	"github.com/pion/turn/v2"
)

func createAuthHandler() turn.AuthHandler {
	return func(username string,  realm string, srcAddr net.Addr) (key []byte, ok bool) {
		return []byte("password"), true
	}
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

	udpListener, err := net.ListenPacket("udp4", serverIP+":"+strconv.Itoa(*port))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm:              realm,
		AuthHandler:        createAuthHandler(),
		ChannelBindTimeout: channelBindTimeout,
		PacketConnConfigs: []PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
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
