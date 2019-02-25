package basic

import (
	//"fmt"
	"encoding/json"

	eos "github.com/eosforce/goeosforce"
)

func HandTransaction(strjson, trxid string) {

	var trx eos.Transaction
	json.Unmarshal([]byte(strjson), &trx)
	for _, act := range trx.Actions {
		//先解析transfer
		if act.Name == toActionName("transfer", "action") {
			DecodeTransfer(act.ActionData.Data.(string), trxid)
		}
	}
}
