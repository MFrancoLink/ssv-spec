package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongMsgID tests a SSVMessage ID which doesn't belong to the validator
func WrongMsgID() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgWrongID(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "wrong SSVMessage ID",
		Runner:                  dr,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b",
		ExpectedError:           "Messages invalid: msg ID doesn't match validator ID",
	}
}