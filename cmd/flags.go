package cmd

import "github.com/urfave/cli/v2"

var (
	Validator = &cli.IntFlag{
		Name:  "validator",
		Usage: "Index of the validator to check duties",
	}

	Timezone = &cli.StringFlag{
		Name:  "timezone",
		Usage: "Timezone for the attestation duties",
		Value: "UTC",
	}
)

var Flags = []cli.Flag{
	Validator,
	Timezone,
}
