package main

import (
    "log"
    "net"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
  //  pb "github.com/eosforce/forcegrpc/force_transfer"
    "google.golang.org/grpc/reflection"
	//"fmt"

	 "github.com/eosforce/bus-service/force-grpc-server/basic"
	// pb_trx "github.com/eosforce/forcegrpc/force_transaction"
	force_relay_commit "github.com/eosforce/bus-service/force_relay_commit"
	 "flag"
	 "github.com/eosforce/bus-service/force-grpc-server/common"
	// "strconv"
)

const (
    port = ":50051"
)

var vault_password = flag.String("vault_password", "123xyp", "the vault password")
var vault_file = flag.String("vault_file", "./eosc-vault.json", "the vault password")
var tcp_port = flag.Int("tcp_port",50051,"the port tcp listen")
var tcp_ip = flag.String("tcp_ip","127.0.0.1","the ip tcp listen")
var url = flag.String("url","http://127.0.0.1:8888","the addr which action send to")
var abipath = flag.String("abipath","./force.token.abi","the path of the abi file")
// server is used to implement helloworld.GreeterServer.
type server struct{}


func (s *server) RpcSendaction(ctx context.Context, in *force_relay_commit.RelayCommitRequest) (*force_relay_commit.RelayCommitReply, error) {
	//接下来解析transaction中的内容	
	basic.HandRelayBlock(in.Block,in.Action)
    return &force_relay_commit.RelayCommitReply{Reply:"get Block"}, nil
}

func main() {
    flag.Parse()

    lis,err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(*tcp_ip), *tcp_port, ""})
   // lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
	}
	
	common.SetVaultPasswd(*vault_password)
    common.SetVaultFile(*vault_file)
    common.SetDestUrl(*url)

    basic.SetAbiFilePath(*abipath)
	
    s := grpc.NewServer()
    force_relay_commit.RegisterRelayCommitServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
