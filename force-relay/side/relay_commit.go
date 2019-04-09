package side

import (
	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	eos "github.com/eosforce/goforceio"
)

type commitParam struct {
	Name     eos.Name              `json:"chain"`
	Transfer eos.AccountName       `json:"transfer"`
	Block    chainhandler.Block    `json:"block"`
	Actions  []chainhandler.Action `json:"actions"`
}

func newCommitAction(b *chainhandler.Block, transfer eos.PermissionLevel, actionsToCommit []chainhandler.Action) *eos.Action {
	seelog.Infof("commit block %d %v %d %v", b.GetNum(), b.ID, len(actionsToCommit), b.Previous)
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
