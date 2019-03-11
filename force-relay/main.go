package main

import (
	"flag"
	"log"
	"net"

	force_relay_commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"github.com/eosforce/goeosforce/ecc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var transferURL = flag.String("url", "0.0.0.0:50051", "transfer service url to listen")
var configPath = flag.String("cfg", "./config.json", "confg file path")
var chain = flag.String("chain", "eosforce", "the name of chain")
var transfer = flag.String("transfer", "eosforce", "the name of transfer")

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	ecc.PublicKeyPrefixCompat = "FOSC"
}

func (s *server) RpcSendaction(ctx context.Context, in *force_relay_commit.RelayCommitRequest) (*force_relay_commit.RelayCommitReply, error) {
	side.HandRelayBlock(in.Block, in.Action)
	return &force_relay_commit.RelayCommitReply{Reply: "get Block"}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *transferURL)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	side.CreateClient(*configPath)

	relayCfg := side.NewCfg(*chain, *transfer)
	relayCfg.AppendActionInfo("force.token", "transfer")
	relayCfg.AppendActionInfo("force", "newaccount")
	side.SetCfg(relayCfg)

	s := grpc.NewServer()
	force_relay_commit.RegisterRelayCommitServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
