package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// LateDecidedBiggerQuorum tests processing a decided msg for a just decided instance (with a bigger quorum)
func LateDecidedBiggerQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := lateDecidedBiggerQuorumStateComparison()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		testingutils.TestingIdentifier,
		qbft.FirstHeight,
		ks,
	)
	msgs = append(
		msgs,
		testingutils.TestingCommitMultiSignerMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]},
			[]types.OperatorID{1, 2, 3, 4},
		),
	)
	return &tests.ControllerSpecTest{
		Name: "decide late decided bigger quorum",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
	}
}

func lateDecidedBiggerQuorumStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		testingutils.TestingIdentifier,
		qbft.FirstHeight,
		ks,
	)
	msgs = append(
		msgs,
		testingutils.TestingCommitMultiSignerMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]},
			[]types.OperatorID{1, 2, 3, 4},
		),
	)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
		testingutils.TestingOperatorSigner(ks),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ID:              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(
				testingutils.TestingProposalMessage(
					ks.OperatorKeys[1],
					types.OperatorID(1),
				),
			),
			LastPreparedRound: qbft.FirstRound,
			LastPreparedValue: testingutils.TestingQBFTFullData,
			Decided:           true,
			DecidedValue:      testingutils.TestingQBFTFullData,
			Round:             qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
