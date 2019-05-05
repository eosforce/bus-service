package main

import (
	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/fanyang1988/force-block-ev/blockdb"
	"github.com/fanyang1988/force-go/p2p"
	"github.com/fanyang1988/force-go/types"
	"github.com/pkg/errors"
)

func startSideService() {
	// frome side need to commit block to relay
	chainCfgs, _ := cfg.GetChainCfg("relay")
	_, p2ps := cfg.GetChainCfg("side")
	side.InitCommitWorker(chainCfgs, cfg.GetTransfers())

	// for p2p chain id
	info, err := relay.Client().GetInfoData()
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

	lastBlockData, err := relay.Client().GetBlockDataByNum(lastNum)
	if err != nil {
		panic(errors.Errorf("get block num %d err by %s", lastNum, err.Error()))
	}

	p2pPeers := p2p.NewP2PClient(types.FORCEIO, p2p.P2PInitParams{
		Name:     "relay",
		ClientID: info.ChainID.String(),
		StartBlock: &p2p.P2PSyncData{
			HeadBlockNum:             lastBlockData.BlockNum,
			HeadBlockID:              lastBlockData.ID,
			HeadBlockTime:            lastBlockData.Timestamp,
			LastIrreversibleBlockNum: lastBlockData.BlockNum,
			LastIrreversibleBlockID:  lastBlockData.ID,
		},
		Peers:  p2ps,
		Logger: logger.Logger(),
	})

	p2pPeers.RegHandler(&handlerImp{
		verifier: blockdb.NewFastBlockVerifier(p2ps, 0, chainhandler.NewChainHandler(
			func(block *chainhandler.Block, actions []chainhandler.Action) {
				relay.HandRelayBlock(block, actions)
			})),
	})
	p2pPeers.Start()
}
