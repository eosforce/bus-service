package side

import (
	"time"

	"github.com/eosforce/bus-service/force-relay/chainhandler"

	"github.com/cihub/seelog"
)

var lastCommittedBlockNum uint32 = 0

// HandSideBlock handle block from side chain
func HandSideBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	const retryTimes int = 3

	num := block.GetNum()

	var blockCommitLast *chainhandler.Block
	var err error

	if lastCommittedBlockNum > 0 && num != 0 && lastCommittedBlockNum > num {
		seelog.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return
	}

	for i := 0; i < retryTimes; i++ {
		if i > 1 {
			time.Sleep(100 * time.Millisecond)
		}
		blockCommitLast, err = GetLastCommittedBlock()
		if err != nil {
			seelog.Errorf("get last commit block err by %v", err.Error())
		} else {
			break
		}
	}

	lastCommittedBlockNum = blockCommitLast.GetNum()
	if blockCommitLast != nil && num != 0 && lastCommittedBlockNum >= num {
		seelog.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return
	}

	commitAct := newCommitAction(block, actions)
	_, err = client.PushActions(commitAct)
	if err != nil {
		seelog.Errorf("push action error  ", err.Error())
		return
	}
}
