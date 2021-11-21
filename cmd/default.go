package main

import (
	"net/http"

	"github.com/AlexGustafsson/tp-link-exporter/internal/tplink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func defaultCommand(context *cli.Context) error {
	broadcaster := tplink.NewBroadcaster()
	broadcaster.BroadcastAddress = context.String("address") + ":9999"
	log.Infof("Finding devices in %s:9999", context.String("address"))
	go broadcaster.Listen()

	collector := tplink.NewCollector()
	if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
		return err
	}

	go func() {
		for {
			device := <-broadcaster.Responses()
			collector.CollectDevice(device)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":2112", nil)
}
