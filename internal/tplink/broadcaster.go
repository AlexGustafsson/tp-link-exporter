package tplink

import (
	"encoding/json"
	"net"
	"time"

	"go.uber.org/zap"
)

type Broadcaster struct {
	targets   []string
	interval  time.Duration
	socket    *net.UDPConn
	responses chan *DeviceResponse
	log       *zap.Logger
}

func NewBroadcaster(targets []string, interval time.Duration, log *zap.Logger) *Broadcaster {
	return &Broadcaster{
		targets:   targets,
		interval:  interval,
		responses: make(chan *DeviceResponse),
		log:       log,
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
	f.log.Info("Listening for responses", zap.String("address", f.socket.LocalAddr().String()))

	go f.broadcastContinously()

	b := make([]byte, 1024)
	for {
		n, _, err := f.socket.ReadFromUDP(b)
		if err != nil {
			f.log.Error("Failed to read from socket", zap.Error(err))
			continue
		}

		message := DecryptMessage(b[:n])

		var response DeviceResponse
		err = json.Unmarshal(message, &response)
		if err != nil {
			f.log.Error("Failed to unmarshal device broadcast response", zap.String("message", string(message)), zap.Error(err))
			continue
		}
		f.responses <- &response
	}
}

func (f *Broadcaster) Responses() chan *DeviceResponse {
	return f.responses
}

func (f *Broadcaster) Broadcast() error {
	for _, address := range f.targets {
		resolvedAddress, err := net.ResolveUDPAddr("udp4", address+":9999")
		if err != nil {
			return err
		}

		f.log.Debug("Broadcasting", zap.String("address", resolvedAddress.String()))
		message := EncryptMessage(InfoMessage)
		_, err = f.socket.WriteTo(message, resolvedAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Broadcaster) broadcastContinously() {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := f.Broadcast()
		if err != nil {
			f.log.Error("Failed to broadcast", zap.Error(err))
		}
	}
}
