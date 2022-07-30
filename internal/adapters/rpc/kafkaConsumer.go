package rpc

import (
	"analytic-service/internal/ports"
	"analytic-service/pkg/kafka"
	"analytic-service/pkg/kafkaSchemes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"os"
)

type KafkaConsumer struct {
	kafka.ConsumerConnection
	domain        ports.CommandDomain
	idempKeyStore ports.IdempotencyKeyStorage
	logger        logrus.FieldLogger
	notify        chan error
	topic         string
}

func NewKafkaConsumer(
	logger logrus.FieldLogger,
	domain ports.CommandDomain,
	idempotencyKeyValidator ports.IdempotencyKeyStorage,
	consumer kafka.ConsumerConnection,
	topic string,
) *KafkaConsumer {
	notify := make(chan error, 1)
	return &KafkaConsumer{
		ConsumerConnection: consumer,
		domain:             domain,
		idempKeyStore:      idempotencyKeyValidator,
		logger:             logger,
		notify:             notify,
		topic:              topic,
	}
}

func (k *KafkaConsumer) checkIdempotencyKey(ctx context.Context, key string) error {
	ok, err := k.idempKeyStore.CheckIdempotencyKeyInStore(ctx, key)
	if err != nil {
		k.logger.Errorf("failed to check idempotency key in storage: %s", err.Error())
		return err
	}
	if ok {
		k.logger.Infof("task with idempotency key {%s} has already been processed, skip", key)
		return errors.New("key exists in storage")
	}
	return nil
}

func (k *KafkaConsumer) onCreateCommand(ctx context.Context, rawMessage *[]byte) error {
	var payload kafkaSchemes.TaskAnalyticsCreateTypePayload
	if err := json.Unmarshal(*rawMessage, &payload); err != nil {
		k.logger.Errorf("error: <%s> while unmarshalling payload, message header: %v", err.Error())
		return err
	}
	if err := k.domain.CreateTask(ctx, payload.TaskID); err != nil {
		return err
	}
	return nil
}

func (k *KafkaConsumer) onSetStartCommand(ctx context.Context, rawMessage *[]byte) error {
	var payload kafkaSchemes.TaskAnalyticsAcceptRejectTypePayload
	if err := json.Unmarshal(*rawMessage, &payload); err != nil {
		k.logger.Errorf("error: <%s> while unmarshalling payload", err.Error())
		return err
	}
	if err := k.domain.SetTimeStart(ctx, payload.TaskID, payload.Email, payload.Time, payload.TaskState); err != nil {
		return err
	}
	return nil
}

func (k *KafkaConsumer) onSetEndCommand(ctx context.Context, rawMessage *[]byte) error {
	var payload kafkaSchemes.TaskAnalyticsAcceptRejectTypePayload
	if err := json.Unmarshal(*rawMessage, &payload); err != nil {
		k.logger.Errorf("error: <%s> while unmarshalling payload", err.Error())
		return err
	}
	if err := k.domain.SetTimeEnd(ctx, payload.TaskID, payload.Email, payload.Time, payload.TaskState); err != nil {
		return err
	}
	return nil
}

func (k *KafkaConsumer) selectCommand(rawMessage []byte) {
	var message kafkaSchemes.BaseTopic
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		k.logger.Errorf("error while unmarshalling rawMessage: %v", err)
		return
	}

	ctx := context.TODO()
	if err := k.checkIdempotencyKey(ctx, message.IdempotencyKey); err != nil {
		return
	}

	switch message.TypeTopic {
	case "create_task":
		if err := k.onCreateCommand(ctx, &rawMessage); err != nil {
			return
		}
	case "set_start":
		if err := k.onSetStartCommand(ctx, &rawMessage); err != nil {
			return
		}
	case "set_end":
		if err := k.onSetEndCommand(ctx, &rawMessage); err != nil {
			return
		}
	default:
		k.logger.Errorf("can't select method, got - %s", message.TypeTopic)
		return
	}
	if err := k.idempKeyStore.Commit(ctx, message.IdempotencyKey); err != nil {
		k.logger.Errorf("error committing message by idempotency key: %s", err.Error())
	}
}

func (k *KafkaConsumer) Run(interruptChan chan os.Signal) {
	go func() {
		consumer, err := k.C.ConsumePartition(k.topic, 0, sarama.OffsetNewest)
		if err != nil {
			k.notify <- err
			close(k.notify)
		}

		for true {
			select {
			case msg := <-consumer.Messages():
				k.selectCommand(msg.Value)
			case err := <-consumer.Errors():
				k.logger.Errorf("error in topic: %s", err.Error())
			case <-interruptChan:
				err := consumer.Close()
				if err != nil {
					k.notify <- err
					close(k.notify)
				}
				return
			}
		}
	}()
}

func (k *KafkaConsumer) Notify() <-chan error {
	return k.notify
}
