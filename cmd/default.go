package main

import (
	"encoding/json"
	"os"

	"github.com/AlexGustafsson/tp-link-exporter/tplink"
	"github.com/urfave/cli/v2"
)

func defaultCommand(context *cli.Context) error {
	finder := tplink.NewDeviceFinder()
	finder.BroadcastAddress = "10.0.0.200:9999"

	go finder.Listen()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	for {
		system := <-finder.Found()
		encoder.Encode(system)
	}
}
