package commit

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MultiSignerWithOverlap tests a multi signer commit msg which does overlap previous valid commit signers
func MultiSignerWithOverlap() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{2, 3}),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "multi signer, with overlap",
		Pre:           pre,
		PostRoot:      "859f2fffb64e1e97f28f656daeed61b97c919eb69cc4f24f6b4831fd8cbde9a3",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		},
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
