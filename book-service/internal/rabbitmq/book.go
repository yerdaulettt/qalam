package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"book-service/internal/book"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type adminPublisher struct {
	publisher *rabbitmqamqp.Publisher
}

func NewAdminPublisher(p *rabbitmqamqp.Publisher) *adminPublisher {
	return &adminPublisher{publisher: p}
}

func (p *adminPublisher) Publish(ctx context.Context, m book.Message) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	res, err := p.publisher.Publish(ctx, rabbitmqamqp.NewMessage(data))
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
