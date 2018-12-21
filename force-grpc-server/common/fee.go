package common

import (
	"fmt"
	"github.com/eosforce/goeosforce"
)

// GetFeeByTrx get fee sum by actions
func GetFeeByTrx(tx *eos.Transaction) (eos.Asset, error) {
	api := GetAPI()

	// if no set will err
	tx.Fee = eos.NewEOSAsset(0)

	resp, err := api.GetFee(tx)
	fmt.Printf("resp %v\n", resp)
	if err != nil{
		return eos.NewEOSAsset(0), err
	}
	return resp.Fee, nil
}
