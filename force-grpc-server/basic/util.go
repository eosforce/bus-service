package basic

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/go-yaml2json"
	"github.com/eosforce/forcegrpc/force-grpc-server/common"
	"io/ioutil"
	"os"
	"strings"

	"github.com/eosforce/forcec/cli"
	eosvault "github.com/eosforce/forcec/vault"
	"github.com/eosforce/goeosforce"
)

func toSHA256Bytes(in, field string) eos.SHA256Bytes {
	return common.ToSHA256Bytes(in, field)
}

func pushEOSCActions(api *eos.API, actions ...*eos.Action) {
	common.PushEOSCActions(api, actions...)
}

func getEOSCTransaction(api *eos.API, actions ...*eos.Action) *eos.Transaction{
	return common.GetTransaction(api, actions...)
}


func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

func mustGetWallet() *eosvault.Vault {
	vault, err := common.SetupWallet()
	errorCheck("wallet setup", err)
	return vault
}

func permissionToPermissionLevel(in string) (out eos.PermissionLevel, err error) {
	return common.PermissionToPermissionLevel(in)
}

func permissionsToPermissionLevels(in []string) (out []eos.PermissionLevel, err error) {
	return common.PermissionsToPermissionLevels(in)
}

func getAPI() *eos.API {
	return common.GetAPI()
}

func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}


func loadYAMLOrJSONFile(filename string, v interface{}) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return json.Unmarshal(cnt, v)
	}
	return yamlUnmarshal(cnt, v)
}

func toName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid name format for %q", field), err)
	}

	return name
}

func toAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid account format for %q", field), err)
	}

	return acct
}

func toPermissionLevel(in, field string) eos.PermissionLevel {
	perm, err := common.PermissionToPermissionLevel(in)
	if err != nil {
		errorCheck(fmt.Sprintf("invalid permission level for %q", field), err)
	}
	return perm
}

func toActionName(in, field string) eos.ActionName {
	return eos.ActionName(toName(in, field))
}

func isStubABI(abi eos.ABI) bool {
	return abi.Version == "" &&
		abi.Actions == nil &&
		abi.ErrorMessages == nil &&
		abi.Extensions == nil &&
		abi.RicardianClauses == nil &&
		abi.Structs == nil && abi.Tables == nil &&
		abi.Types == nil
}

