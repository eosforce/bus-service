package basic

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	eos "github.com/eosforce/goeosforce"
	"github.com/tidwall/gjson"
)

var transferabi *eos.ABI

func SetAbiFilePath(path string) {
	b, err := ioutil.ReadFile(path)
	errorCheck("get ABI", err)
	str := string(b)
	r := strings.NewReader(str)
	transferabi, _ = eos.NewABI(r)
}

func DecodeTransfer(actionjson, trxid string) {
	hexData, _ := hex.DecodeString(actionjson)
	data, _ := transferabi.DecodeAction(hexData, eos.ActN("transfer"))
	from := gjson.GetBytes(data, "from").String()
	to := gjson.GetBytes(data, "to").String()
	quantity := gjson.GetBytes(data, "quantity").String()
	memo := gjson.GetBytes(data, "memo").String()
	fmt.Println("trxid:", trxid, ",from:", from, ",to:", to, ",quantity:", quantity, ",memo:", memo)
	//Transfer(from,to,quantity,memo,trxid)
}
