package side

import (
	"time"

	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
	forceio "github.com/fanyang1988/force-go/forceio"
	"github.com/fanyang1988/force-go/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
)

// client client to force relay chain
var client types.ClientInterface

// CreateClient create client to force relay chain
func CreateClient(typ types.ClientType, cfg *config.ConfigData) {
	for {
		var err error
		logger.Logger().Info("create client cfg",
			zap.String("url", cfg.URL),
			zap.String("chainID", cfg.ChainID))
		client, err = force.NewClient(typ, cfg)
		if err != nil {
			logger.LogError("create client error, need retry", err)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}

func Client() types.ClientInterface {
	return client
}

type lastCommitBlockInfo struct {
	Chain eos.Name       `json:"chain"`
	Last  BlockToForceio `json:"last"`
}

// GetLastCommittedBlock get last committed block to relay chain
func GetLastCommittedBlock() (*BlockToForceio, error) {
	req := eos.GetTableRowsRequest{
		Code:  "force.relay",
		Scope: cfg.GetRelayCfg().Chain,
		Table: "relaystat",
	}

	forceioClient, ok := client.(*forceio.API)
	if !ok {
		return nil, types.ErrNoSupportChain
	}

	res, err := forceioClient.GetTableRows(req)
	if err != nil {
		return nil, errors.Wrapf(err, "get table")
	}

	rspBlock := make([]lastCommitBlockInfo, 0, 32)
	err = res.BinaryToStructs(&rspBlock)
	if err != nil {
		return nil, errors.Wrapf(err, "to struct")
	}

	if len(rspBlock) == 0 {
		return nil, errors.New("rsp block info no find")
	}

	logger.Debugf("get last cm block %s from %d", cfg.GetRelayCfg().Chain, rspBlock[0].Last.Num)

	return &rspBlock[0].Last, nil
}

// ActionsToRelay actions need to relay
type ActionsToRelay struct {
	names     map[string]int
	Contracts []string
	Actions   []string
	keys      map[string]bool
}

func NewActionsToRelay() *ActionsToRelay {
	return &ActionsToRelay{
		names:     make(map[string]int),
		Contracts: make([]string, 0, 8),
		Actions:   make([]string, 0, 8),
		keys:      make(map[string]bool),
	}
}

func (a *ActionsToRelay) Append(name, actContract, actName string) {
	_, ok := a.names[name]
	if ok {
		return
	}
	a.names[name] = len(a.Contracts)
	a.Contracts = append(a.Contracts, actContract)
	a.Actions = append(a.Actions, actName)

	a.keys[actContract+"::"+actName] = true
}

func (a *ActionsToRelay) IsNeedCommit(contract, name string) bool {
	_, ok := a.keys[contract+"::"+name]
	return ok
}

/*
  "rows": [{
      "chain": "eosforce",
      "name": "eosforce.ct",
      "actaccount": "eosio",
      "actname": "transfer",
      "relayacc": "rs",
      "account": "relay.token",
      "data": ""
    },{
      "chain": "eosforce",
      "name": "eosforce.t",
      "actaccount": "eosio.token",
      "actname": "transfer",
      "relayacc": "rs",
      "account": "relay.token",
      "data": ""
    }
  ],
*/

type handlersInfo struct {
	Chain             eos.Name `json:"chain"`
	Name              eos.Name `json:"name"`
	ActionContract    eos.Name `json:"actaccount"`
	ActionName        eos.Name `json:"actname"`
	SideRelayName     eos.Name `json:"relayacc"`
	RelayContractName eos.Name `json:"account"`
	Data              string   `json:"data"`
}

// GetRelayActions get actions need to relay
func GetRelayActions() (*ActionsToRelay, error) {
	req := eos.GetTableRowsRequest{
		Code:  "force.relay",
		Scope: cfg.GetRelayCfg().Chain,
		Table: "handlers",
	}

	forceioClient, ok := client.(*forceio.API)
	if !ok {
		return nil, types.ErrNoSupportChain
	}

	// TODO get all
	tableRes, err := forceioClient.GetTableRows(req)
	if err != nil {
		return nil, errors.Wrapf(err, "get table")
	}

	rspHandlers := make([]handlersInfo, 0, 32)
	err = tableRes.BinaryToStructs(&rspHandlers)
	if err != nil {
		return nil, errors.Wrapf(err, "to struct err")
	}

	if len(rspHandlers) == 0 {
		return nil, errors.New("rsp block info no find")
	}

	logger.Debugf("get handlers %s from %v", cfg.GetRelayCfg().Chain, rspHandlers)

	res := NewActionsToRelay()
	for _, h := range rspHandlers {
		logger.Debugf("handler %s from %s:%s", h.Name, h.ActionContract, h.ActionName)
		res.Append(string(h.Name), string(h.ActionContract), string(h.ActionName))
	}

	return res, nil
}
