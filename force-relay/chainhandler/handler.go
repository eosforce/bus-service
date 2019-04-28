package chainhandler

import (
	"sync"

	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
	"go.uber.org/zap"
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
		logger.Logger().Info("start chain handler")
		for {
			bi, ok := <-ch.blockQueue
			if !ok {
				logger.Logger().Error("handler chan close")
				return
			}
			ch.handler(&bi.block, bi.actions)
		}
	}(res)

	return res
}

func (c *ChainHandler) OnBlock(blockNum uint32, blockID eos.Checksum256, block *eos.SignedBlock) error {
	logger.Logger().Debug("on block",
		zap.Uint32("num", blockNum),
		zap.String("id", blockID.String()))
	var bqi blockQueueItem
	bqi.block = Block{
		Producer:         block.Producer,
		Num:              blockNum,
		ID:               blockID,
		Previous:         block.Previous,
		Confirmed:        uint16(block.Confirmed),
		TransactionMRoot: block.TransactionMRoot,
		ActionMRoot:      block.ActionMRoot,
	}

	actions := make([]Action, 0, 1024)
	for _, trx := range block.Transactions {
		if trx.Status != eos.TransactionStatusExecuted {
			continue
		}

		st, err := trx.Transaction.Packed.Unpack()
		if err != nil {
			continue
		}

		for _, act := range st.Actions {
			auth := make([]PermissionLevel, 0, 8)
			for _, authToTrx := range act.Authorization {
				auth = append(auth, PermissionLevel{
					Actor:      authToTrx.Actor,
					Permission: authToTrx.Permission,
				})
			}
			actions = append(actions, Action{
				Account:       act.Account,
				Name:          act.Name,
				Authorization: auth,
				Data:          act.HexData,
			})
		}
	}
	bqi.actions = actions[:]
	c.blockQueue <- bqi
	return nil
}

func (c *ChainHandler) Close() {
	close(c.blockQueue)
	c.wg.Wait()
}
