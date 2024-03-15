package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

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

func QueryAttesterDuties(epoch uint64, validators []uint64) AttesterDutiesResponse {
	var epochDuties AttesterDutiesResponse

	// Request the attester duties for the epoch
	err := PostRequest(
		fmt.Sprintf("validator/duties/attester/%d", epoch),
		&epochDuties,
		validators,
	)
	if err != nil {
		panic(err)
	}

	return epochDuties
}

func main() {
	// Get the head block of the beacon chain
	headBeaconBlock := QueryBeaconHeadBlock()

	// Get the current slot
	currentSlot := headBeaconBlock.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Validators whose duties should be retrieved
	validators := []uint64{811475}

	// Create a validator duties manager
	DutiesManager := New()

	// Validator duties responses
	currEpochDuties := QueryAttesterDuties(currentEpoch, validators)
	nextEpochDuties := QueryAttesterDuties(currentEpoch+1, validators)

	// Process each response
	DutiesManager.ProcessDuties(currEpochDuties)
	DutiesManager.ProcessDuties(nextEpochDuties)

	// Load the specified timezone. Default to UTC
	loc, _ := time.LoadLocation("America/New_York")

	// The timestamp where the previous attestion ended. At the start, it will just be
	// the current time.
	prevAttestEnd := time.Now().In(loc)

	// For each slot in the mapping, get the timestamp of the slot start
	for s, validators := range DutiesManager.Slots {
		// Determine the timestamp the slot started
		slotStart := time.Unix(int64(BeaconChainGenesis+s*SecondsPerSlot), 0)

		// Gap until the next attestion must be made
		gapUntilNextAttest := int64(math.Floor(slotStart.Sub(prevAttestEnd).Seconds()))

		// If there are validators at this slot, list those validators
		var validatorStrs []string
		for _, v := range validators {
			validatorStrs = append(validatorStrs, strconv.Itoa(int(v)))
		}

		// If the attestation has already occurred in the current epoch, then it
		// can be skipped
		if gapUntilNextAttest < 0 {
			// Display output if already attested
			fmt.Printf(
				"slot %d - already attested at this epoch - Validators: %s\n",
				s,
				strings.Join(validatorStrs, ", "),
			)
		} else {
			// Display output if havent yet attested
			fmt.Printf(
				"slot %d - gap of %d seconds - from %s to %s - Validators: %s\n",
				s,
				gapUntilNextAttest,
				prevAttestEnd.Format(time.Kitchen),
				slotStart.Format(time.Kitchen),
				strings.Join(validatorStrs, ", "),
			)
		}

		// Set the previous attestion time
		prevAttestEnd = slotStart
	}

	// TODO: I should add the epoch in the output as well
	// TODO: I should make the validator the key in the mapping. Then, for each validator
	// we can loop through and calculate the unique gaps for each one.
}
