package basic

import (
	"fmt"
	force_relay_commit "github.com/eosforce/bus-service/force_relay_commit"
)

// func Handblock(blockNum int32,trans []*pb_block.BlockTransRequest) {
// 	for _, transValue := range trans {
// 		//先解析transfer
// 		HandTransaction(transValue.Trx,transValue.Trxid)
// 	}
// }
//处理先关的块信息
//
func HandRelayBlock(block *force_relay_commit.RelayBlock,Action []*force_relay_commit.RelayAction) {
	fmt.Println(block.Producer,"--",block.Id,"--",block.Previous,"--",block.Confirmed,"--",block.TransactionMroot,"--",block.ActionMroot,"--",block.Mroot,"--")
	for _, ActionValue := range Action  {
		fmt.Println("before print Action info -------------------")
		fmt.Println(ActionValue.Account,"---",ActionValue.ActionName,"---",ActionValue.Data)
	}
}