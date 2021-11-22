package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AlexGustafsson/tp-link-exporter/internal/tplink"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func defaultCommand(ctx *cli.Context) error {
	verbose := ctx.Bool("verbose")

	address := ctx.String("address")
	if address == "" {
		address = ":8080"
	}

	interval := ctx.Duration("interval")
	if interval == 0 {
		interval = 5 * time.Second
	}

	targets := ctx.StringSlice("target")
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

	broadcaster := tplink.NewBroadcaster(targets, log)
	log.Info("Finding devices", zap.String("address", address), zap.Int("port", 9999))
	go broadcaster.Listen()

	collector := tplink.NewCollector()
	if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
		return err
	}

	log.Info("Listening", zap.String("address", address))
	handler := promhttp.Handler()
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := broadcaster.Broadcast(); err != nil {
			log.Error("Failed to send broadcast", zap.Error(err))
		}

		// Get the preferred timeout or default to 2s
		scrapeTimeoutHeader := r.Header.Get("X-Prometheus-Scrape-Timeout-Seconds")
		if scrapeTimeoutHeader == "" {
			scrapeTimeoutHeader = "2"
		}
		scrapeTimeoutSeconds, err := strconv.ParseInt(scrapeTimeoutHeader, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Don't allocate more than 2 seconds - that's highly unlikely in a well-functioning network
		scrapeTimeoutSeconds = int64(math.Min(float64(scrapeTimeoutSeconds), float64(2)))
		scrapeTimeout := time.Duration(scrapeTimeoutSeconds) * time.Second

		// Allocate some time for creating the response
		scrapeTimeout -= 200 * time.Millisecond

		// Wait for devices to responds
		log.Debug("Waiting for devices to respond", zap.Duration("timeout", scrapeTimeout))
		devices := make([]*tplink.DeviceResponse, 0)
		timeout, cancel := context.WithTimeout(context.Background(), scrapeTimeout)
	loop:
		for {
			select {
			case device := <-broadcaster.Responses():
				devices = append(devices, device)
				log.Debug("Found device", zap.String("name", device.Device.Info.Alias))
			case <-timeout.Done():
				break loop
			}
		}
		cancel()

		// Reset metrics and update them to remove no longer found, or reconfigured devices
		collector.Reset()
		for _, device := range devices {
			collector.CollectDevice(device)
		}

		handler.ServeHTTP(w, r)
	}))
	return http.ListenAndServe(address, nil)
}
