package relay

import (
	"time"

	"github.com/cihub/seelog"
	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
)

// client client to force relay chain
var client *force.Client

// CreateSideClient create client to force side chain
func CreateSideClient(cfg *config.Config) {
	for {
		var err error
		seelog.Tracef("cfg %v", *cfg)
		client, err = force.NewClient(cfg)
		if err != nil {
			seelog.Warnf("create client error by %s , need retry", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}

func Client() *force.Client {
	return client
}
