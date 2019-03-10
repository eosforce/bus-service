package basic

import (
	"fmt"

	"github.com/fanyang1988/force-go"
)

// client client to force relay chain
var client *force.Client

// CreateClient create client to force relay chain
func CreateClient(configPath string) {
	var err error
	client, err = force.NewClientFromFile(configPath)
	if err != nil {
		fmt.Println("create client error  ", err.Error())
		return
	}
}

// GetLastCommittedBlock get last committed block to relay chain
func GetLastCommittedBlock() uint32 {

	return 0
}
