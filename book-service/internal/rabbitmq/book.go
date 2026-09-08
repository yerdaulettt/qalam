package rabbitmq

import (
	"context"
	"errors"
	"log"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type adminPublisher struct {
	publisher  *rabbitmqamqp.Publisher
	management *rabbitmqamqp.AmqpManagement
}

func NewAdminPublisher(p *rabbitmqamqp.Publisher, m *rabbitmqamqp.AmqpManagement) *adminPublisher {
	return &adminPublisher{publisher: p, management: m}
}

func (p *adminPublisher) Publish(ctx context.Context, message []byte, queue string) error {
	_, err := p.management.DeclareQueue(ctx, &rabbitmqamqp.QuorumQueueSpecification{Name: queue})
	if err != nil {
		return err
	}

	msg, err := rabbitmqamqp.NewMessageWithAddress(message, &rabbitmqamqp.QueueAddress{Queue: queue})
	if err != nil {
		return err
	}

	res, err := p.publisher.Publish(ctx, msg)
	if err != nil {
		return err
	}

	switch res.Outcome.(type) {
	case *rabbitmqamqp.StateAccepted:
		return nil
	case *rabbitmqamqp.StateRejected:
		return res.MessageRejectedError
	default:
		log.Println(res.Outcome)
		return errors.New("Publish error")
	}
}
