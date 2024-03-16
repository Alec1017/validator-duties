package duties

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

// Simple wrapper struct for an epoch/slot pair
type EpochSlot struct {
	Slot  uint64
	Epoch uint64
}

type ValidatorDuties struct {
	EpochSlots []EpochSlot
	Validator  uint64
}

// Create a new instance of validator duties
func New(validator uint64) *ValidatorDuties {
	// Create a new validator duties struct
	validatorDuties := &ValidatorDuties{
		Validator:  validator,
		EpochSlots: []EpochSlot{},
	}

	return validatorDuties
}

// Processes an attester duty response
func (d *ValidatorDuties) ProcessDuties(epoch uint64) {
	// Query the attester duties from the beacon node
	epochDuties := QueryAttesterDuties(epoch, d.Validator)

	// Pull out the data from the response
	for _, epochData := range epochDuties.Data {
		// add the slot for the validator
		d.EpochSlots = append(d.EpochSlots, EpochSlot{epochData.Slot, epoch})
	}
}

// Queries the head block of the beacon chain
func QueryBeaconHeadBlock() SignedBeaconBlockResponse {
	// Signed beacon block response
	var signedBeaconBlock SignedBeaconBlockResponse

	// Request the head block of the beacon chain
	err := GetRequest("beacon/headers/head", &signedBeaconBlock)
	if err != nil {
		panic(err)
	}

	return signedBeaconBlock
}

// Queries the duties of the attester for a given epoch
func QueryAttesterDuties(epoch uint64, validator uint64) AttesterDutiesResponse {
	// Attester duties response
	var epochDuties AttesterDutiesResponse

	// Request the attester duties for the epoch
	err := PostRequest(
		fmt.Sprintf("validator/duties/attester/%d", epoch),
		&epochDuties,
		[]uint64{validator},
	)
	if err != nil {
		panic(err)
	}

	return epochDuties
}

func Start(*cli.Context) error {
	// Get the head block of the beacon chain
	headBeaconBlock := QueryBeaconHeadBlock()

	// Get the current slot
	currentSlot := headBeaconBlock.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Create a validator duties manager with the validator whose duties should
	// be retrieved
	dutiesManager := New(811475)

	// Process the duties for the current epoch and the next epoch
	dutiesManager.ProcessDuties(currentEpoch)
	dutiesManager.ProcessDuties(currentEpoch + 1)

	// Load the specified timezone. Default to UTC
	loc, _ := time.LoadLocation("America/New_York")

	// The timestamp where the previous attestion ended. At the start, it will just be
	// the current time.
	prevAttestEnd := time.Now().In(loc)

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
