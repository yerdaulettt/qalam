package rabbitmq

import (
	"context"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type rabbitConn struct {
	env  *rabbitmqamqp.Environment
	conn *rabbitmqamqp.AmqpConnection
}

func NewConn(ctx context.Context, mqUrl string) (*rabbitConn, error) {
	env := rabbitmqamqp.NewEnvironment(mqUrl, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		return nil, err
	}

	return &rabbitConn{
		env:  env,
		conn: conn,
	}, nil
}

func (r *rabbitConn) Close(ctx context.Context) error {
	err := r.env.CloseConnections(ctx)
	if err != nil {
		return err
	}

	err = r.conn.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitConn) NewConsumer(ctx context.Context, bindKey, queue string) (*rabbitmqamqp.Consumer, error) {
	topic, err := r.conn.Management().DeclareExchange(ctx, &rabbitmqamqp.TopicExchangeSpecification{Name: "events"})
	if err != nil {
		return nil, err
	}

	_, err = r.conn.Management().DeclareQueue(ctx, &rabbitmqamqp.QuorumQueueSpecification{Name: queue})
	if err != nil {
		return nil, err
	}

	_, err = r.conn.Management().Bind(ctx, &rabbitmqamqp.ExchangeToQueueBindingSpecification{
		SourceExchange:   topic.Name(),
		DestinationQueue: queue,
		BindingKey:       bindKey,
	})
	if err != nil {
		return nil, err
	}

	c, err := r.conn.NewConsumer(ctx, queue, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}
