package main

/////////////////////////////////////////////////////////////////////////////////////////
//                                 Beacon Block Headers                                //
/////////////////////////////////////////////////////////////////////////////////////////

type BeaconBlockHeader struct {
	Slot          uint64 `json:"slot,string"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader `json:"message"`
	Signature string            `json:"signature"`
}

type SignedBeaconBlockHeaderContainer struct {
	Root      string                  `json:"root"`
	Canonical bool                    `json:"canonical"`
	Header    SignedBeaconBlockHeader `json:"header"`
}

type SignedBeaconBlockResponse struct {
	ExecutionOptimistic bool                             `json:"execution_optimistic"`
	Finalized           bool                             `json:"finalized"`
	Data                SignedBeaconBlockHeaderContainer `json:"data"`
}

/////////////////////////////////////////////////////////////////////////////////////////
//									Attester Duties									   //
/////////////////////////////////////////////////////////////////////////////////////////

type AttesterDuty struct {
	Pubkey                 string `json:"pubkey"`
	ValidatorIndex         string `json:"validator_index"`
	CommitteeIndex         string `json:"committee_index"`
	CommitteeLength        string `json:"committee_length"`
	CommitteesAtSlot       string `json:"committees_at_slot"`
	ValidatorCommiteeIndex string `json:"validator_committee_index"`
	Slot                   string `json:"slot"`
}

type AttesterDutiesResponse struct {
	DependentRoot       string         `json:"dependent_root"`
	ExecutionOptimistic bool           `json:"execution_optimistic"`
	Data                []AttesterDuty `json:"data"`
}
