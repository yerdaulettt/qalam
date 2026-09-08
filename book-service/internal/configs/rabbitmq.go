package configs

import "fmt"

func NewRabbitUrl() string {
	user := getEnv("RMQ_USER", "guest")
	password := getEnv("RMQ_PASSWORD", "guest")
	host := getEnv("RMQ_HOST", "localhost")
	port := getEnv("RMQ_PORT", "5672")

	mqUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	return mqUrl
}
