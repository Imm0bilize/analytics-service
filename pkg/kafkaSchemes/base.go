package kafkaSchemes

type BaseTopic struct {
	TypeTopic      string `json:"type"`
	Timestamp      string `json:"timestamp"`
	IdempotencyKey string `json:"idempotency_key"`
}
