package relay

import (
	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	eos "github.com/eosforce/goforceio"
)

// Destroy action data for relay.token::destroy action
type Destroy struct {
	Chain    eos.Name        `json:"chain"`
	From     eos.AccountName `json:"from"`
	Quantity eos.Asset       `json:"quantity"`
	Memo     string          `json:"memo"`
}

func HandRelayBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	for _, act := range actions {
		if act.Account != eos.AN("relay.token") || act.Name != eos.ActN("destroy") {
			continue
		}

		var actData Destroy
		err := eos.UnmarshalBinary(act.Data, &actData)
		if err != nil {
			seelog.Errorf("UnmarshalBinary act err by %s", err.Error())
			continue
		}

		onTokenReturnSideChain(block, &actData)
	}
}

func onTokenReturnSideChain(block *chainhandler.Block, act *Destroy) {
	seelog.Debugf("on return in block %d : %s %v by %v in %s",
		block.GetNum(), act.Chain, act.From, act.Quantity, act.Memo)
}
