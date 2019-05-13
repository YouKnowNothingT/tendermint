package consensus

import (
	"bytes"
	"encoding/gob"
	"fmt"
	dkg "go.dedis.ch/kyber/share/dkg/rabin"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"
)

type DKGMockDontSendOneCommit struct {
	Dealer
}

func NewDKGMockDealerNoCommit(validators *types.ValidatorSet, pubKey crypto.PubKey, sendMsgCb func(*types.DKGData), logger log.Logger) Dealer {
	return &DKGMockDontSendOneCommit{NewDKGDealer(validators, pubKey, sendMsgCb, logger)}
}

func (m *DKGMockDontSendOneCommit) Start() error {
	err := m.Dealer.Start()
	if err != nil {
		return err
	}
	m.GenerateTransitions()
	return nil
}

func (m *DKGMockDontSendOneCommit) GenerateTransitions() {
	m.Dealer.SetTransitions([]transition{
		// Phase I
		m.Dealer.SendDeals,
		m.Dealer.ProcessDeals,
		m.Dealer.ProcessResponses,
		m.ProcessJustifications,
		// Phase II
		m.Dealer.ProcessCommits,
		m.Dealer.ProcessComplaints,
		m.Dealer.ProcessReconstructCommits,
	})
}

func (m *DKGMockDontSendOneCommit) ProcessJustifications() (err error, ready bool) {
	if !m.IsJustificationsReady() {
		return nil, false
	}

	commits, err := m.GetCommits()
	if err != nil {
		return err, true
	}

	var (
		buf = bytes.NewBuffer(nil)
		enc = gob.NewEncoder(buf)
	)
	if err := enc.Encode(commits); err != nil {
		return fmt.Errorf("failed to encode response: %v", err), true
	}

	state := m.Dealer.GetState()

	message := &types.DKGData{
		Type:        types.DKGCommits,
		RoundID:     state.roundID,
		Addr:        state.addrBytes,
		Data:        buf.Bytes(),
		NumEntities: len(commits.Commitments),
	}

	m.SendMsgCb(message)

	return nil, true
}

func (m *DKGMockDontSendOneCommit) GetCommits() (*dkg.SecretCommits, error) {
	commits, err := m.Dealer.GetCommits()

	// remove one response message
	commits.Commitments = commits.Commitments[:len(commits.Commitments)-1]

	return commits, err
}

type DKGMockDontSendAnyCommits struct {
	Dealer
}

func NewDKGMockDealerAnyCommits(validators *types.ValidatorSet, pubKey crypto.PubKey, sendMsgCb func(*types.DKGData), logger log.Logger) Dealer {
	return &DKGMockDontSendAnyCommits{NewDKGDealer(validators, pubKey, sendMsgCb, logger)}
}

func (m *DKGMockDontSendAnyCommits) Start() error {
	err := m.Dealer.Start()
	if err != nil {
		return err
	}
	m.GenerateTransitions()
	return nil
}

func (m *DKGMockDontSendAnyCommits) GenerateTransitions() {
	m.Dealer.SetTransitions([]transition{
		// Phase I
		m.Dealer.SendDeals,
		m.Dealer.ProcessDeals,
		m.Dealer.ProcessResponses,
		m.ProcessJustifications,
		// Phase II
		m.Dealer.ProcessCommits,
		m.Dealer.ProcessComplaints,
		m.Dealer.ProcessReconstructCommits,
	})
}

func (m *DKGMockDontSendAnyCommits) ProcessJustifications() (error, bool) {
	if !m.Dealer.IsJustificationsReady() {
		return nil, false
	}

	fmt.Println("dkgState: sending commits", "commits", 0)

	return nil, true
}
