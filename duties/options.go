package duties

import (
	"time"

	"github.com/Alec1017/validator-duties/cmd"
	"github.com/urfave/cli/v2"
)

// Option type that will operate on the validator duties and
// modify its state
type Option func(d *ValidatorDuties) error

func WithValidator(validator uint64) Option {
	return func(d *ValidatorDuties) error {
		// set the validator
		d.Validator = validator

		return nil
	}
}

func WithTimezone(timezone *time.Location) Option {
	return func(d *ValidatorDuties) error {
		// set the timezone
		d.Timezone = timezone

		return nil
	}
}

func WithBeaconNodeEndpoint(endpoint string) Option {
	return func(d *ValidatorDuties) error {
		// set the beacon node endpoint
		d.BeaconNodeEndpoint = endpoint

		return nil
	}
}

// Pulls the CLI options from the context and converts them
// into executable option functions to be processed by the
// validator duties manager
func FlagOptions(c *cli.Context) ([]Option, error) {
	// parse the validator
	validator := c.Uint64(cmd.Validator.Name)

	// parse the timezone
	timezoneStr := c.String(cmd.Timezone.Name)

	// parse the beacon node endpoint
	endpoint := c.String(cmd.BeaconNodeEndpoint.Name)

	// Load the specified timezone. Default to UTC
	timezone, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, err
	}

	// Create an array of options
	opts := []Option{
		WithValidator(validator),
		WithTimezone(timezone),
		WithBeaconNodeEndpoint(endpoint),
	}

	return opts, nil
}
