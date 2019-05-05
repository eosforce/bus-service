package relay

import (
	"fmt"

	"github.com/fanyang1988/force-go/types"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
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
	logger.Debugf("on block from relay %d", block.GetNum())
	for _, act := range actions {
		if act.Account != "relay.token" || act.Name != "destroy" {
			continue
		}

		var actData Destroy
		err := eos.UnmarshalBinary(act.Data, &actData)
		if err != nil {
			logger.LogError("UnmarshalBinary act err", err)
			continue
		}

		onTokenReturnSideChain(block, &actData)
	}
}

func onTokenReturnSideChain(block *chainhandler.Block, act *Destroy) {
	logger.Debugf("on return in block %d : %s %v by %v in %s",
		block.GetNum(), act.Chain, act.From, act.Quantity, act.Memo)

	num = num + 1
	for _, w := range cfg.GetWatchers() {
		commitOutAction(w, act)
	}
}

// OutAction  capi_name committer, uint64_t num, capi_name to, name chain, name contract, const asset& quantity, const std::string& memo
type OutAction struct {
	Committer eos.Name  `json:"committer"`
	Num       uint64    `json:"num"`
	To        eos.Name  `json:"to"`
	Chain     eos.Name  `json:"chain"`
	Contract  eos.Name  `json:"contract"`
	Quantity  eos.Asset `json:"quantity"`
	Memo      string    `json:"memo"`
}

// TODO
var num uint64

func commitOutAction(committer cfg.Relayer, act *Destroy) error {
	actToCommit := &types.Action{
		Account: cfg.GetRelayCfg().RelayContract,
		Name:    "out",
		Authorization: []types.PermissionLevel{
			committer.SideAccount,
		},
		Data: OutAction{
			Committer: eos.Name(committer.RelayAccount.Actor),
			Num:       num,
			To:        eos.Name(act.From),
			Chain:     act.Chain,
			Contract:  eos.Name("force.token"),
			Quantity:  act.Quantity,
			Memo:      act.Memo,
		},
	}

	_, err := client.PushActions(actToCommit)
	if err != nil {
		fmt.Println("push action error  ", err.Error())
	}

	return err
}
