package side

import (
	"strings"
	"time"

	"github.com/eosforce/bus-service/force-relay/cfg"
	"github.com/eosforce/bus-service/force-relay/chainhandler"
	"github.com/eosforce/bus-service/force-relay/logger"
	"github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
	"github.com/fanyang1988/force-go/types"
)

const (
	sideBlockNumPerSec  = 2
	relayBlockNumPerSec = 2
	maxBlockPerTrx      = 1
	timeBetweenTrx      = relayBlockNumPerSec*3 + 1
	retryTimes          = 100
)

type commitWorker struct {
	committer types.PermissionLevel
	works     chan commitParam
	client    types.ClientInterface
}

type commitWorkers struct {
	cws []*commitWorker
}

func newCommitWorkers(clientCfg *config.ConfigData, committers []cfg.Relayer) *commitWorkers {
	res := &commitWorkers{
		cws: make([]*commitWorker, 0, len(committers)),
	}
	for _, c := range committers {
		cw := &commitWorker{
			committer: c.RelayAccount,
		}
		cw.Start(clientCfg)
		res.cws = append(res.cws, cw)
	}

	return res
}

var commitWorkerMng *commitWorkers

func InitCommitWorker(clientCfg *config.ConfigData, committers []cfg.Relayer) {
	commitWorkerMng = newCommitWorkers(clientCfg, committers)
}

func (c *commitWorkers) OnBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	for _, cw := range c.cws {
		cw.OnBlock(block, actions)
	}
}

func (c *commitWorker) Start(cfg *config.ConfigData) {
	c.works = make(chan commitParam, 4096)
	for {
		client, err := force.NewClient(types.FORCEIO, cfg)
		if err != nil {
			logger.LogError("create client error, need retry", err)
			time.Sleep(1 * time.Second)
		} else {
			c.client = client
			break
		}
	}

	go func(cc *commitWorker) {
		cc.Loop()
	}(c)
}

func (c *commitWorker) OnBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	c.works <- commitParam{
		Name:     c.client.Name(cfg.GetRelayCfg().Chain),
		Transfer: c.client.Name(c.committer.Actor),
		Block:    *block,
		Actions:  actions,
	}
}

func (c *commitWorker) Loop() {
	ticker := time.NewTicker(timeBetweenTrx * time.Second)
	defer ticker.Stop()
	works2do := make([]commitParam, 0, 4096)[:]
	for {
		select {
		case w := <-c.works:
			works2do = append(works2do, w)
			if len(works2do) >= maxBlockPerTrx {
				c.CommitTrx(works2do)
				works2do = works2do[:0]
				time.Sleep(50 * time.Millisecond)
			}
		default:
		}

		select {
		case <-ticker.C:
			if len(works2do) > 0 {
				c.CommitTrx(works2do)
				works2do = works2do[:0]
			}
		default:
		}
	}
}

func (c *commitWorker) CommitTrx(cps []commitParam) {
	actions := make([]*types.Action, 0, len(cps))

	for _, cp := range cps {
		actions = append(actions, &types.Action{
			Account: "force.relay",
			Name:    "commit",
			Authorization: []types.PermissionLevel{
				c.committer,
			},
			Data: cp,
		})
	}

	logger.Debugf("commit %s blocks num : %d -> %d",
		string(c.committer.Actor), cps[0].Block.Num, cps[len(cps)-1].Block.Num)

	for i := 0; i < retryTimes; i++ {
		if i > 1 {
			time.Sleep(10 * time.Millisecond)
		}
		_, err := c.client.PushActions(actions...)

		if err != nil {
			logger.LogError("commit action err", err)
			if strings.Contains(err.Error(), "Transaction took too long") {
				logger.Warnf("need wait chain err by took too long")
				time.Sleep(8 * time.Second)
			}

			if strings.Contains(err.Error(), "RAM") {
				logger.Warnf("need wait other chain err by RAM")
				time.Sleep(8 * time.Second)
			}
		} else {
			break
		}
	}

}
