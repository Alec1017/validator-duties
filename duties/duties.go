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
	ProposalSlots   []uint64
	AttestationSlot uint64
	Epoch           uint64
	SyncCommittee   bool
}

type ValidatorDuties struct {
	Timezone           *time.Location
	CurrEpochDuties    *ValidatorDutyEpoch
	NextEpochDuties    *ValidatorDutyEpoch
	BeaconNodeEndpoint string
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
func (d *ValidatorDuties) ProcessProposerDuties(epoch uint64) ([]uint64, error) {
	// Define an empty proposal slots array
	proposalSlots := make([]uint64, 0)

	// Query the proposer duties from the beacon node
	proposerDuties, err := d.QueryProposerDuties(epoch)
	if err != nil {
		return proposalSlots, err
	}

	// Pull out the data from the response
	for _, proposerData := range proposerDuties.Data {
		// Only pull out the slot if ithe validator index matches
		if d.Validator == proposerData.ValidatorIndex {
			// Append the slot
			proposalSlots = append(proposalSlots, proposerDuties.Data[0].Slot)
		}
	}

	return proposalSlots, nil
}

// Processes sync committee duty response
func (d *ValidatorDuties) ProcessSyncCommitteeDuties(epoch uint64) (bool, error) {
	// Query the sync committee duties from the beacon node
	syncCommitteeDuties, err := d.QuerySyncCommitteeDuties(epoch)
	if err != nil {
		return false, err
	}

	// No sync committee duties found for this epoch
	if len(syncCommitteeDuties.Data) == 0 {
		return false, nil
	}

	// Ensure that the data response is only of size 1
	if len(syncCommitteeDuties.Data) != 1 {
		return false, errors.New("should have received only 1 sync committee duty response")
	}

	// Ensure the response is for the correct validator
	if d.Validator != syncCommitteeDuties.Data[0].ValidatorIndex {
		return false, errors.New("received sync committee response for the wrong validator")
	}

	return true, nil
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

// Queries the sync committee duties for a given epoch
func (d *ValidatorDuties) QuerySyncCommitteeDuties(
	epoch uint64,
) (*SyncCommitteeDutiesResponse, error) {
	// Sync committee duties response
	var epochDuties SyncCommitteeDutiesResponse

	// Request the sync committee duties for the epoch
	err := PostRequest(
		fmt.Sprintf(d.BeaconNodeEndpoint+"/eth/v1/validator/duties/sync/%d", epoch),
		&epochDuties,
		[]uint64{d.Validator},
	)
	if err != nil {
		return nil, err
	}

	return &epochDuties, nil
}

func (d *ValidatorDuties) DisplayDuties() {
	// The timestamp where the previous attestion ended. At the start, it will just be
	// the current time.
	prevAttestEnd := time.Now().In(d.Timezone)

	// Display the validator
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Printf(
		"|                          Validator: %s                            |\n",
		strconv.Itoa(int(d.Validator)),
	)
	fmt.Print("-------------------------------------------------------------------------\n\n")
	fmt.Println("----------------------------- Attestations ------------------------------")

	// Loop through the current epoch and the next epoch
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

	fmt.Println("\n------------------------------ Proposals --------------------------------")

	// Loop through the current epoch and the next epoch
	for _, epochDuties := range []*ValidatorDutyEpoch{d.CurrEpochDuties, d.NextEpochDuties} {
		// Display output if proposals were found
		if len(epochDuties.ProposalSlots) > 0 {
			fmt.Printf(
				"epoch %d - WARNING: at least one proposal is schedule in this epoch!\n",
				epochDuties.Epoch,
			)
		} else {
			fmt.Printf(
				"epoch %d - not proposing any blocks in this epoch \n",
				epochDuties.Epoch,
			)
		}
	}

	fmt.Println("\n---------------------------- Sync Committee -----------------------------")

	// Loop through the current epoch and the next epoch
	for _, epochDuties := range []*ValidatorDutyEpoch{d.CurrEpochDuties, d.NextEpochDuties} {
		// Display output if sync committee duties were found
		if epochDuties.SyncCommittee {
			fmt.Printf(
				"epoch %d - WARNING: validator is part of a sync committee!\n",
				epochDuties.Epoch,
			)
		} else {
			fmt.Printf(
				"epoch %d - not part of a sync committee in this epoch \n",
				epochDuties.Epoch,
			)
		}
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
	currEpochProposerDutySlots, err := d.ProcessProposerDuties(currentEpoch)
	if err != nil {
		return err
	}

	// Process the proposer duties for the next epoch
	nextEpochProposerDutySlots, err := d.ProcessProposerDuties(currentEpoch + 1)
	if err != nil {
		return err
	}

	// Process the sync committee duties for the current epoch
	currEpochSyncCommitteeDuty, err := d.ProcessSyncCommitteeDuties(currentEpoch)
	if err != nil {
		return err
	}

	// Process the sync committee duties for the next epoch
	nextEpochSyncCommitteeDuty, err := d.ProcessSyncCommitteeDuties(currentEpoch + 1)
	if err != nil {
		return err
	}

	// Set the current epoch duties
	d.CurrEpochDuties = &ValidatorDutyEpoch{
		Epoch:           currentEpoch,
		AttestationSlot: currEpochAttesterDutySlot,
		ProposalSlots:   currEpochProposerDutySlots,
		SyncCommittee:   currEpochSyncCommitteeDuty,
	}

	// Set the next epoch duties
	d.NextEpochDuties = &ValidatorDutyEpoch{
		Epoch:           currentEpoch + 1,
		AttestationSlot: nextEpochAttesterDutySlot,
		ProposalSlots:   nextEpochProposerDutySlots,
		SyncCommittee:   nextEpochSyncCommitteeDuty,
	}

	// Display the output
	d.DisplayDuties()

	return nil
}
