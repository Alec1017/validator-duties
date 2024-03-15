package main

import (
	"fmt"
)

func main() {
	// Signed beacon block response
	var signedBeaconBlock SignedBeaconBlockResponse

	// Request the head block of the beacon chain
	err := GetRequest("beacon/headers/head", &signedBeaconBlock)
	if err != nil {
		panic(err)
	}

	// Get the current slot
	currentSlot := signedBeaconBlock.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Validators whose duties should be retrieved
	validators := []uint64{811475}

	// Validator duties responses
	var currEpochDuties AttesterDutiesResponse
	var nextEpochDuties AttesterDutiesResponse

	// Request the attester duties for the current epoch
	currEpochErr := PostRequest(
		fmt.Sprintf("validator/duties/attester/%d", currentEpoch),
		&currEpochDuties,
		validators,
	)
	if currEpochErr != nil {
		panic(currEpochErr)
	}

	// Request the attester duties for the next epoch
	nextEpochErr := PostRequest(
		fmt.Sprintf("validator/duties/attester/%d", currentEpoch+1),
		&nextEpochDuties,
		validators,
	)
	if nextEpochErr != nil {
		panic(nextEpochErr)
	}

	// Using both epoch responses, we want to build up a mapping that contains
	// all the slots from the two responses as keys. Each key will contain a list of
	// validator indices who are attesting at that slot.

	// Then, at the end, assume that all validators will be attesting in currEpoch + 2 at
	// the very first slot. This is just so we assume the worst case.

	// Print it out
	fmt.Printf("%+v", signedBeaconBlock)
	fmt.Println()
	fmt.Printf("%+v", currEpochDuties)
}
