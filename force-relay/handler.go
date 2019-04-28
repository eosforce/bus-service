package main

import (
	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"go.uber.org/zap"
)

type handlerImp struct {
	verifier *blockdb.FastBlockVerifier
}

func (h *handlerImp) OnBlock(peer string, msg *eos.SignedBlock) error {
	logger.Logger().Info("onblock", zap.Uint32("num", msg.BlockNumber()))
	logger.Logger().Error("onblock", zap.Uint32("num", msg.BlockNumber()))
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
