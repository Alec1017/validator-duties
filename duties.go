package main

type ValidatorDuties struct {
	Slots map[uint64][]uint64
}

// Create a new instance of validator duties
func New() *ValidatorDuties {
	// Create a new validator duties struct
	validatorDuties := &ValidatorDuties{
		Slots: make(map[uint64][]uint64),
	}

	return validatorDuties
}

// Processes an attester duty response
func (d *ValidatorDuties) ProcessDuties(duties AttesterDutiesResponse) {
	// Pull out the data from the response
	for _, epochData := range duties.Data {
		// At the slot, append the validator index
		d.Slots[epochData.Slot] = append(d.Slots[epochData.Slot], epochData.ValidatorIndex)
	}
}
