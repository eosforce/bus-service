package relay

import (
	"fmt"

	force "github.com/fanyang1988/force-go"
	"github.com/fanyang1988/force-go/config"
)

// client client to force relay chain
var client *force.Client

// CreateSideClient create client to force side chain
func CreateSideClient(cfg *config.Config) {
	var err error
	client, err = force.NewClient(cfg)
	if err != nil {
		fmt.Println("create client error  ", err.Error())
		panic(err)
		return
	}
}
