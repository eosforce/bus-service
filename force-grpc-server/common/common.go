package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	eosvault "github.com/eosforce/forcec/vault"
	eos "github.com/eosforce/goeosforce"
	"github.com/eosforce/goeosforce/ecc"
	"github.com/eosforce/goeosforce/sudo"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var global_vault_file string="/home/xuyapeng/go_workspace/src/github.com/eosforce/forcec/forcec/eosc-vault.json"

var vault_password string

func SetVaultPasswd(passwd string) {
	vault_password = passwd
}

func SetVaultFile (file string) {
	global_vault_file = file
}

func GetAPI() *eos.API {
	httpHeaders := viper.GetStringSlice("global-http-header")
	//viper.GetString("global-api-url")
	api := eos.New("http://127.0.0.1:8888")
	for _, header := range httpHeaders {
		headerArray := strings.SplitN(header, ": ", 2)
		if len(headerArray) != 2 || strings.Contains(headerArray[0], " ") {
			errorCheck("validating http headers", fmt.Errorf("invalid HTTP Header format"))
		}
		api.Header.Add(headerArray[0], headerArray[1])
	}
	return api
}

func SetupWallet() (*eosvault.Vault, error) {
	walletFile := global_vault_file//viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("wallet file %q missing: %s", walletFile, err)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault: %s", err)
	}
	boxer := eosvault.NewPassphraseBoxer(vault_password)
	if err != nil {
		return nil, fmt.Errorf("secret boxer: %s", err)
	}

	if err := vault.Open(boxer); err != nil {
		return nil, err
	}

	return vault, nil
}

func attachWallet(api *eos.API) {
	walletURLs := viper.GetStringSlice("global-wallet-url")
	if len(walletURLs) == 0 {
		vault, err := SetupWallet()
		errorCheck("setting up wallet", err)

		api.SetSigner(vault.KeyBag)
	} else {
		if len(walletURLs) == 1 {
			// If a `walletURLs` has a Username in the path, use instead of `default`.
			api.SetSigner(eos.NewWalletSigner(eos.New(walletURLs[0]), "default"))
		} else {
			fmt.Println("Multi-signer not yet implemented.  Please choose only one `--wallet-url`")
			os.Exit(1)
		}
	}
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

func PermissionToPermissionLevel(in string) (out eos.PermissionLevel, err error) {
	return eos.NewPermissionLevel(in)
}

func PermissionsToPermissionLevels(in []string) (out []eos.PermissionLevel, err error) {
	// loop all parameters
	for _, singleArg := range in {

		// if they specified "account@active,account2", handle that too..
		for _, val := range strings.Split(singleArg, ",") {
			level, err := PermissionToPermissionLevel(strings.TrimSpace(val))
			if err != nil {
				return out, err
			}

			out = append(out, level)
		}
	}

	return
}

func GetTransaction(api *eos.API, actions ...*eos.Action) *eos.Transaction{
	opts := &eos.TxOptions{}

	if chainID := viper.GetString("global-offline-chain-id"); chainID != "" {
		opts.ChainID = ToSHA256Bytes(chainID, "--offline-chain-id")
	}

	if headBlockID := viper.GetString("global-offline-head-block"); headBlockID != "" {
		opts.HeadBlockID = ToSHA256Bytes(headBlockID, "--offline-head-block")
	}

	if delaySec := viper.GetInt("global-delay-sec"); delaySec != 0 {
		opts.DelaySecs = uint32(delaySec)
	}

	if err := opts.FillFromChain(api); err != nil {
		fmt.Println("Error fetching tapos + chain_id from the chain (specify --offline flags for offline operations):", err)
		os.Exit(1)
	}

	tx := eos.NewTransaction(actions, opts)

	tx = optionallySudoWrap(tx, opts)

	tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)

	fee, err := GetFeeByTrx(tx)
	if err != nil {
		fmt.Println("Error get fee:", err)
		os.Exit(1)
	}
	tx.Fee = fee
	return tx
}

func PushEOSCActions(api *eos.API, actions ...*eos.Action) {
	PushEOSCActionsAndContextFreeActions(api, nil, actions)
}

func PushEOSCActionsAndContextFreeActions(api *eos.API, contextFreeActions []*eos.Action, actions []*eos.Action) {
	for _, act := range contextFreeActions {
		act.Authorization = nil
	}
	permissions := viper.GetStringSlice("global-permission")
	if len(permissions) != 0 {
		levels, err := PermissionsToPermissionLevels(permissions)
		errorCheck("specified --permission(s) invalid", err)

		for _, act := range actions {
			act.Authorization = levels
		}
	}
	opts := &eos.TxOptions{}

	if chainID := viper.GetString("global-offline-chain-id"); chainID != "" {
		opts.ChainID = ToSHA256Bytes(chainID, "--offline-chain-id")
	}

	if headBlockID := viper.GetString("global-offline-head-block"); headBlockID != "" {
		opts.HeadBlockID = ToSHA256Bytes(headBlockID, "--offline-head-block")
	}

	if delaySec := viper.GetInt("global-delay-sec"); delaySec != 0 {
		opts.DelaySecs = uint32(delaySec)
	}

	if err := opts.FillFromChain(api); err != nil {
		fmt.Println("Error fetching tapos + chain_id from the chain (specify --offline flags for offline operations):", err)
		os.Exit(1)
	}
	
	tx := eos.NewTransaction(actions, opts)
	if len(contextFreeActions) > 0 {
		tx.ContextFreeActions = contextFreeActions
	}
	tx = optionallySudoWrap(tx, opts)

	tx.SetExpiration(time.Duration(120) * time.Second)
	fee, err := GetFeeByTrx(tx)
	if err != nil {
		fmt.Println("Error get fee:", err)
		os.Exit(1)
	}
	tx.Fee = fee
	signedTx, packedTx := OptionallySignTransaction(tx, opts.ChainID, api)
	OptionallyPushTransaction(signedTx, packedTx, opts.ChainID, api)
}

func optionallySudoWrap(tx *eos.Transaction, opts *eos.TxOptions) *eos.Transaction {
	if viper.GetBool("global-sudo-wrap") {
		return eos.NewTransaction([]*eos.Action{sudo.NewExec(eos.AccountName("eosio"), *tx)}, opts)
	}
	return tx
}

func OptionallySignTransaction(tx *eos.Transaction, chainID eos.SHA256Bytes, api *eos.API) (signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction) {
	if !viper.GetBool("global-skip-sign") {
		textSignKeys := viper.GetStringSlice("global-offline-sign-key")
		if len(textSignKeys) > 0 {
			var signKeys []ecc.PublicKey
			for _, key := range textSignKeys {
				pubKey, err := ecc.NewPublicKey(key)
				errorCheck(fmt.Sprintf("parsing public key %q", key), err)

				signKeys = append(signKeys, pubKey)
			}
			api.SetCustomGetRequiredKeys(func(tx *eos.Transaction) ([]ecc.PublicKey, error) {
				return signKeys, nil
			})
		}

		attachWallet(api)

		var err error
		signedTx, packedTx, err = api.SignTransaction(tx, chainID, eos.CompressionNone)
		errorCheck("signing transaction", err)
	} else {
		signedTx = eos.NewSignedTransaction(tx)
	}
	return signedTx, packedTx
}

func OptionallyPushTransaction(signedTx *eos.SignedTransaction, packedTx *eos.PackedTransaction, chainID eos.SHA256Bytes, api *eos.API) {
	writeTrx := viper.GetString("global-write-transaction")

	if writeTrx != "" {
		cnt, err := json.MarshalIndent(signedTx, "", "  ")
		errorCheck("marshalling json", err)

		annotatedCnt, err := sjson.Set(string(cnt), "chain_id", hex.EncodeToString(chainID))
		errorCheck("adding chain_id", err)

		err = ioutil.WriteFile(writeTrx, []byte(annotatedCnt), 0644)
		errorCheck("writing output transaction", err)

		fmt.Printf("Transaction written to %q\n", writeTrx)
	} else {
		if packedTx == nil {
			fmt.Println("A signed transaction is required if you want to broadcast it. Remove --skip-sign (or add --write-transaction ?)")
			os.Exit(1)
		}

		// TODO: print the traces
		PushTransaction(api, packedTx, chainID)
	}
}

func PushTransaction(api *eos.API, packedTx *eos.PackedTransaction, chainID eos.SHA256Bytes) {
	resp, err := api.PushTransaction(packedTx)
	errorCheck("pushing transaction", err)

	//fmt.Println("Transaction submitted to the network. Confirm at https://eosq.app/tx/" + resp.TransactionID)
	trxURL := transactionURL(chainID, resp.TransactionID)
	fmt.Printf("\nTransaction submitted to the network.\n  %s\n", trxURL)
	if resp.BlockID != "" {
		blockURL := blockURL(chainID, resp.BlockID)
		fmt.Printf("Server says transaction was included in block %d:\n  %s\n", resp.BlockNum, blockURL)
	}
}

func transactionURL(chainID eos.SHA256Bytes, trxID string) string {
	hexChain := hex.EncodeToString(chainID)
	switch hexChain {
	case "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906":
		return fmt.Sprintf("https://eosq.app/tx/%s", trxID)
	case "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191":
		return fmt.Sprintf("https://kylin.eosq.app/tx/%s", trxID)
	}
	return trxID
}

func blockURL(chainID eos.SHA256Bytes, blockID string) string {
	hexChain := hex.EncodeToString(chainID)
	switch hexChain {
	case "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906":
		return fmt.Sprintf("https://eosq.app/block/%s", blockID)
	case "5fff1dae8dc8e2fc4d5b23b2c7665c97f9e9d8edf2b6485a86ba311c25639191":
		return fmt.Sprintf("https://kylin.eosq.app/block/%s", blockID)
	}
	return blockID
}

func ToSHA256Bytes(in, field string) eos.SHA256Bytes {
	if len(in) != 64 {
		errorCheck(fmt.Sprintf("%q invalid", field), errors.New("should be 64 hexadecimal characters"))
	}

	bytes, err := hex.DecodeString(in)
	errorCheck(fmt.Sprintf("invalid hex in %q", field), err)

	return bytes
}
