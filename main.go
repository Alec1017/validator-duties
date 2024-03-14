package main

import (
	"fmt"
)

func main() {
	var b SignedBeaconBlockResponse
	var d AttesterDutiesResponse

	// Request the head block of the beacon chain
	err := GetRequest("beacon/headers/head", &b)
	if err != nil {
		panic(err)
	}

	// Get the current slot
	currentSlot := b.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Validators whose duties should be retrieved
	validators := []uint64{811475}

	// Request the attester duties for the current epoch
	postErr := PostRequest(
		fmt.Sprintf("validator/duties/attester/%d", currentEpoch),
		&d,
		validators,
	)
	if postErr != nil {
		panic(postErr)
	}

	// Print it out
	fmt.Printf("%+v", b)
	fmt.Println()
	fmt.Printf("%+v", d)
}
