package rabbitmq

import (
	"context"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func NewPublisherConn(ctx context.Context, mqUrl string) (*rabbitmqamqp.Publisher, error) {
	env := rabbitmqamqp.NewEnvironment(mqUrl, nil)
	mq, err := env.NewConnection(ctx)
	if err != nil {
		return nil, err
	}

	_, err = mq.Management().DeclareQueue(ctx, &rabbitmqamqp.QuorumQueueSpecification{Name: "book"})
	if err != nil {
		return nil, err
	}

	publisher, err := mq.NewPublisher(ctx, &rabbitmqamqp.QueueAddress{Queue: "book"}, nil)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
