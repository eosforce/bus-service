package main

import (
	"errors"
	"fmt"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-block-ev/blockev"
)

func startSideService() {
	// frome side need to commit block to relay
	chainCfgs, _ := cfg.GetChainCfg("relay")
	_, p2ps := cfg.GetChainCfg("side")
	side.InitCommitWorker(chainCfgs, cfg.GetTransfers())

	// for p2p chain id
	info, err := relay.Client().GetInfo()
	if err != nil {
		panic(errors.New("get info err"))
	}

	lastCommitted, err := side.GetLastCommittedBlock()
	if err != nil {
		panic(errors.New("get info err"))
	}

	lastNum := lastCommitted.GetNum()
	if lastNum > 3 {
		lastNum -= 2
	}

	lastBlock, err := relay.Client().GetBlockByNum(lastNum)
	if err != nil {
		panic(fmt.Errorf("err by %s", err.Error()))
	}

	p2pPeers := blockev.NewP2PPeers("relay", info.ChainID.String(), &lastBlock.BlockHeader, p2ps)
	p2pPeers.RegisterHandler(blockev.NewP2PMsgHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				side.HandSideBlock(block, actions)
			})),
	}))
	p2pPeers.Start()

}
