package tests

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []struct {
		InputValue    []byte
		InputMessages []*qbft.SignedMessage
		Decided       bool
		DecidedVal    []byte
		DecidedCnt    uint
		SavedDecided  *qbft.SignedMessage
	}
	ValCheck       qbft.ProposedValueCheckF
	OutputMessages []*qbft.SignedMessage
	ExpectedError  string
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		test.ValCheck,
		func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 1
		},
	)

	var lastErr error
	for _, runData := range test.RunInstanceData {
		startedInstance := false
		err := contr.StartNewInstance(runData.InputValue)
		if err != nil {
			lastErr = err
		} else {
			startedInstance = true
		}

		if !startedInstance {
			continue
		}

		decidedCnt := 0
		for _, msg := range runData.InputMessages {
			decided, _, err := contr.ProcessMsg(msg)
			if err != nil {
				lastErr = err
			}
			if decided {
				decidedCnt++
			}
		}

		require.EqualValues(t, runData.DecidedCnt, decidedCnt)

		isDecided, decidedVal := contr.InstanceForHeight(contr.Height).IsDecided()
		require.EqualValues(t, runData.Decided, isDecided)
		require.EqualValues(t, runData.DecidedVal, decidedVal)

		if runData.SavedDecided != nil {
			// test saved to storage
			decided, err := contr.GenerateConfig().GetStorage().GetHighestDecided(identifier[:])
			require.NoError(t, err)
			require.NotNil(t, decided)
			r1, err := decided.GetRoot()
			require.NoError(t, err)

			r2, err := runData.SavedDecided.GetRoot()
			require.NoError(t, err)

			require.EqualValues(t, r2, r1)
			require.EqualValues(t, runData.SavedDecided.Signers, decided.Signers)
			require.EqualValues(t, runData.SavedDecided.Signature, decided.Signature)

			// test broadcasted
			broadcastedMsgs := contr.GenerateConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
			require.Greater(t, len(broadcastedMsgs), 0)
			found := false
			for _, msg := range broadcastedMsgs {
				if msg.MsgType == types.SSVDecidedMsgType && bytes.Equal(identifier[:], msg.MsgID[:]) {
					msg2 := &qbft.DecidedMessage{}
					require.NoError(t, msg2.Decode(msg.Data))
					r1, err := msg2.SignedMessage.GetRoot()
					require.NoError(t, err)

					if bytes.Equal(r1, r2) &&
						reflect.DeepEqual(runData.SavedDecided.Signers, msg2.SignedMessage.Signers) &&
						reflect.DeepEqual(runData.SavedDecided.Signature, msg2.SignedMessage.Signature) {
						require.False(t, found)
						found = true
					}
				}
			}
			require.True(t, found)
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSpecTest) TestName() string {
	return "qbft controller " + test.Name
}