package rabbitmq

import (
	"context"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type rabbitConn struct {
	env        *rabbitmqamqp.Environment
	conn       *rabbitmqamqp.AmqpConnection
	Management *rabbitmqamqp.AmqpManagement
}

func NewConn(ctx context.Context, mqUrl string) (*rabbitConn, error) {
	env := rabbitmqamqp.NewEnvironment(mqUrl, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		return nil, err
	}

	return &rabbitConn{
		env:        env,
		conn:       conn,
		Management: conn.Management(),
	}, nil
}

func (r *rabbitConn) NewPublisher(ctx context.Context) (*rabbitmqamqp.Publisher, error) {
	_, err := r.Management.DeclareExchange(ctx, &rabbitmqamqp.TopicExchangeSpecification{Name: "events"})
	if err != nil {
		return nil, err
	}

	p, err := r.conn.NewPublisher(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	return p, nil
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
