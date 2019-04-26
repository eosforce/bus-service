package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/cihub/seelog"
	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/eosforce/goforceio/ecc"
	"github.com/fanyang1988/force-block-ev/log"
)

var configPath = flag.String("cfg", "./config.json", "confg file path")
var chain = flag.String("chain", "eosforce", "the name of chain")
var transfer = flag.String("transfer", "eosforce", "the name of transfer")

func init() {
	ecc.PublicKeyPrefixCompat = "FOSC"
}

func main() {
	flag.Parse()
	defer seelog.Flush()

	runtime.GOMAXPROCS(8)

	log.EnableLogging(false)

	err := cfg.LoadCfgs(*configPath)
	if err != nil {
		seelog.Errorf("load cfg err by %s", err.Error())
		return
	}

	seelog.Infof("dd %s", ecc.PublicKeyPrefixCompat)

	sideChainCfgs, _ := cfg.GetChainCfg("side")
	relay.CreateSideClient(sideChainCfgs)
	relayChainCfgs, _ := cfg.GetChainCfg("relay")
	side.CreateClient(relayChainCfgs)

	go func() {
		if len(cfg.GetWatchers()) == 0 {
			seelog.Infof("no need start relay")
			return
		}
		seelog.Infof("start relay service")
		startRelayService()
	}()

	go func() {
		if len(cfg.GetTransfers()) == 0 {
			seelog.Infof("no need start side")
			return
		}
		seelog.Infof("start side service")
		startSideService()
	}()

	for {
		time.Sleep(1 * time.Second)
		// TODO check status
	}
}
