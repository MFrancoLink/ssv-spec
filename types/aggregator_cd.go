package types

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type AssignedAggregator struct {
	ValidatorIndex phase0.ValidatorIndex
	SelectionProof phase0.BLSSignature
	CommitteeIndex uint64
}

func (a *AssignedAggregator) Encode() ([]byte, error) {
	return a.MarshalSSZ()
}
func (a *AssignedAggregator) Decode(data []byte) error {
	return a.UnmarshalSSZ(data)
}

type AggregatorConsensusData struct {
	Version spec.DataVersion // Beacon version

	// Aggregator duties
	Aggregators []AssignedAggregator `ssz-max:"3000"`

	CommitteeIndex []uint64 `ssz-max:"64"`         // ordered unique list of the existing beacon committees, i.e. a subset of the [1,...,64] list
	Attestation    [][]byte `ssz-max:"64,1000000"` // encoded phase0.Attestation or electra.Attestation (depending on the version), one for each beacon committee index

	// Sync Committee Duties
	Contributors []AssignedAggregator `ssz-max:"12000"`

	SubCommitteeIndexes         []uint64                           `ssz-max:"4"` // ordered unique list of the existing sync committee subnets, i.e. a subset of the [1,2,3,4] list
	SyncCommitteeSelectionProof []altair.SyncCommitteeContribution `ssz-max:"4"` // one for each sync committee subnet
}

func (a *AggregatorConsensusData) Encode() ([]byte, error) {
	return a.MarshalSSZ()
}
func (a *AggregatorConsensusData) Decode(data []byte) error {
	return a.UnmarshalSSZ(data)
}
