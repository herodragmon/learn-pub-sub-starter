package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"context"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: b,
	})
	if err != nil{
		return err
	}
	return nil
}

type SimpleQueueType string

const (
    QueueTransient SimpleQueueType = "transient"
    QueueDurable   SimpleQueueType = "durable"
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err:= conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var durable, autoDelete, exclusive bool

	switch queueType {
	case QueueDurable:
  	durable = true
    autoDelete = false
    exclusive = false
	case QueueTransient:
    durable = false
    autoDelete = true
    exclusive = true
	}

	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, err
	}
	return ch, q, nil
}
