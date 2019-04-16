package side

import (
	"time"

	"github.com/eosforce/bus-service/force-relay/cfg"

	"github.com/eosforce/bus-service/force-relay/chainhandler"

	"github.com/cihub/seelog"
)

var lastCommittedBlockNum uint32

// HandSideBlock handle block from side chain
func HandSideBlock(block *chainhandler.Block, actions []chainhandler.Action) {
	if handSideBlockImp(block, actions) {
		time.Sleep(5 * time.Millisecond)
	}
}

// HandSideBlock handle block from side chain
func handSideBlockImp(block *chainhandler.Block, actions []chainhandler.Action) bool {
	const retryTimes int = 32

	num := block.GetNum()

	var blockCommitLast *chainhandler.Block
	var err error

	if lastCommittedBlockNum > 0 && num != 0 && lastCommittedBlockNum >= num {
		seelog.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return false
	}

	if lastCommittedBlockNum == 0 {
		blockCommitLast, err = GetLastCommittedBlock()
		if err != nil {
			seelog.Errorf("get last commit block err by %v", err.Error())
		}
		lastCommittedBlockNum = blockCommitLast.GetNum()
	}

	if blockCommitLast != nil && num != 0 && lastCommittedBlockNum >= num {
		seelog.Debugf("no need commit %v to %v", lastCommittedBlockNum, num)
		return false
	}

	transfers := cfg.GetTransfers()

	for _, t := range transfers {
		commitAct := newCommitAction(block, t.RelayAccount, actions)
		for i := 0; i < retryTimes; i++ {
			if i > 1 {
				time.Sleep(3500 * time.Millisecond)
			}
			_, err = client.PushActions(commitAct)
			if err != nil {
				seelog.Errorf("push action error %s", err.Error())

			} else {
				break
			}
		}
	}

	return true
}
