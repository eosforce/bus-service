package side

import (
	"time"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// client client to force relay chain
var client *force.Client

// CreateClient create client to force relay chain
func CreateClient(cfg *config.Config) {
	for {
		var err error
		logger.Logger().Debug("create client cfg",
			zap.String("url", cfg.URL),
			zap.String("chainID", cfg.ChainID.String()),
			zap.Bool("isDebug", cfg.IsDebug))
		client, err = force.NewClient(cfg)
		if err != nil {
			logger.LogError("create client error, need retry", err)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}

func Client() *force.Client {
	return client
}

type lastCommitBlockInfo struct {
	Chain eos.Name           `json:"chain"`
	Last  chainhandler.Block `json:"last"`
}

// GetLastCommittedBlock get last committed block to relay chain
func GetLastCommittedBlock() (*chainhandler.Block, error) {
	req := eos.GetTableRowsRequest{
		Code:  "force.relay",
		Scope: cfg.GetRelayCfg().Chain,
		Table: "relaystat",
	}

	res, err := client.GetTableRows(req)
	if err != nil {
		return nil, err
	}

	rspBlock := make([]lastCommitBlockInfo, 0, 32)
	err = res.BinaryToStructs(&rspBlock)
	if err != nil {
		return nil, err
	}

	if len(rspBlock) == 0 {
		return nil, errors.New("rsp block info no find")
	}

	return &rspBlock[0].Last, nil
}
