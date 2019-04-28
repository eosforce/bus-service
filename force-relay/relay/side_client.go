package relay

import (
	"time"

	"github.com/eosforce/bus-service/force-relay/logger"
	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
	"go.uber.org/zap"
)

// client client to force relay chain
var client *force.Client

// CreateSideClient create client to force side chain
func CreateSideClient(cfg *config.Config) {
	for {
		var err error
		logger.Logger().Info("create client cfg",
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
