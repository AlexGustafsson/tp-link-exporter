package tplink

import (
	"encoding/json"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Broadcaster struct {
	BroadcastAddress string

	socket    *net.UDPConn
	responses chan *DeviceResponse
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		BroadcastAddress: "255.255.255.255",

		responses: make(chan *DeviceResponse),
	}
}

func (f *Broadcaster) Listen() error {
	listenAddress, err := net.ResolveUDPAddr("udp4", ":0")
	if err != nil {
		return err
	}

	f.socket, err = net.ListenUDP("udp4", listenAddress)
	if err != nil {
		return err
	}
	log.Infof("Listen address: %v", f.socket.LocalAddr())

	go f.broadcastContinously()

	b := make([]byte, 1024)
	for {
		n, _, err := f.socket.ReadFromUDP(b)
		if err != nil {
			log.Errorf("Failed to read from socket: %v", err)
			continue
		}

		message := DecryptMessage(b[:n])

		var response DeviceResponse
		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Errorf("Failed to unmarshal device broadcast response: %s, %v", message, err)
			continue
		}
		f.responses <- &response
	}
}

func (f *Broadcaster) Responses() chan *DeviceResponse {
	return f.responses
}

func (f *Broadcaster) Broadcast() error {
	address, err := net.ResolveUDPAddr("udp4", f.BroadcastAddress)
	if err != nil {
		return err
	}

	log.Debugf("Broadcasting")
	message := EncryptMessage(InfoMessage)
	_, err = f.socket.WriteTo(message, address)
	if err != nil {
		return err
	}

	return nil
}

func (f *Broadcaster) broadcastContinously() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := f.Broadcast()
		if err != nil {
			log.Errorf("Failed to broadcast: %v", err)
		}
	}
}
