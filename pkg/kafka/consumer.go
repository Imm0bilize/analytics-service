package kafka

import (
	"github.com/Shopify/sarama"
)

type ConsumerConnection struct {
	C sarama.Consumer
}

// NewConsumerConnection creates a new kafka consumer
// hostPortPairs in format: "0.0.0.0:2345", ...
func NewConsumerConnection(clientId string, hostPortPairs ...string) (*ConsumerConnection, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientId
	cfg.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(hostPortPairs, cfg)
	if err != nil {
		return nil, err
	}
	return &ConsumerConnection{consumer}, nil
}

func (c *ConsumerConnection) Shutdown() error {
	return c.C.Close()
}
