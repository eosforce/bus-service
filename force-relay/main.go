package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/eosforce/goforceio/p2p"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/logger"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/eosforce/goforceio/ecc"
	blockevlog "github.com/fanyang1988/force-block-ev/log"
)

var configPath = flag.String("cfg", "./config.json", "config file path")
var isDebug = flag.Bool("d", false, "run in debug mode")

func init() {
	ecc.PublicKeyPrefixCompat = "FOSC"
}

func main() {
	flag.Parse()
	logger.EnableLogging(*isDebug)
	blockevlog.SetLogger(logger.Logger())
	if *isDebug {
		p2p.EnableP2PLogging()
	}

	defer func() {
		err := logger.Logger().Sync()
		if err != nil {
			fmt.Printf("logger sync err by %s", err.Error())
		}
	}()

	runtime.GOMAXPROCS(2)

	err := cfg.LoadCfgs(*configPath)
	if err != nil {
		logger.Sugar().Errorf("load cfg err by %s", err.Error())
		return
	}

	sideChainCfg, _ := cfg.GetChainCfg("side")
	sideChainCfg.IsDebug = *isDebug
	relay.CreateSideClient(sideChainCfg)

	relayChainCfg, _ := cfg.GetChainCfg("relay")
	relayChainCfg.IsDebug = *isDebug
	side.CreateClient(relayChainCfg)

	go func() {
		if len(cfg.GetWatchers()) == 0 {
			logger.Sugar().Infof("no need start relay")
			return
		}
		logger.Sugar().Infof("start relay service")
		startRelayService()
	}()

	go func() {
		if len(cfg.GetTransfers()) == 0 {
			logger.Sugar().Infof("no need start side")
			return
		}
		logger.Sugar().Infof("start side service")
		startSideService()
	}()

	for {
		time.Sleep(1 * time.Second)
		// TODO check status
	}
}
