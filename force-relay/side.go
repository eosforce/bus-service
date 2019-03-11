package main

import (
	"context"
	"net"

	"github.com/cihub/seelog"
	commit "github.com/eosforce/bus-service/force-relay/pbs/relay"
	"github.com/eosforce/bus-service/force-relay/side"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type sideServer struct{}

func (s *sideServer) RpcSendaction(ctx context.Context, in *commit.RelayCommitRequest) (*commit.RelayCommitReply, error) {
	side.HandRelayBlock(in.Block, in.Action)
	return &commit.RelayCommitReply{Reply: "get Block"}, nil
}

func startSideService() {
	lis, err := net.Listen("tcp", *transferURL)
	if err != nil {
		seelog.Errorf("failed to listen: %v", err)
		return
	}

	side.CreateClient(*configPath)
	relayCfg := side.NewCfg(*chain, *transfer)
	relayCfg.AppendActionInfo("force.token", "transfer")
	relayCfg.AppendActionInfo("force", "newaccount")
	side.SetCfg(relayCfg)

	service := grpc.NewServer()
	commit.RegisterRelayCommitServer(service, &sideServer{})
	reflection.Register(service)
	if err := service.Serve(lis); err != nil {
		seelog.Errorf("failed to serve: %v", err.Error())
	}
}
