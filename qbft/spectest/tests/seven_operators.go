package tests

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a simple full happy flow until decided
func SevenOperators() SpecTest {
	pre := testingutils.SevenOperatorsInstance()
	ks := testingutils.Testing7SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[5], types.OperatorID(5)),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[5], types.OperatorID(5)),
	}
	return &MsgProcessingSpecTest{
		Name:          "happy flow seven operators",
		Pre:           pre,
		PostRoot:      "eceb1d9689217606ddef19d2c73565006892eef1bb25697312184d9bc6909301",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
	}
}
