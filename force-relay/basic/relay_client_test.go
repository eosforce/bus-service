package basic

import (
	"testing"

	"github.com/cihub/seelog"
)

func initCfgForTest() {
	CreateClient("../config.json")
	SetCfg(NewCfg("eosforce", "eosforce"))
}

func TestGetLastCommittedBlock(t *testing.T) {
	defer seelog.Flush()
	initCfgForTest()

	rsp, err := GetLastCommittedBlock()
	if err != nil {
		t.Errorf("err by %v", err.Error())
		t.FailNow()
	}

	t.Logf("rsp info %v", rsp)
	t.Logf("rsp info %v", rsp.ID)
	t.Logf("rsp info %v", rsp.Num)
	t.Logf("rsp info %v", rsp.Producer)
}
