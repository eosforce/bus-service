package main

import (
	"errors"

	"github.com/eosforce/bus-service/force-relay/logger"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-go/p2p"
	"github.com/fanyang1988/force-go/types"
)

func startRelayService() {
	// from relay to side, so create side client
	_, p2ps := cfg.GetChainCfg("relay")

	// for chain id
	info, err := side.Client().GetInfoData()
	if err != nil {
		panic(errors.New("get info err"))
	}

	p2pPeers := p2p.NewP2PClient(types.EOSForce, p2p.P2PInitParams{
		Name:       "testNode",
		ClientID:   info.ChainID.String(),
		StartBlock: nil,
		Peers:      p2ps,
		Logger:     logger.Logger(),
	})

	p2pPeers.RegHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, 0, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				relay.HandRelayBlock(block, actions)
			})),
	})
	p2pPeers.Start()
}
