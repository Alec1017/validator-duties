package cmd

import (
	"log"
	"os"

	"github.com/Alec1017/validator-duties/duties"
	"github.com/urfave/cli/v2"
)

func Run() {
	// Define the configuration for the CLI
	app := cli.App{}
	app.Name = "validator-duties"
	app.Usage = "CLI that allows you to check the upcoming epochs for gaps between attestations. Useful for maintenance of the validator."
	app.Version = "0.1.0"
	app.Flags = Flags
	app.Action = validatorDuties
	// Kick off the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func validatorDuties(ctx *cli.Context) error {
	// Extract all flag options
	flagOptions, err := FlagOptions(ctx)
	if err != nil {
		return err
	}

	// Execute the app logic
	if err := duties.Start(flagOptions...); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return nil
}
