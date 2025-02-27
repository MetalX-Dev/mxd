package main

import (
	"log"

	"github.com/MetalX-Dev/mxd/controller"
	"github.com/alecthomas/kong"
)

var cli struct {
	// Run the controller server
	Start struct {
		Port int `help:"Port to listen for API requests on" default:"8080" short:"p"`
	} `cmd:"" help:"Start the controller server"`
}

func runStart(port int) {
	log.Printf("Starting controller on port %d", port)
	controller.StartServer(port)
}

func main() {
	ctx := kong.Parse(&cli, kong.Name("mxd"), kong.Description("MetalX Controller"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{
		Compact: true,
		Summary: true,
	}))
	switch ctx.Command() {
	case "start":
		runStart(cli.Start.Port)
	default:
		ctx.FatalIfErrorf(ctx.PrintUsage(true))
	}
}
