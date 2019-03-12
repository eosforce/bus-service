package main

import (
	"flag"

	"github.com/cihub/seelog"
	"github.com/eosforce/goforceio/ecc"
)

var transferURL = flag.String("url", "0.0.0.0:50051", "transfer service url to listen")
var relayURL = flag.String("relay-url", "0.0.0.0:50052", "relay service url to listen")
var configPath = flag.String("cfg", "./config.json", "confg file path")
var chain = flag.String("chain", "eosforce", "the name of chain")
var transfer = flag.String("transfer", "eosforce", "the name of transfer")

func init() {
	ecc.PublicKeyPrefixCompat = "FOSC"
}

func main() {
	flag.Parse()
	defer seelog.Flush()

	go func() {
		startRelayService()
	}()
	// start service for side chain
	startSideService()
}
