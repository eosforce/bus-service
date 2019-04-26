package main

import (
	"errors"

	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/side"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-block-ev/blockev"

	"github.com/eosforce/bus-service/force-relay/relay"
)

func startRelayService() {
	// from relay to side, so create side client
	_, p2ps := cfg.GetChainCfg("relay")

	// for chain id
	info, err := side.Client().GetInfo()
	if err != nil {
		panic(errors.New("get info err"))
	}

	p2pPeers := blockev.NewP2PPeers("side", info.ChainID.String(), nil, p2ps)
	p2pPeers.RegisterHandler(blockev.NewP2PMsgHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				relay.HandRelayBlock(block, actions)
			})),
	}))
	p2pPeers.Start()

}
