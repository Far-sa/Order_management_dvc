package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/commons/discovery"
	"github.com/Far-sa/commons/discovery/consul"
	"github.com/Far-sa/payment/adapter/consumer"
	"github.com/Far-sa/payment/adapter/gateway"
	stripeProcessor "github.com/Far-sa/payment/adapter/processor/stripe"
	"github.com/Far-sa/payment/handler"
	"github.com/Far-sa/payment/service"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	serviceName          = "payment"
	amqpUser             = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass             = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost             = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.EnvString("RABBITMQ_PORT", "5672")
	grpcAddr             = common.EnvString("GRPC_ADDRESS", "localhost:2001")
	consulAddr           = common.EnvString("CONSUL_ADDR", "localhost:8500")
	stripeKey            = common.EnvString("STRIPE_KEY", "")
	httpAddr             = common.EnvString("HTTP_ADDRESS", "localhost:8081")
	endpointStripeSecret = common.EnvString("STRIPE_ENDPOINT_SECRET", "whsec_...")
)

func main() {

	// Register consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)

	ctx := context.Background()
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatalf("failed to health check %v", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Unregister(ctx, instanceID, serviceName)

	//* stripe setup
	stripe.Key = stripeKey

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := stripeProcessor.NewProcessor()

	gateway := gateway.NewGateway(registry)
	paymentSvc := service.NewService(stripeProcessor, gateway)
	consumer := consumer.NewConsumer(paymentSvc)
	go consumer.Listen(ch)

	// http server
	mux := http.NewServeMux()
	httpServer := handler.NewPaymentHTTPHandler(ch)
	httpServer.RegisterRoutes(mux)

	go func() {
		log.Printf("Starting payment http server on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to start http payment server ")
		}
	}()

	// gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC Server Started at ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
