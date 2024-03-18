package duties

import (
	"time"
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
