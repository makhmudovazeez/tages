package main

import (
	"flag"
	"fmt"
	"github.com/makhmudovazeez/tages/internal/config"
	"github.com/makhmudovazeez/tages/internal/server"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configFile = flag.String("f", "cmd/tages.yaml", "the config file")

func main() {
	svcCtx := svc.NewServiceContext(*config.LoadConfig(*configFile))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", svcCtx.Config.Port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	ts := server.NewTagesServer(svcCtx)
	tages.RegisterTagesServer(s, ts)

	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
