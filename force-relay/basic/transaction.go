package basic

import (
	"encoding/json"

	eos "github.com/eosforce/goeosforce"
)

func HandTransaction(strjson, trxid string) {
	var trx eos.Transaction
	json.Unmarshal([]byte(strjson), &trx)
	for _, act := range trx.Actions {
		//先解析transfer
		if act.Name == eos.ActN("transfer") {
			DecodeTransfer(act.ActionData.Data.(string), trxid)
		}
	}
}
