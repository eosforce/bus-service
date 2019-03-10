package basic

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/bronze1man/go-yaml2json"
	"github.com/eosforce/goeosforce"
)

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

func isStubABI(abi eos.ABI) bool {
	return abi.Version == "" &&
		abi.Actions == nil &&
		abi.ErrorMessages == nil &&
		abi.Extensions == nil &&
		abi.RicardianClauses == nil &&
		abi.Structs == nil && abi.Tables == nil &&
		abi.Types == nil
}
