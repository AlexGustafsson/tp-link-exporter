package tplink

import (
	"encoding/json"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type DeviceBroadcastResponse struct {
	System struct {
		SystemInfo *SystemInfo `json:"get_sysinfo"`
	} `json:"system"`
}

type SystemInfo struct {
	SoftwareVersion string `json:"sw_ver"`
	HardwareVersion string `json:"hw_ver"`
	Type            string `json:"type"`
	Model           string `json:"model"`
	DevelopmentName string `json:"dev_name"`
	IconHash        string `json:"icon_hash"`
	RelayState      int    `json:"relay_state"`
	OnTime          uint64 `json:"on_time"`
	ActiveMode      string `json:"active_mode"`
	Feature         string `json:"feature"`
	Updating        int    `json:"updating"`
	RSSI            int    `json:"rssi"`
	LedOff          int    `json:"led_off"`
	Alias           string `json:"alias"`
	DeviceID        string `json:"deviceId"`
	HardwareID      string `json:"hwId"`
	OEMID           string `json:"oemId"`
	Latitude        string `json:"latitude"`
	LatitudeInt     int    `json:"latitude_i"`
	Longitude       string `json:"longitude"`
	LongitudeInt    int    `json:"longitude_i"`
	MAC             string `json:"mac"`
	ErrorCode       int    `json:"err_code"`
}

type DeviceFinder struct {
	BroadcastAddress string

	socket    *net.UDPConn
	responses chan *SystemInfo
}

func NewDeviceFinder() *DeviceFinder {
	return &DeviceFinder{
		BroadcastAddress: "255.255.255.255",

		responses: make(chan *SystemInfo),
	}
}

func (f *DeviceFinder) Listen() error {
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

		var response DeviceBroadcastResponse
		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Errorf("Failed to unmarshal device broadcast response: %v", err)
			continue
		}

		f.responses <- response.System.SystemInfo
	}
}

func (f *DeviceFinder) Found() chan *SystemInfo {
	return f.responses
}

func (f *DeviceFinder) Broadcast() error {
	address, err := net.ResolveUDPAddr("udp4", f.BroadcastAddress)
	if err != nil {
		return err
	}

	log.Debugf("Broadcasting")
	message := EncryptMessage(DiscoveryMessage)
	_, err = f.socket.WriteTo(message, address)
	if err != nil {
		return err
	}

	return nil
}

func (f *DeviceFinder) broadcastContinously() {
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
