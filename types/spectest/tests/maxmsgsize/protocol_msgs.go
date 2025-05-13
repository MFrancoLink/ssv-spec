package maxmsgsize

import (
	"math"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

func testingRsaSig() []byte {
	s := [256]byte{1}
	return s[:]
}
func testingIdentifier() []byte {
	s := [56]byte{1}
	return s[:]
}

func signedQbftMsg(qbftMsg *qbft.Message) *types.SignedSSVMessage {
	encodedQbftMsg, err := qbftMsg.Encode()
	if err != nil {
		panic(err)
	}
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   [56]byte{1},
		Data:    encodedQbftMsg,
	}
	signedSSVMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{testingRsaSig()},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,
	}
	return signedSSVMsg
}

func decided(qbftMsg *qbft.Message, count int) *types.SignedSSVMessage {
	encodedQbftMsg, err := qbftMsg.Encode()
	if err != nil {
		panic(err)
	}
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   [56]byte{1},
		Data:    encodedQbftMsg,
	}
	operatorIDs := make([]types.OperatorID, 0)
	signatures := make([][]byte, 0)
	for i := 0; i < count; i++ {
		operatorIDs = append(operatorIDs, types.OperatorID(i))
		signature := [256]byte{1}
		signatures = append(signatures, signature[:])
	}

	signedSSVMsg := &types.SignedSSVMessage{
		SSVMessage:  ssvMsg,
		Signatures:  signatures,
		OperatorIDs: operatorIDs,
	}
	return signedSSVMsg
}

func signedPSigMsg(pSigMsgs *types.PartialSignatureMessages) *types.SignedSSVMessage {
	encodedPSigsMsgs, err := pSigMsgs.Encode()
	if err != nil {
		panic(err)
	}
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   [56]byte{1},
		Data:    encodedPSigsMsgs,
	}
	signedSSVMsg := &types.SignedSSVMessage{
		Signatures:  [][]byte{testingRsaSig()},
		OperatorIDs: []types.OperatorID{1},
		SSVMessage:  ssvMsg,
	}
	return signedSSVMsg
}

func prepareMsg() *qbft.Message {
	return &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Identifier: testingIdentifier(),
		Round:      1,
		Height:     1,
		Root:       [32]byte{1},
	}
}

func proposalMsg() *qbft.Message {
	return &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Identifier: testingIdentifier(),
		Round:      1,
		Height:     1,
		Root:       [32]byte{1},
	}
}

func PartialSigMsgs(count int) *types.PartialSignatureMessages {

	partialSigMsgs := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     1,
		Messages: make([]*types.PartialSignatureMessage, 0),
	}

	if count > 0 {
		for i := 0; i < count; i++ {
			signature := [96]byte{1}
			pSigMsg := &types.PartialSignatureMessage{
				Signer:           1,
				ValidatorIndex:   1,
				SigningRoot:      [32]byte{1},
				PartialSignature: signature[:],
			}
			partialSigMsgs.Messages = append(partialSigMsgs.Messages, pSigMsg)
		}
	}

	return partialSigMsgs
}

func WithFullData(msg *types.SignedSSVMessage, fullData []byte) *types.SignedSSVMessage {
	msg.FullData = fullData
	return msg
}

func BeaconVote() []byte {
	bv := &types.BeaconVote{
		BlockRoot: [32]byte{1},
		Source: &phase0.Checkpoint{
			Epoch: 1,
			Root:  [32]byte{1},
		},
		Target: &phase0.Checkpoint{
			Epoch: 1,
			Root:  [32]byte{1},
		},
	}
	encodedBeaconVote, err := bv.Encode()
	if err != nil {
		panic(err)
	}
	return encodedBeaconVote
}

func assignedAggregator() types.AssignedAggregator {
	assignedAggregator := types.AssignedAggregator{
		ValidatorIndex: 1,
		SelectionProof: [96]byte{1},
		CommitteeIndex: 1,
	}
	return assignedAggregator
}

func expectedCommitteeIndexes(totalIndexes, numValidators int) int {
	numCommitteeIndexes := float64(totalIndexes) * (1.0 - math.Pow(float64((totalIndexes-1)/totalIndexes), float64(numValidators)))
	return int(numCommitteeIndexes)
}

func AggConsensuData(attAggregators int, numAttestations int, scAggregators int, numSCContributions int) *types.AggregatorConsensusData {

	aggregators := make([]types.AssignedAggregator, 0)
	for i := 0; i < attAggregators; i++ {
		aggregators = append(aggregators, assignedAggregator())
	}
	contributors := make([]types.AssignedAggregator, 0)
	for i := 0; i < scAggregators; i++ {
		contributors = append(contributors, assignedAggregator())
	}

	// Expected number of committee index: 64 * (1 - pow(63/64,attAggregators))
	committeeIndexes := make([]uint64, 0)
	for i := 0; i < int(numAttestations); i++ {
		committeeIndexes = append(committeeIndexes, uint64(i))
	}
	attestations := make([][]byte, 0)
	for i := 0; i < int(numAttestations); i++ {
		att := &electra.Attestation{
			AggregationBits: bitfield.NewBitlist(31250),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           phase0.CommitteeIndex(0),
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature:     [96]byte{1},
			CommitteeBits: bitfield.NewBitvector64(),
		}
		attEncoded, err := att.MarshalSSZ()
		if err != nil {
			panic(err)
		}
		attestations = append(attestations, attEncoded)
	}

	// Expected number of subcommittee index: 4 * (1 - pow(3/4,scAggregators))
	subCommitteeIndexes := make([]uint64, 0)
	for i := 0; i < int(numSCContributions); i++ {
		subCommitteeIndexes = append(subCommitteeIndexes, uint64(i))
	}
	syncCommitteeSelectionProofs := make([]altair.SyncCommitteeContribution, 0)
	for i := 0; i < int(numSCContributions); i++ {
		syncCommitteeSelectionProofs = append(syncCommitteeSelectionProofs, altair.SyncCommitteeContribution{
			Slot:              1,
			BeaconBlockRoot:   [32]byte{1},
			SubcommitteeIndex: uint64(i),
			AggregationBits:   bitfield.NewBitvector128(),
			Signature:         [96]byte{1},
		})
	}

	aggConsensusData := &types.AggregatorConsensusData{
		Version:                     spec.DataVersionPhase0,
		Aggregators:                 aggregators,
		CommitteeIndex:              committeeIndexes,
		Attestation:                 attestations,
		Contributors:                contributors,
		SubCommitteeIndexes:         subCommitteeIndexes,
		SyncCommitteeSelectionProof: syncCommitteeSelectionProofs,
	}
	return aggConsensusData
}

func AggConsensuDataBytes(attAggregators int, numAttestations int, scAggregators int, numSCContributions int) []byte {
	aggConsensusData := AggConsensuData(attAggregators, numAttestations, scAggregators, numSCContributions)
	encodedAggConsensusData, err := aggConsensusData.MarshalSSZ()
	if err != nil {
		panic(err)
	}
	return encodedAggConsensusData
}

func ProposerBlock() []byte {
	block := [92000]byte{1}
	return block[:]
}

func PrepareProtocolMsg() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "prepare protocol msg",
		Object:                signedQbftMsg(prepareMsg()),
		ExpectedEncodedLength: 484,
		IsMaxSize:             false,
	}
}

func DecidedMsg4() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "decided msg 4",
		Object:                decided(prepareMsg(), 4),
		ExpectedEncodedLength: 1288,
		IsMaxSize:             false,
	}
}

func DecidedMsg7() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "decided msg 7",
		Object:                decided(prepareMsg(), 7),
		ExpectedEncodedLength: 2092,
		IsMaxSize:             false,
	}
}

func DecidedMsg10() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "decided msg 10",
		Object:                decided(prepareMsg(), 10),
		ExpectedEncodedLength: 2896,
		IsMaxSize:             false,
	}
}

func DecidedMsg13() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "decided msg 13",
		Object:                decided(prepareMsg(), 13),
		ExpectedEncodedLength: 3700,
		IsMaxSize:             false,
	}
}

func EmptyPSigMsgs() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "empty psig protocol msg",
		Object:                signedPSigMsg(PartialSigMsgs(0)),
		ExpectedEncodedLength: 372,
		IsMaxSize:             false,
	}
}

func PSigMsgsWith1() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "psig protocol msg w 1",
		Object:                signedPSigMsg(PartialSigMsgs(1)),
		ExpectedEncodedLength: 372 + 144,
		IsMaxSize:             false,
	}
}
func PSigMsgsWith2() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "psig protocol msg w 2",
		Object:                signedPSigMsg(PartialSigMsgs(2)),
		ExpectedEncodedLength: 372 + 144*2,
		IsMaxSize:             false,
	}
}
func PSigMsgsWith1000() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "psig protocol msg w 1k",
		Object:                signedPSigMsg(PartialSigMsgs(1000)),
		ExpectedEncodedLength: 372 + 144*1000,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgEmpty() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg empty",
		Object:                signedQbftMsg(proposalMsg()),
		ExpectedEncodedLength: 484,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgCommittee() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg committee",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), BeaconVote()),
		ExpectedEncodedLength: 596,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgProposer() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg proposer",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), ProposerBlock()),
		ExpectedEncodedLength: 92484,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgAggCommittee00() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg committee 00",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), AggConsensuDataBytes(0, 0, 0, 0)),
		ExpectedEncodedLength: 516,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgAggCommittee1Aggregator() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg committee 1 agg",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), AggConsensuDataBytes(1, 0, 0, 0)),
		ExpectedEncodedLength: 628,
		IsMaxSize:             false,
	}
}
func ProposalProtocolMsgAggCommittee1SCAggregator() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg committee 1 sc agg",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), AggConsensuDataBytes(0, 0, 1, 0)),
		ExpectedEncodedLength: 628,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgAggCommittee1Att() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg committee 1 att",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), AggConsensuDataBytes(0, 1, 0, 0)),
		ExpectedEncodedLength: 4671,
		IsMaxSize:             false,
	}
}
func ProposalProtocolMsgAggCommittee1SCC() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg committee 1 sc c",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), AggConsensuDataBytes(0, 0, 0, 1)),
		ExpectedEncodedLength: 684,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgAggregatorOnly() *StructureSizeTest {
	cdData := electra.AggregateAndProof{
		AggregatorIndex: 1,
		SelectionProof:  [96]byte{1},
		Aggregate: &electra.Attestation{
			AggregationBits: bitfield.NewBitlist(31250),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           phase0.CommitteeIndex(0),
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature:     [96]byte{1},
			CommitteeBits: bitfield.NewBitvector64(),
		},
	}
	encodedData, err := cdData.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	cd := types.ValidatorConsensusData{
		Duty: types.ValidatorDuty{
			Type:           types.BNRoleAggregator,
			Slot:           1,
			PubKey:         [48]byte{1},
			ValidatorIndex: 1,
			CommitteeIndex: 1,
		},
		Version: spec.DataVersionPhase0,
		DataSSZ: encodedData,
	}
	cdBytes, err := cd.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	return &StructureSizeTest{
		Name:                  "proposal protocol msg agg only",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), cdBytes),
		ExpectedEncodedLength: 4859,
		IsMaxSize:             false,
	}
}

func ProposalProtocolMsgSCOnly() *StructureSizeTest {
	cdData := make(types.Contributions, 0)
	cdData = append(cdData, &types.Contribution{
		SelectionProofSig: [96]byte{1},
		Contribution: altair.SyncCommitteeContribution{
			Slot:              1,
			BeaconBlockRoot:   [32]byte{1},
			SubcommitteeIndex: 1,
			AggregationBits:   bitfield.NewBitvector128(),
			Signature:         [96]byte{1},
		},
	})
	encodedData, err := cdData.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	cd := types.ValidatorConsensusData{
		Duty: types.ValidatorDuty{
			Type:           types.BNRoleAggregator,
			Slot:           1,
			PubKey:         [48]byte{1},
			ValidatorIndex: 1,
			CommitteeIndex: 1,
		},
		Version: spec.DataVersionPhase0,
		DataSSZ: encodedData,
	}
	cdBytes, err := cd.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	return &StructureSizeTest{
		Name:                  "proposal protocol msg sc only",
		Object:                WithFullData(signedQbftMsg(proposalMsg()), cdBytes),
		ExpectedEncodedLength: 868,
		IsMaxSize:             false,
	}
}
