package consumer

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Far-sa/commons/api"

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

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := rc.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)

				continue
			}

			log.Printf("Payment link created %s", paymentLink)
		}
	}()

	<-forever
}
