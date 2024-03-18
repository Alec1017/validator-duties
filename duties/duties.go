package duties

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// Simple wrapper struct for an epoch/slot pair
type EpochSlot struct {
	Slot  uint64
	Epoch uint64
}

type ValidatorDuties struct {
	Timezone           *time.Location
	BeaconNodeEndpoint string
	EpochSlots         []EpochSlot
	Validator          uint64
}

// Create a new instance of validator duties
func New(opts ...Option) (*ValidatorDuties, error) {
	// Create a new validator duties struct
	validatorDuties := &ValidatorDuties{}

	// Process all options
	for _, opt := range opts {
		if err := opt(validatorDuties); err != nil {
			return nil, err
		}
	}

	return validatorDuties, nil
}

// Processes an attester duty response
func (d *ValidatorDuties) ProcessDuties(epoch uint64) error {
	// Query the attester duties from the beacon node
	epochDuties, err := d.QueryAttesterDuties(epoch)
	if err != nil {
		return err
	}

	// Pull out the data from the response
	for _, epochData := range epochDuties.Data {
		// add the slot for the validator
		d.EpochSlots = append(d.EpochSlots, EpochSlot{epochData.Slot, epoch})
	}

	return nil
}

// Queries the head block of the beacon chain
func (d *ValidatorDuties) QueryBeaconHeadBlock() (*SignedBeaconBlockResponse, error) {
	// Signed beacon block response
	var signedBeaconBlock SignedBeaconBlockResponse

	// Request the head block of the beacon chain
	err := GetRequest(d.BeaconNodeEndpoint+"/eth/v1/beacon/headers/head", &signedBeaconBlock)
	if err != nil {
		return nil, err
	}

	return &signedBeaconBlock, nil
}

// Queries the duties of the attester for a given epoch
func (d *ValidatorDuties) QueryAttesterDuties(epoch uint64) (*AttesterDutiesResponse, error) {
	// Attester duties response
	var epochDuties AttesterDutiesResponse

	// Request the attester duties for the epoch
	err := PostRequest(
		fmt.Sprintf(d.BeaconNodeEndpoint+"/eth/v1/validator/duties/attester/%d", epoch),
		&epochDuties,
		[]uint64{d.Validator},
	)
	if err != nil {
		return nil, err
	}

	return &epochDuties, nil
}

func Start(opts ...Option) error {
	// Create a validator duties manager with the validator whose duties should
	// be retrieved
	dutiesManager, err := New(opts...)
	if err != nil {
		return err
	}

	// Get the head block of the beacon chain
	headBeaconBlock, err := dutiesManager.QueryBeaconHeadBlock()
	if err != nil {
		return err
	}

	// Get the current slot
	currentSlot := headBeaconBlock.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Process the duties for the current epoch and the next epoch
	dutiesManager.ProcessDuties(currentEpoch)
	dutiesManager.ProcessDuties(currentEpoch + 1)

	// The timestamp where the previous attestion ended. At the start, it will just be
	// the current time.
	prevAttestEnd := time.Now().In(dutiesManager.Timezone)

	// Display the validator
	fmt.Printf("Validator: %s\n", strconv.Itoa(int(dutiesManager.Validator)))

	// For each slot in the mapping, get the timestamp of the slot start
	for _, epochSlot := range dutiesManager.EpochSlots {
		// Determine the timestamp the slot started
		slotStart := time.Unix(int64(BeaconChainGenesis+epochSlot.Slot*SecondsPerSlot), 0)

		// Gap until the next attestion must be made
		gapUntilNextAttest := int64(math.Floor(slotStart.Sub(prevAttestEnd).Seconds()))

		// If the attestation has already occurred in the current epoch, then it
		// can be skipped
		if gapUntilNextAttest < 0 {
			// Display output if already attested
			fmt.Printf(
				"epoch %d - slot %d - already attested at this epoch\n",
				epochSlot.Epoch,
				epochSlot.Slot,
			)
		} else {
			// Display output if havent yet attested
			fmt.Printf(
				"epoch %d - slot %d - gap of %d seconds - from %s to %s\n",
				epochSlot.Epoch,
				epochSlot.Slot,
				gapUntilNextAttest,
				prevAttestEnd.Format(time.Kitchen),
				slotStart.Format(time.Kitchen),
			)
		}

		// Set the previous attestion time
		prevAttestEnd = slotStart
	}

	return nil
}
