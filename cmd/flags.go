package cmd

import "github.com/urfave/cli/v2"

var (
	Validator = &cli.Uint64Flag{
		Name:     "validator",
		Usage:    "Index of the validator to check duties",
		Required: true,
	}

	Timezone = &cli.StringFlag{
		Name:  "timezone",
		Usage: "Timezone for the attestation duties",
		Value: "UTC",
	}

	BeaconNodeEndpoint = &cli.StringFlag{
		Name:  "beacon-node-endpoint",
		Usage: "A consensus client http endpoint",
		Value: "http://localhost:5052",
	}
)

var Flags = []cli.Flag{
	Validator,
	Timezone,
	BeaconNodeEndpoint,
}
