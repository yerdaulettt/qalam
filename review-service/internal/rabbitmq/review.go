package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"review-service/internal/review"

	"github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type reviewConsumer struct {
	consumer *rabbitmqamqp.Consumer
	service  *review.ReviewService
}

func NewReviewConsumer(c *rabbitmqamqp.Consumer, s *review.ReviewService) *reviewConsumer {
	return &reviewConsumer{consumer: c, service: s}
}

func (c *reviewConsumer) DeleteBookId(ctx context.Context) error {
	for {
		receive, err := c.consumer.Receive(ctx)
		if err != nil {
			log.Println(err)
		}

		var m review.Message

		err = json.Unmarshal(receive.Message().GetData(), &m)
		if err != nil {
			log.Println(err)
			continue
		}

		err = receive.Accept(ctx)
		if err != nil {
			log.Println(err)
		}

		err = c.service.DeleteBookId(ctx, m.BookId)
		if err != nil {
			log.Println(err)
		}
	}
}
