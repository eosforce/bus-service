package basic

import (
	"fmt"
	force_relay_commit "github.com/eosforce/bus-service/force_relay_commit"
	force "github.com/fanyang1988/force-go"
	//"github.com/eosforce/goeosforce/system"
)

//处理先关的块信息
var client *force.Client

//接下来构造Action并发送Action  在service上不做任何校验
func Createclient(configPath string) {
	var err error 
	client, err = force.NewClientFromFile(configPath)
	if err != nil {
		fmt.Println("create client error  ",err.Error())
		return
	}
}

func HandRelayBlock(block *force_relay_commit.RelayBlock,Action []*force_relay_commit.RelayAction) {
	commitAct := newCommitAction(block,Action)
	_, err := client.PushActions(commitAct)
	if err != nil {
		fmt.Println("push action error  ",err.Error())
		return
	}
}