package main

import (
	"encoding/json"
	"os"

	"github.com/AlexGustafsson/tp-link-exporter/tplink"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func defaultCommand(context *cli.Context) error {
	broadcaster := tplink.NewBroadcaster()
	broadcaster.BroadcastAddress = context.String("address") + ":9999"
	log.Infof("Finding devices in %s:9999", context.String("address"))
	go broadcaster.Listen()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	for {
		device := <-broadcaster.Responses()
		encoder.Encode(device)
	}
}
