package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AlexGustafsson/tp-link-exporter/internal/tplink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func defaultCommand(context *cli.Context) error {
	address := context.String("address")
	verbose := context.Bool("verbose")

	// Configure base logging
	logConfig := zap.NewProductionConfig()
	if verbose {
		logConfig.Level.SetLevel(zap.DebugLevel)
	}
	log, err := logConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logging: %v", err)
		os.Exit(1)
	}
	defer log.Sync()

	broadcaster := tplink.NewBroadcaster(log)
	broadcaster.BroadcastAddress = address + ":9999"
	log.Info("Finding devices", zap.String("address", address), zap.Int("port", 9999))
	go broadcaster.Listen()

	collector := tplink.NewCollector()
	if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
		return err
	}

	go func() {
		for {
			device := <-broadcaster.Responses()
			log.Debug("Found device", zap.String("name", device.Device.Info.Alias))
			collector.CollectDevice(device)
		}
	}()

	log.Info("Listening", zap.String("address", ":8080"))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8080", nil)
}
