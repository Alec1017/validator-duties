package main

import "fmt"

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
