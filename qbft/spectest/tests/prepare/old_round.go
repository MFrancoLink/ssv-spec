package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// OldRound tests prepare for signedProposal.Message.Round < state.Round
func OldRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 10

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 10),
	}
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications(rcMsgs), nil)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], 1, 9),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare prev round",
		Pre:           pre,
		PostRoot:      "f16c3384489a8d2ee68e57e8fb5194498c90f694a142a6dda397cc192fd7bdac",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
