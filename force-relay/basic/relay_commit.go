package basic

import (
	"encoding/binary"

	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/pbs/relay"
	eos "github.com/eosforce/goeosforce"
)

type block struct {
	Producer         eos.AccountName `json:"producer"`
	Num              uint32          `json:"num"`
	ID               eos.Checksum256 `json:"id"`
	Previous         eos.Checksum256 `json:"previous"`
	Confirmed        uint16          `json:"confirmed"`
	TransactionMRoot eos.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      eos.Checksum256 `json:"action_mroot"`
	MRoot            eos.Checksum256 `json:"mroot"`
}

type commitParam struct {
	Name     eos.Name        `json:"chain"`
	Transfer eos.AccountName `json:"transfer"`
	Block    block           `json:"block"`
	Actions  []action        `json:"actions"`
}

type action struct {
	Account       eos.AccountName   `json:"account"`
	Name          eos.ActionName    `json:"name"`
	Authorization []permissionLevel `json:"authorization"`
	Data          []byte            `json:"data"`
}

type permissionLevel struct {
	Actor      eos.AccountName    `json:"actor"`
	Permission eos.PermissionName `json:"permission"`
}

var chain eos.Name
var transfer eos.AccountName

func SetChain(chainName string) {
	chain = eos.Name(chainName)
}

func SetTransfer(transferName string) {
	transfer = eos.AccountName(transferName)
}

func blockNum(blockID []byte) uint32 {
	if len(blockID) < 32 {
		return 0
	}
	return binary.BigEndian.Uint32(blockID[:32])
}

func newCommitAction(relayblock *force_relay_commit.RelayBlock, actionsToCommit []*force_relay_commit.RelayAction) *eos.Action {
	b := block{
		Producer:         eos.AN(relayblock.Producer),
		Num:              blockNum(relayblock.Id),
		ID:               relayblock.Id,
		Previous:         relayblock.Previous,
		Confirmed:        uint16(relayblock.Confirmed),
		TransactionMRoot: relayblock.TransactionMroot,
		ActionMRoot:      relayblock.ActionMroot,
		MRoot:            relayblock.Mroot,
	}

	seelog.Infof("commit block %d %v %d", b.Num, b.ID, len(actionsToCommit))

	acts := make([]action, 0, len(actionsToCommit)+1)
	for _, act := range actionsToCommit {
		auth := make([]permissionLevel, 0, 8)
		for _, authori := range act.Authorization {
			auth = append(auth, permissionLevel{
				Actor:      eos.AN(authori.Actor),
				Permission: eos.PN(authori.Permission),
			})
		}
		acts = append(acts, action{
			Account:       eos.AN(act.Account),
			Name:          eos.ActN(act.ActionName),
			Authorization: auth,
			Data:          act.Data,
		})
		//seelog.Infof("action %v", act)
	}

	return &eos.Action{
		Account: eos.AN("force.relay"),
		Name:    eos.ActN("commit"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AccountName(chain), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(commitParam{
			Name:     chain,
			Transfer: transfer,
			Block:    b,
			Actions:  acts,
		}),
	}
}
