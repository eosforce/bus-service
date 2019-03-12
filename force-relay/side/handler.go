package side

import (
	"fmt"

	"github.com/eosforce/bus-service/force-relay/chainhandler"

	"github.com/cihub/seelog"
)

// HandSideBlock handle block from side chain
func HandSideBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	blockCommitLast, err := GetLastCommittedBlock()
	if err != nil {
		seelog.Errorf("get last commit block err by %v", err.Error())
	} else {
		lastNum := blockCommitLast.GetNum()
		num := block.GetNum()
		if blockCommitLast != nil && num != 0 && lastNum >= num {
			seelog.Debugf("no need commit %v to %v", lastNum, num)
			return
		}
	}

	commitAct := newCommitAction(block, actions)
	_, err = client.PushActions(commitAct)
	if err != nil {
		fmt.Println("push action error  ", err.Error())
		return
	}
}
