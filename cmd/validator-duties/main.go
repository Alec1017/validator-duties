package main

import (
	"log"
	"os"

	"github.com/Alec1017/validator-duties/cmd"
	"github.com/Alec1017/validator-duties/duties"
	"github.com/urfave/cli/v2"
)

func main() {
	// Define the configuration for the CLI
	app := cli.App{}
	app.Name = "validator-duties"
	app.Usage = "CLI that allows you to check the upcoming epochs for gaps between attestations. Useful for maintenance of the validator."
	app.Version = "0.0.1"
	app.Flags = cmd.Flags
	app.Action = func(ctx *cli.Context) error {
		// Execute the app logic
		if err := duties.Start(ctx); err != nil {
			return cli.Exit(err.Error(), 1)
		}

		return nil
	}

	// Kick off the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
