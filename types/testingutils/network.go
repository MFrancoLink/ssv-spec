package testingutils

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
	OperatorID      types.OperatorID
	OperatorSK      *rsa.PrivateKey
	Domain          types.DomainType
}

func NewTestingNetwork(operatorID types.OperatorID, sk *rsa.PrivateKey) *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
		OperatorID:      operatorID,
		OperatorSK:      sk,
		Domain:          TestingSSVDomainType,
	}
}

func (net *TestingNetwork) Broadcast(msgID types.MessageID, message *types.SignedSSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func ConvertBroadcastedMessagesToSSVMessages(signedMessages []*types.SignedSSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range signedMessages {
		ret = append(ret, msg.SSVMessage)
	}
	return ret
}

// GetDomainType returns the Domain type used for signatures
func (c *TestingNetwork) GetDomainType() types.DomainType {
	return c.Domain
}
