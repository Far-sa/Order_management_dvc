package main

import (
	"context"
	"log"
	"net"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/order/handler"
	"github.com/Far-sa/order/repository"
	"github.com/Far-sa/order/service"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
)

func main() {

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("error listening : %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	repo := repository.New()
	svc := service.New(repo)

	handler.NewGRPC(grpcServer, svc)
	svc.CreateOrder(context.Background())

	log.Println("GRPC server started at:", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
