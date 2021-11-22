package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/AlexGustafsson/tp-link-exporter/internal/version"
	"github.com/urfave/cli/v2"
)

var appHelpTemplate = `Usage: {{.Name}} [global options] command [command options] [arguments]

{{.Usage}}

Version: {{.Version}}

Options:
  {{range .Flags}}{{.}}
  {{end}}
Commands:
  {{range .Commands}}{{.Name}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} help command' for more information on a command.
`

var commandHelpTemplate = `Usage: tp-link-exporter {{.Name}} [options] {{if .ArgsUsage}}{{.ArgsUsage}}{{end}}

{{.Usage}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .Flags}}

Options:{{range .Flags}}
   {{.}}{{end}}{{end}}
`

func commandNotFound(context *cli.Context, command string) {
	fmt.Fprintf(os.Stderr,
		"%s: '%s' is not a %s command. See '%s help'.",
		context.App.Name,
		command,
		context.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}

func main() {
	cli.AppHelpTemplate = appHelpTemplate
	cli.CommandHelpTemplate = commandHelpTemplate

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "A prometheus exporter for TP-Link smart home devices"
	app.Version = version.FullVersion()
	app.CommandNotFound = commandNotFound
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "version",
			Usage:  "Show the application's version",
			Action: versionCommand,
		},
	}
	app.Action = defaultCommand
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Enable verbose logging",
		},
		&cli.StringFlag{
			Name:        "address",
			Usage:       "Address to serve metrics on",
			DefaultText: ":8080",
		},
		&cli.StringSliceFlag{
			Name:        "target",
			Usage:       "Target address to talk to. May be specified multiple times. May be a broadcast address.",
			DefaultText: "192.168.1.255",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
