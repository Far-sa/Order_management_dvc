package cosumer

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Far-sa/commons/api"
	"github.com/Far-sa/commons/broker"
	"github.com/Far-sa/order/contract"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service contract.OrderService
}

func NewConsumer(service contract.OrderService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
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
			log.Printf("Received message: %s", d.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			_, err := c.service.UpdateOrder(context.Background(), o)
			if err != nil {
				log.Printf("failed to update order: %v", err)
				continue
			}

			log.Println("Order has been updated from AMQP")
			d.Ack(false)
		}
	}()

	<-forever
}
