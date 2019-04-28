package main

import (
	eos "github.com/eosforce/goforceio"
	"github.com/fanyang1988/force-block-ev/blockdb"
)

type handlerImp struct {
	verifier *blockdb.FastBlockVerifier
}

func (h *handlerImp) OnBlock(peer string, msg *eos.SignedBlock) error {
	return h.verifier.OnBlock(peer, msg)
}
func (h *handlerImp) OnGoAway(peer string, msg *eos.GoAwayMessage) error {
	return nil
}
func (h *handlerImp) OnHandshake(peer string, msg *eos.HandshakeMessage) error {
	return nil
}
func (h *handlerImp) OnTimeMsg(peer string, msg *eos.TimeMessage) error {
	return nil
}
