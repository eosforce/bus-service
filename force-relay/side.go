package main

import (
	"github.com/eosforce/bus-service/force-relay/logger"
	"github.com/pkg/errors"

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
		panic(errors.New("GetLastCommittedBlock info err"))
	}

	logger.Debugf("get last committed block %v", lastCommitted)

	lastNum := lastCommitted.GetNum()
	if lastNum > 3 {
		lastNum -= 2
	}

	if lastNum == 0 {
		// no committed
		lastNum = 1
	}

	lastBlock, err := relay.Client().GetBlockByNum(lastNum)
	if err != nil {
		panic(errors.Errorf("get block num %d err by %s", lastNum, err.Error()))
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
