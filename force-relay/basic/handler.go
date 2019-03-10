package basic

import (
	"fmt"

	"github.com/eosforce/bus-service/force-relay/pbs/relay"
)

// HandRelayBlock handle block from side chain
func HandRelayBlock(block *force_relay_commit.RelayBlock, Action []*force_relay_commit.RelayAction) {
	blockCommitLast, err := GetLastCommittedBlock()
	num := blockNum(block.Id)
	if err == nil && blockCommitLast != nil && num != 0 && blockCommitLast.Num >= num {
		return
	}

	commitAct := newCommitAction(block, Action)
	_, err = client.PushActions(commitAct)
	if err != nil {
		fmt.Println("push action error  ", err.Error())
		return
	}
}
