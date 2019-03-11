package chainhandler

import (
	"context"

	eos "github.com/eosforce/goeosforce"

	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
)

type HandlerFunc func(block *Block, actions []Action)

type ChainHandler struct {
	handler HandlerFunc
}

func NewChainHandler(h HandlerFunc) *ChainHandler {
	return &ChainHandler{
		handler: h,
	}
}

func (s *ChainHandler) RpcSendaction(ctx context.Context, in *commit.RelayCommitRequest) (*commit.RelayCommitReply, error) {
	block := Block{
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

	s.handler(&block, actions)
	return &commit.RelayCommitReply{Reply: "get Block"}, nil
}
