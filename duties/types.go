package duties

/////////////////////////////////////////////////////////////////////////////////////////
//                                 Beacon Block Headers                                //
/////////////////////////////////////////////////////////////////////////////////////////

type BeaconBlockHeader struct {
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
	Slot          uint64 `json:"slot,string"`
}

type SignedBeaconBlockHeader struct {
	Signature string            `json:"signature"`
	Message   BeaconBlockHeader `json:"message"`
}

type SignedBeaconBlockHeaderContainer struct {
	Root      string                  `json:"root"`
	Header    SignedBeaconBlockHeader `json:"header"`
	Canonical bool                    `json:"canonical"`
}

type SignedBeaconBlockResponse struct {
	Data                SignedBeaconBlockHeaderContainer `json:"data"`
	ExecutionOptimistic bool                             `json:"execution_optimistic"`
	Finalized           bool                             `json:"finalized"`
}

/////////////////////////////////////////////////////////////////////////////////////////
//									Attester Duties									   //
/////////////////////////////////////////////////////////////////////////////////////////

type AttesterDuty struct {
	Pubkey                 string `json:"pubkey"`
	CommitteeIndex         string `json:"committee_index"`
	CommitteeLength        string `json:"committee_length"`
	CommitteesAtSlot       string `json:"committees_at_slot"`
	ValidatorCommiteeIndex string `json:"validator_committee_index"`
	ValidatorIndex         uint64 `json:"validator_index,string"`
	Slot                   uint64 `json:"slot,string"`
}

type AttesterDutiesResponse struct {
	DependentRoot       string          `json:"dependent_root"`
	Data                []*AttesterDuty `json:"data"`
	ExecutionOptimistic bool            `json:"execution_optimistic"`
}

/////////////////////////////////////////////////////////////////////////////////////////
//									Proposer Duties									   //
/////////////////////////////////////////////////////////////////////////////////////////

type ProposerDuty struct {
	Pubkey         string `json:"pubkey"`
	ValidatorIndex uint64 `json:"validator_index,string"`
	Slot           uint64 `json:"slot,string"`
}

type ProposerDutiesResponse struct {
	DependentRoot       string          `json:"dependent_root"`
	Data                []*ProposerDuty `json:"data"`
	ExecutionOptimistic bool            `json:"execution_optimistic"`
}
