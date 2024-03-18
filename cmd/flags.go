package cmd

import (
	"time"

	"github.com/Alec1017/validator-duties/duties"
	"github.com/urfave/cli/v2"
)

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

// Pulls the CLI options from the context and converts them
// into executable option functions to be processed by the
// validator duties manager
func FlagOptions(c *cli.Context) ([]duties.Option, error) {
	// parse the validator
	validator := c.Uint64(Validator.Name)

	// parse the timezone
	timezoneStr := c.String(Timezone.Name)

	// parse the beacon node endpoint
	endpoint := c.String(BeaconNodeEndpoint.Name)

	// Load the specified timezone. Default to UTC
	timezone, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, err
	}

	// Create an array of options
	opts := []duties.Option{
		duties.WithValidator(validator),
		duties.WithTimezone(timezone),
		duties.WithBeaconNodeEndpoint(endpoint),
	}

	return opts, nil
}
