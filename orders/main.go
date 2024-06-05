package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/commons/discovery"
	"github.com/Far-sa/commons/discovery/consul"
	"github.com/Far-sa/commons/tracer"
	consumer "github.com/Far-sa/order/cosumer"
	"github.com/Far-sa/order/handler"
	"github.com/Far-sa/order/repository"
	"github.com/Far-sa/order/service"
	"github.com/Far-sa/order/telemetry"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr  = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := tracer.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
	}

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
				logger.Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Unregister(ctx, instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	repo := repository.New()
	svc := service.New(repo)

	//! decorator pattern- useful for metric,tracing and logging
	svcWithTelemetry := telemetry.NewTelemetryMiddleware(svc)
	svcWithLogging := telemetry.NewTelemetryMiddleware(svcWithTelemetry)

	handler.NewGRPC(grpcServer, svcWithLogging, ch)

	consumer := consumer.NewConsumer(svcWithLogging)
	go consumer.Listen(ch)

	logger.Info("Starting HTTP server", zap.String("port", grpcAddr))

	log.Println("GRPC server started at:", grpcAddr)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
