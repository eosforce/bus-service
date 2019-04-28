package side

import (
	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
	"go.uber.org/zap"
)

type commitParam struct {
	Name     eos.Name              `json:"chain"`
	Transfer eos.AccountName       `json:"transfer"`
	Block    chainhandler.Block    `json:"block"`
	Actions  []chainhandler.Action `json:"actions"`
}

func newCommitAction(b *chainhandler.Block, transfer eos.PermissionLevel, actionsToCommit []chainhandler.Action) *eos.Action {
	logger.Logger().Info("commit block",
		zap.Uint32("num", b.GetNum()),
		zap.String("id", b.ID.String()),
		zap.Int("action", len(actionsToCommit)),
		zap.String("previous", b.Previous.String()))
	return &eos.Action{
		Account: eos.AN("force.relay"),
		Name:    eos.ActN("commit"),
		Authorization: []eos.PermissionLevel{
			transfer,
		},
		ActionData: eos.NewActionData(commitParam{
			Name:     eos.Name(cfg.GetRelayCfg().Chain),
			Transfer: transfer.Actor,
			Block:    *b,
			Actions:  actionsToCommit,
		}),
	}
}
