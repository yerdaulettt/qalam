package rabbitmq

import (
	"context"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

func NewConsumerConn(ctx context.Context, mqUrl string) (*rabbitmqamqp.Consumer, error) {
	env := rabbitmqamqp.NewEnvironment(mqUrl, nil)
	mq, err := env.NewConnection(ctx)
	if err != nil {
		return nil, err
	}

	_, err = mq.Management().DeclareQueue(ctx, &rabbitmqamqp.QuorumQueueSpecification{Name: "book"})
	if err != nil {
		return nil, err
	}

	consumer, err := mq.NewConsumer(ctx, "book", nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
