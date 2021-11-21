package tplink

// InfoMessage requests all available information from the device. It has a 4 byte unimportant padding which is also present in the response message.
var InfoMessage = []byte(`    {"emeter":{"get_realtime":{}},"system":{"get_sysinfo":{}}}`)

type DeviceResponse struct {
	Device struct {
		Info *DeviceInfo `json:"get_sysinfo"`
	} `json:"system"`
	EnergyMeter struct {
		Info *EnergyMeterInfo `json:"get_realtime"`
	} `json:"emeter"`
}

type DeviceInfo struct {
	SoftwareVersion string  `json:"sw_ver"`
	HardwareVersion string  `json:"hw_ver"`
	Type            string  `json:"type"`
	Model           string  `json:"model"`
	DevelopmentName string  `json:"dev_name"`
	IconHash        string  `json:"icon_hash"`
	RelayState      int     `json:"relay_state"`
	OnTime          uint64  `json:"on_time"`
	ActiveMode      string  `json:"active_mode"`
	Feature         string  `json:"feature"`
	Updating        int     `json:"updating"`
	RSSI            float64 `json:"rssi"`
	LedOff          int     `json:"led_off"`
	Alias           string  `json:"alias"`
	DeviceID        string  `json:"deviceId"`
	HardwareID      string  `json:"hwId"`
	OEMID           string  `json:"oemId"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	MAC             string  `json:"mac"`
	ErrorCode       int     `json:"err_code"`
}

type EnergyMeterInfo struct {
	Current      float64 `json:"current"`
	Voltage      float64 `json:"voltage"`
	Power        float64 `json:"power"`
	Total        float64 `json:"total"`
	ErrorCode    int     `json:"err_code"`
	ErrorMessage string  `json:"err_msg"`
}
