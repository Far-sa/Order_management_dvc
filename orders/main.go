package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/commons/discovery"
	"github.com/Far-sa/commons/discovery/consul"
	"github.com/Far-sa/order/handler"
	"github.com/Far-sa/order/repository"
	"github.com/Far-sa/order/service"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Unregister(ctx, instanceID, serviceName)

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.New()
	svc := service.New(repo)

	handler.NewGRPC(grpcServer, svc)
	svc.CreateOrder(context.Background())

	log.Println("GRPC server started at:", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
