package duties

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

// Contains the duties of a validator for a given epoch
type ValidatorDutyEpoch struct {
	AttestationSlot uint64
	ProposalSlot    uint64
	Epoch           uint64
}

type ValidatorDuties struct {
	Timezone           *time.Location
	BeaconNodeEndpoint string
	CurrEpochDuties    *ValidatorDutyEpoch
	NextEpochDuties    *ValidatorDutyEpoch
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
func (d *ValidatorDuties) ProcessAttesterDuties(epoch uint64) (uint64, error) {
	// Query the attester duties from the beacon node
	attesterDuties, err := d.QueryAttesterDuties(epoch)
	if err != nil {
		return 0, err
	}

	// Ensure that the data response is only of size 1
	if len(attesterDuties.Data) != 1 {
		return 0, errors.New("should have received only 1 attester duty response")
	}

	// Pull out the data from the response
	attestationSlot := attesterDuties.Data[0].Slot

	return attestationSlot, nil
}

// Processes a proposer duty response
func (d *ValidatorDuties) ProcessProposerDuties(epoch uint64) (uint64, error) {
	// Query the proposer duties from the beacon node
	proposerDuties, err := d.QueryProposerDuties(epoch)
	if err != nil {
		return 0, err
	}

	// Pull out the data from the response
	for _, proposerData := range proposerDuties.Data {
		// Only pull out the slot if ithe validator index matches
		if d.Validator == proposerData.ValidatorIndex {
			proposalSlot := proposerDuties.Data[0].Slot

			return proposalSlot, nil
		}
	}

	return 0, nil
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

// Queries the duties of the proposer for a given epoch
func (d *ValidatorDuties) QueryProposerDuties(epoch uint64) (*ProposerDutiesResponse, error) {
	// Proposer duties response
	var proposerDuties ProposerDutiesResponse

	// Request the proposer duties for the epoch
	err := GetRequest(
		fmt.Sprintf(d.BeaconNodeEndpoint+"/eth/v1/validator/duties/proposer/%d", epoch),
		&proposerDuties,
	)
	if err != nil {
		return nil, err
	}

	return &proposerDuties, nil
}

func (d *ValidatorDuties) DisplayDuties() {
	// The timestamp where the previous attestion ended. At the start, it will just be
	// the current time.
	prevAttestEnd := time.Now().In(d.Timezone)

	// Display the validator
	fmt.Printf("Validator: %s\n", strconv.Itoa(int(d.Validator)))

	// For each slot in the mapping, get the timestamp of the slot start
	for _, epochDuties := range []*ValidatorDutyEpoch{d.CurrEpochDuties, d.NextEpochDuties} {
		// Determine the timestamp the attestation slot started
		slotStart := time.Unix(
			int64(BeaconChainGenesis+epochDuties.AttestationSlot*SecondsPerSlot),
			0,
		)

		// Gap until the next attestion must be made
		gapUntilNextAttest := int64(math.Floor(slotStart.Sub(prevAttestEnd).Seconds()))

		// If the attestation has already occurred in the current epoch, then it
		// can be skipped
		if gapUntilNextAttest < 0 {
			// Display output if already attested
			fmt.Printf(
				"epoch %d - slot %d - already attested at this epoch\n",
				epochDuties.Epoch,
				epochDuties.AttestationSlot,
			)
		} else {
			// Display output if havent yet attested
			fmt.Printf(
				"epoch %d - slot %d - gap of %d seconds - from %s to %s\n",
				epochDuties.Epoch,
				epochDuties.AttestationSlot,
				gapUntilNextAttest,
				prevAttestEnd.Format(time.Kitchen),
				slotStart.Format(time.Kitchen),
			)
		}

		// Set the previous attestion time
		prevAttestEnd = slotStart
	}
}

func Start(opts ...Option) error {
	// Create a validator duties manager with the validator whose duties should
	// be retrieved
	d, err := New(opts...)
	if err != nil {
		return err
	}

	// Get the head block of the beacon chain
	headBeaconBlock, err := d.QueryBeaconHeadBlock()
	if err != nil {
		return err
	}

	// Get the current slot
	currentSlot := headBeaconBlock.Data.Header.Message.Slot

	// Determine current epoch using current slot
	currentEpoch := currentSlot / SlotsPerEpoch

	// Process the attester duties for the current epoch
	currEpochAttesterDutySlot, err := d.ProcessAttesterDuties(currentEpoch)
	if err != nil {
		return err
	}

	// Process the attester duties for the next epoch
	nextEpochAttesterDutySlot, err := d.ProcessAttesterDuties(currentEpoch + 1)
	if err != nil {
		return err
	}

	// Process the proposer duties for the current epoch
	currEpochProposerDutySlot, err := d.ProcessProposerDuties(currentEpoch)
	if err != nil {
		return err
	}

	// Process the proposer duties for the next epoch
	nextEpochProposerDutySlot, err := d.ProcessProposerDuties(currentEpoch + 1)
	if err != nil {
		return err
	}

	// Set the current epoch duties
	d.CurrEpochDuties = &ValidatorDutyEpoch{
		Epoch:           currentEpoch,
		AttestationSlot: currEpochAttesterDutySlot,
		ProposalSlot:    currEpochProposerDutySlot,
	}

	// Set the next epoch duties
	d.NextEpochDuties = &ValidatorDutyEpoch{
		Epoch:           currentEpoch + 1,
		AttestationSlot: nextEpochAttesterDutySlot,
		ProposalSlot:    nextEpochProposerDutySlot,
	}

	// Display the output
	d.DisplayDuties()

	return nil
}
