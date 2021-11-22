package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AlexGustafsson/tp-link-exporter/internal/tplink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func defaultCommand(context *cli.Context) error {
	verbose := context.Bool("verbose")

	address := context.String("address")
	if address == "" {
		address = ":8080"
	}

	interval := context.Duration("interval")
	if interval == 0 {
		interval = 5 * time.Second
	}

	targets := context.StringSlice("target")
	if targets == nil {
		targets = []string{"192.168.1.255"}
	}

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

	broadcaster := tplink.NewBroadcaster(targets, interval, log)
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

	log.Info("Listening", zap.String("address", address))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(address, nil)
}
