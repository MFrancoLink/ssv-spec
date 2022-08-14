package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests unknown SignedPostConsensusMessage signer
func UnknownSigner() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(testingutils.Testing7SharesSet().Shares[5], 5)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao unknown signer",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid randao message: failed to verify PartialSignature: signer not found in operators",
	}
}