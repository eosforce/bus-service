package basic

import (
	"fmt"

	"github.com/eosforce/bus-service/force-relay/pbs/relay"
)

// HandRelayBlock handle block from side chain
func HandRelayBlock(block *force_relay_commit.RelayBlock, Action []*force_relay_commit.RelayAction) {
	commitAct := newCommitAction(block, Action)
	_, err := client.PushActions(commitAct)
	if err != nil {
		fmt.Println("push action error  ", err.Error())
		return
	}
}
