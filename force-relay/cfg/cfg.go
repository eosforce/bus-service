package cfg

import (
	"errors"

	"github.com/eosforce/bus-service/force-relay/logger"
	eos "github.com/eosforce/goforceio"
	"github.com/fanyang1988/force-go/config"
)

// Relayer transfer, watcher and checker
type Relayer struct {
	SideAccount  eos.PermissionLevel
	RelayAccount eos.PermissionLevel
}

// RelayCfg cfg for relay
type RelayCfg struct {
	Chain         string `json:"chain"`
	RelayContract string `json:"relaycontract"`
}

var relayCfg RelayCfg

// ChainCfgs cfg for each chain
var chainCfgs map[string]*config.Config
var chainP2PCfgs map[string][]string

var transfers []Relayer
var watchers []Relayer

// GetChainCfg get chain cfg
func GetChainCfg(name string) (*config.Config, []string) {
	c, ok := chainCfgs[name]
	if !ok || c == nil {
		panic(errors.New("no find chain cfg "))
	}
	p, ok := chainP2PCfgs[name]
	if !ok {
		panic(errors.New("no find chain p2p cfg "))
	}
	return c, p
}

// GetTransfers get transfers
func GetTransfers() []Relayer {
	return transfers
}

// GetRelayCfg get cfg for relay
func GetRelayCfg() RelayCfg {
	return relayCfg
}

// GetWatchers get watchers
func GetWatchers() []Relayer {
	return watchers
}

// LoadCfgs load cfg for force-relay
func LoadCfgs(path string) error {
	cfgInFile := struct {
		Chains []struct {
			Name string            `json:"name"`
			P2P  []string          `json:"p2p"`
			Cfg  config.ConfigData `json:"cfg"`
		} `json:"chains"`
		Transfer []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"transfer"`
		Watcher []struct {
			SideAcc  string `json:"sideacc"`
			RelayAcc string `json:"relayacc"`
		} `json:"watcher"`
		Relay RelayCfg `json:"relay"`
	}{}

	err := config.LoadJSONFile(path, &cfgInFile)
	if err != nil {
		return err
	}

	chainCfgs = make(map[string]*config.Config)
	for _, c := range cfgInFile.Chains {
		cc := config.Config{}
		err := cc.Parse(&c.Cfg)
		logger.Debugf("load cfg %v", cc)
		if err != nil {
			return err
		}
		chainCfgs[c.Name] = &cc
	}

	chainP2PCfgs = make(map[string][]string)
	for _, c := range cfgInFile.Chains {
		logger.Debugf("load p2p cfg %s -> %v", c.Name, c.P2P)
		chainP2PCfgs[c.Name] = c.P2P
	}

	for _, t := range cfgInFile.Transfer {
		transfers = append(transfers, Relayer{
			SideAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.SideAcc),
				Permission: eos.PN("active"),
			},
			RelayAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.RelayAcc),
				Permission: eos.PN("active"),
			},
		})
	}

	for _, t := range cfgInFile.Watcher {
		watchers = append(watchers, Relayer{
			SideAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.SideAcc),
				Permission: eos.PN("active"),
			},
			RelayAccount: eos.PermissionLevel{
				Actor:      eos.AN(t.RelayAcc),
				Permission: eos.PN("active"),
			},
		})
	}

	relayCfg = cfgInFile.Relay

	return nil
}
