package chainhandler

import (
	"context"
	"sync"

	"github.com/cihub/seelog"

	eos "github.com/eosforce/goforceio"

	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
)

type HandlerFunc func(block *Block, actions []Action)

type ChainHandler struct {
	handler    HandlerFunc
	blockQueue chan blockQueueItem
	wg         sync.WaitGroup
}

func NewChainHandler(h HandlerFunc) *ChainHandler {
	res := &ChainHandler{
		handler:    h,
		blockQueue: make(chan blockQueueItem, 4096),
	}

	res.wg.Add(1)
	go func(ch *ChainHandler) {
		defer ch.wg.Done()
		seelog.Infof("start chain handler")
		for {
			bi, ok := <-ch.blockQueue
			if !ok {
				seelog.Warnf("handler chan close")
				return
			}
			seelog.Tracef("process block %d", bi.block.Num)
			ch.handler(&bi.block, bi.actions)
		}
	}(res)

	return res
}

func (s *ChainHandler) RpcSendaction(ctx context.Context, in *commit.RelayCommitRequest) (*commit.RelayCommitReply, error) {
	var bqi blockQueueItem
	bqi.block = Block{
		Producer:         eos.AN(in.Block.Producer),
		Num:              BlockID2Num(in.Block.Id),
		ID:               in.Block.Id,
		Previous:         in.Block.Previous,
		Confirmed:        uint16(in.Block.Confirmed),
		TransactionMRoot: in.Block.TransactionMroot,
		ActionMRoot:      in.Block.ActionMroot,
		MRoot:            in.Block.Mroot,
	}

	actions := make([]Action, 0, len(in.Action))
	for _, act := range in.Action {
		auth := make([]PermissionLevel, 0, 8)
		for _, authToTrx := range act.Authorization {
			auth = append(auth, PermissionLevel{
				Actor:      eos.AN(authToTrx.Actor),
				Permission: eos.PN(authToTrx.Permission),
			})
		}
		actions = append(actions, Action{
			Account:       eos.AN(act.Account),
			Name:          eos.ActN(act.ActionName),
			Authorization: auth,
			Data:          act.Data,
		})
	}
	bqi.actions = actions[:]
	s.blockQueue <- bqi

	return &commit.RelayCommitReply{Reply: "get Block"}, nil
}

func (s *ChainHandler) Close() {
	close(s.blockQueue)
	s.wg.Wait()
}
