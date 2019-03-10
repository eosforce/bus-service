package basic

import (
	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/pbs/relay"
	eos "github.com/eosforce/goeosforce"
)

type commitParam struct {
	Name     eos.Name        `json:"chain"`
	Transfer eos.AccountName `json:"transfer"`
	Block    block           `json:"block"`
	Actions  []action        `json:"actions"`
}

func newCommitAction(relayBlock *force_relay_commit.RelayBlock, actionsToCommit []*force_relay_commit.RelayAction) *eos.Action {
	b := block{
		Producer:         eos.AN(relayBlock.Producer),
		Num:              blockNum(relayBlock.Id),
		ID:               relayBlock.Id,
		Previous:         relayBlock.Previous,
		Confirmed:        uint16(relayBlock.Confirmed),
		TransactionMRoot: relayBlock.TransactionMroot,
		ActionMRoot:      relayBlock.ActionMroot,
		MRoot:            relayBlock.Mroot,
	}

	seelog.Infof("commit block %d %v %d", b.Num, b.ID, len(actionsToCommit))

	acts := make([]action, 0, len(actionsToCommit)+1)
	for _, act := range actionsToCommit {
		auth := make([]permissionLevel, 0, 8)
		for _, authToTrx := range act.Authorization {
			auth = append(auth, permissionLevel{
				Actor:      eos.AN(authToTrx.Actor),
				Permission: eos.PN(authToTrx.Permission),
			})
		}
		acts = append(acts, action{
			Account:       eos.AN(act.Account),
			Name:          eos.ActN(act.ActionName),
			Authorization: auth,
			Data:          act.Data,
		})
	}

	return &eos.Action{
		Account: eos.AN("force.relay"),
		Name:    eos.ActN("commit"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(cfg.Chain), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(commitParam{
			Name:     cfg.Chain,
			Transfer: cfg.TransferAccount,
			Block:    b,
			Actions:  acts,
		}),
	}
}
