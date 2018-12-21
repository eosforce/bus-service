package basic

import (
	//"fmt"
	pb_block "github.com/eosforce/forcegrpc/force_block"
)

func Handblock(blockNum int32,trans []*pb_block.BlockTransRequest) {
	
	for _, transValue := range trans {
		//先解析transfer
		HandTransaction(transValue.Trx,transValue.Trxid)
	}
}