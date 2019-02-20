package main

import (
    "log"
    "net"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
	//"fmt"

	 "github.com/eosforce/bus-service/force-grpc-server/basic"
	force_relay_commit "github.com/eosforce/bus-service/force_relay_commit"
	 "flag"
     "github.com/eosforce/goeosforce/ecc"
)

const (
    port = ":50051"
)

var tcp_port = flag.Int("tcp_port",50051,"the port tcp listen")
var tcp_ip = flag.String("tcp_ip","127.0.0.1","the ip tcp listen")
var configPath = flag.String("cfg", "./config.json", "confg file path")
var chain = flag.String("chain name","eosforce","the name of chain")
var transfer = flag.String("transfer name","eosforce","the name of transfer")

// server is used to implement helloworld.GreeterServer.
type server struct{}

func init() {
	ecc.PublicKeyPrefixCompat = "FOSC"
}

func (s *server) RpcSendaction(ctx context.Context, in *force_relay_commit.RelayCommitRequest) (*force_relay_commit.RelayCommitReply, error) {
	//接下来解析transaction中的内容	
	basic.HandRelayBlock(in.Block,in.Action)
    return &force_relay_commit.RelayCommitReply{Reply:"get Block"}, nil
}

func main() {
    flag.Parse()

    lis,err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(*tcp_ip), *tcp_port, ""})
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
	}
	
    basic.Createclient(*configPath)
    basic.SetChain(*chain)
    basic.SetTransfer(*transfer)
	
    s := grpc.NewServer()
    force_relay_commit.RegisterRelayCommitServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
