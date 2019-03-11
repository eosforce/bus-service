package side

import (
	"fmt"

	"github.com/cihub/seelog"
	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
)

// HandRelayBlock handle block from side chain
func HandRelayBlock(block *commit.RelayBlock, Action []*commit.RelayAction) {
	blockCommitLast, err := GetLastCommittedBlock()
	if err != nil {
		seelog.Errorf("get last commit block err by %v", err.Error())
		return
	}
	lastNum := blockNum(blockCommitLast.ID)
	num := blockNum(block.Id)
	if err == nil && blockCommitLast != nil && num != 0 && lastNum >= num {
		seelog.Debugf("no need commit %v to %v", lastNum, num)
		return
	}

	commitAct := newCommitAction(block, Action)
	_, err = client.PushActions(commitAct)
	if err != nil {
		fmt.Println("push action error  ", err.Error())
		return
	}
}
