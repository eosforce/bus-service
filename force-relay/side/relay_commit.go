package side

import (
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	forceio "github.com/eosforce/goforceio"
	"github.com/fanyang1988/force-go/types"
)

type BlockToForceio struct {
	Producer         forceio.AccountName `json:"producer"`
	Num              uint32              `json:"num"`
	ID               forceio.Checksum256 `json:"id"`
	Previous         forceio.Checksum256 `json:"previous"`
	Confirmed        uint16              `json:"confirmed"`
	TransactionMRoot forceio.Checksum256 `json:"transaction_mroot"`
	ActionMRoot      forceio.Checksum256 `json:"action_mroot"`
	MRoot            forceio.Checksum256 `json:"mroot"`
}

func (b *BlockToForceio) FromGeneral(sw types.SwitcherInterface, bk *chainhandler.Block) {
	b.Producer = forceio.AN(bk.Producer)
	b.Num = bk.Num
	b.Confirmed = bk.Confirmed
	b.ID = forceio.Checksum256(bk.ID)
	b.Previous = forceio.Checksum256(bk.Previous)
	b.TransactionMRoot = forceio.Checksum256(bk.TransactionMRoot)
	b.ActionMRoot = forceio.Checksum256(bk.ActionMRoot)
	b.MRoot = forceio.Checksum256(bk.MRoot)
}

type commitParam struct {
	Name     interface{}           `json:"chain"`
	Transfer interface{}           `json:"transfer"`
	Block    BlockToForceio        `json:"block"`
	Actions  []chainhandler.Action `json:"actions"`
}
