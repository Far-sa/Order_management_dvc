package discovery

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		log.Printf("Error discovering service %s: %v", serviceName, err)
		return nil, err
	}

	log.Printf("Discovered addresses for service %s: %v", serviceName, addrs)

	if len(addrs) == 0 {
		return nil, fmt.Errorf("no instances of service %s found", serviceName)
	}

	// Randomly select an instance
	return grpc.Dial(
		addrs[rand.Intn(len(addrs))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//* grpc middlewares
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}
