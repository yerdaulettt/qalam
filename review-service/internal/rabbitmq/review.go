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

		message := receive.Message()

		switch message.Annotations["x-routing-key"] {
		case "book.deleted":
			var m review.Message

			err = json.Unmarshal(message.GetData(), &m)
			if err != nil {
				log.Println(err)
				continue
			}

			err = c.service.DeleteBookId(ctx, m.BookId)
			if err != nil {
				log.Println(err)
			}
		case "book.updated":
			var bm review.BookMessage

			err = json.Unmarshal(message.GetData(), &bm)
			if err != nil {
				log.Println(err)
				continue
			}

			err = c.service.UpdateBook(ctx, bm.BookId, bm.Name)
			if err != nil {
				log.Println(err)
			}
		}

		err = receive.Accept(ctx)
		if err != nil {
			log.Println(err)
		}
	}
}
