package basic

import (
	//"fmt"
	"time"
	eos "github.com/eosforce/goeosforce"
	"github.com/eosforce/goeosforce/system"
	"github.com/tidwall/gjson"
	"encoding/hex"
	//"github.com/eosforce/forcegrpc/force-grpc-server/common"
)

func Transfer(from,to,amount,memo,trxid string) {
	quantity, err := eos.NewEOSAssetFromString("1.0000 EOS")
	errorCheck("invalid amount", err)

	account_from := toAccount("eosforce","Transfer.go toaccount")
	account_to := toAccount("biosbpa","Transfer.go toaccount")
	
	api := getAPI()
	memo = "from:"+from+" to:"+to+" amount:"+amount+" trxid:" + trxid+" time"+time.Now().Format("2006/1/2 15:04:05")
	action := system.NewTransfer(account_from, account_to, quantity, memo)
	pushEOSCActions(api, action)
}

func DecodeTransfer(actionjson,trxid string) {
	accountName := toAccount("eosio","DecodeTransfer")
	action_name := toActionName("transfer","DecodeTransfer")
	api := getAPI()
	abi, err := api.GetABI(accountName)
	errorCheck("get ABI", err)
	hexData,_ := hex.DecodeString(actionjson)
	data, _ := abi.ABI.DecodeAction( hexData,action_name)
	from := gjson.GetBytes(data, "from").String()
	to := gjson.GetBytes(data, "to").String()
	quantity := gjson.GetBytes(data, "quantity").String()
	memo := gjson.GetBytes(data, "memo").String()
	Transfer(from,to,quantity,memo,trxid)
}