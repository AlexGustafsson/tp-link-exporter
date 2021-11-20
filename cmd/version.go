package main

import (
	"fmt"

	"github.com/AlexGustafsson/tp-link-exporter/internal/version"
	"github.com/urfave/cli/v2"
)

func versionCommand(context *cli.Context) error {
	return printVersion()
}

func printVersion() error {
	fmt.Println(version.FullVersion())

	return nil
}
