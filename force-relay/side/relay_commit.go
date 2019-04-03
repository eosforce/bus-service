package side

import (
	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	eos "github.com/eosforce/goforceio"
)

type commitParam struct {
	Name     eos.Name              `json:"chain"`
	Transfer eos.AccountName       `json:"transfer"`
	Block    chainhandler.Block    `json:"block"`
	Actions  []chainhandler.Action `json:"actions"`
}

func newCommitAction(b *chainhandler.Block, actionsToCommit []chainhandler.Action) *eos.Action {
	seelog.Infof("commit block %d %v %d", b.GetNum(), b.ID, len(actionsToCommit))
	return &eos.Action{
		Account: eos.AN("force.relay"),
		Name:    eos.ActN("commit"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(cfg.TransferAccount), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(commitParam{
			Name:     cfg.Chain,
			Transfer: cfg.TransferAccount,
			Block:    *b,
			Actions:  actionsToCommit,
		}),
	}
}
