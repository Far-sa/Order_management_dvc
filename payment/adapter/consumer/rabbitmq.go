package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/Far-sa/commons/api"
	"go.opentelemetry.io/otel"

	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/payment/contract"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service contract.PaymentsService
}

func NewConsumer(service contract.PaymentsService) *consumer {
	return &consumer{service}
}

func (rc *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("received message: %s", d.Body)

			// Extract the headers
			ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := rc.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)

				if err := broker.HandleDelivery(ch, &d); err != nil {
					log.Printf("error handling retry %v", err)
				}

				d.Nack(false, false)

				continue
			}

			messageSpan.AddEvent(fmt.Sprintf("payment.created: %s", paymentLink))
			messageSpan.End()

			log.Printf("Payment link created %s", paymentLink)
			d.Ack(false)
		}
	}()

	<-forever
}
