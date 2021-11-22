package tplink

import "github.com/prometheus/client_golang/prometheus"

type Collector struct {
	CurrentGauge    *prometheus.GaugeVec
	VoltageGauge    *prometheus.GaugeVec
	PowerGauge      *prometheus.GaugeVec
	RelayStateGauge *prometheus.GaugeVec
	RSSIGauge       *prometheus.GaugeVec
}

func NewCollector() *Collector {
	labels := []string{"name", "device_id", "model", "type"}

	return &Collector{
		CurrentGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "tplink",
				Subsystem: "energy",
				Name:      "current",
				Help:      "Current",
			},
			labels,
		),
		VoltageGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "tplink",
				Subsystem: "energy",
				Name:      "voltage",
				Help:      "Voltage",
			},
			labels,
		),
		PowerGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "tplink",
				Subsystem: "energy",
				Name:      "power_watts",
				Help:      "Power draw in watts",
			},
			labels,
		),
		RelayStateGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "tplink",
				Subsystem: "relay",
				Name:      "state",
				Help:      "State of the relay. 1 is on, 0 is off",
			},
			labels,
		),
		RSSIGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "tplink",
				Subsystem: "statistics",
				Name:      "rssi",
				Help:      "Received Signal Strength Indication",
			},
			labels,
		),
	}
}

func (c *Collector) CollectDevice(device *DeviceResponse) {
	if device.EnergyMeter.Info != nil && device.EnergyMeter.Info.ErrorCode == 0 {
		c.CurrentGauge.WithLabelValues(device.Device.Info.Alias, device.Device.Info.DeviceID, device.Device.Info.Model, device.Device.Info.Type).Set(device.EnergyMeter.Info.Current)
		c.VoltageGauge.WithLabelValues(device.Device.Info.Alias, device.Device.Info.DeviceID, device.Device.Info.Model, device.Device.Info.Type).Set(device.EnergyMeter.Info.Voltage)
		c.PowerGauge.WithLabelValues(device.Device.Info.Alias, device.Device.Info.DeviceID, device.Device.Info.Model, device.Device.Info.Type).Set(device.EnergyMeter.Info.Power)
	}

	c.RelayStateGauge.WithLabelValues(device.Device.Info.Alias, device.Device.Info.DeviceID, device.Device.Info.Model, device.Device.Info.Type).Set(float64(device.Device.Info.RelayState))
	c.RSSIGauge.WithLabelValues(device.Device.Info.Alias, device.Device.Info.DeviceID, device.Device.Info.Model, device.Device.Info.Type).Set(device.Device.Info.RSSI)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.CurrentGauge.Collect(ch)
	c.VoltageGauge.Collect(ch)
	c.PowerGauge.Collect(ch)
	c.RelayStateGauge.Collect(ch)
	c.RSSIGauge.Collect(ch)
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}
