package main

import (
    "log"
    "net"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
  //  pb "github.com/eosforce/forcegrpc/force_transfer"
    "google.golang.org/grpc/reflection"
//	"fmt"

	 "github.com/eosforce/forcegrpc/force-grpc-server/basic"
	// pb_trx "github.com/eosforce/forcegrpc/force_transaction"
	pb_block "github.com/eosforce/forcegrpc/force_block"
	 "flag"
	 "github.com/eosforce/forcegrpc/force-grpc-server/common"
	 "strconv"
)

const (
    port = ":50051"
)

var vault_password = flag.String("vault_password", "123xyp", "the vault password")
var vault_file = flag.String("vault_file", "./eosc-vault.json", "the vault password")

// server is used to implement helloworld.GreeterServer.
type server struct{}


func (s *server) RpcSendaction(ctx context.Context, in *pb_block.BlockRequest) (*pb_block.BlockReply, error) {
	//接下来解析transaction中的内容	
	basic.Handblock(in.Blocknum,in.Trans)
    return &pb_block.BlockReply{Reply:"get Block",Message: "BlockNum:"+strconv.Itoa(int(in.Blocknum))}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
	}
	flag.Parse()
	
	common.SetVaultPasswd(*vault_password)
	common.SetVaultFile(*vault_file)
	
    s := grpc.NewServer()
    pb_block.RegisterGrpcBlockServer(s, &server{})
    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
