package side

import (
	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
	"github.com/fanyang1988/force-go/types"
	"go.uber.org/zap"
)

type commitParam struct {
	Name     interface{}           `json:"chain"`
	Transfer interface{}           `json:"transfer"`
	Block    chainhandler.Block    `json:"block"`
	Actions  []chainhandler.Action `json:"actions"`
}

func newCommitAction(b *chainhandler.Block, transfer types.PermissionLevel, actionsToCommit []chainhandler.Action) *types.Action {
	logger.Logger().Info("commit block",
		zap.Uint32("num", b.GetNum()),
		zap.String("id", b.ID.String()),
		zap.Int("action", len(actionsToCommit)),
		zap.String("previous", b.Previous.String()))
	return &types.Action{
		Account: "force.relay",
		Name:    "commit",
		Authorization: []types.PermissionLevel{
			transfer,
		},
		Data: commitParam{
			Name:     cfg.GetRelayCfg().Chain,
			Transfer: transfer.Actor,
			Block:    *b,
			Actions:  actionsToCommit,
		},
	}
}
