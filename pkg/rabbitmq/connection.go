package rabbitmq

import (
	"errors"
	"fmt"

	domain "github.com/Psykepro/item-storage-client/_domain"

	"github.com/Psykepro/item-storage-client/config"

	"github.com/streadway/amqp"
)

// NewRabbitMQConn Initializing new RabbitMQ connection
func NewRabbitMQConn(cfg *config.RabbitMQ) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	return amqp.Dial(connAddr)
}

func InitRabbitMQConn(
	amqpConn *amqp.Connection,
	rabbitMqCfg *config.RabbitMQ,
	logger domain.Logger,
) (*amqp.Channel, *amqp.Queue, error) {

	ch, err := openRabbitMqChannel(amqpConn, logger)
	if err != nil {
		return nil, nil, err
	}

	q, err := declareRabbitMqQueue(ch, rabbitMqCfg, logger)
	if err != nil {
		return ch, nil, err
	}

	return ch, q, nil
}

func openRabbitMqChannel(amqpConn *amqp.Connection, logger domain.Logger) (*amqp.Channel, error) {
	logger.Debugf("Opening RabbitMQ channel ...")
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open a channel. Err: [%s]", err))
	}
	logger.Debugf("Successfully opened RabbitMQ channel!")

	return ch, nil
}

func declareRabbitMqQueue(ch *amqp.Channel, rabbitMqCfg *config.RabbitMQ, logger domain.Logger) (*amqp.Queue, error) {
	logger.Debugf("Declaring RabbitMQ Queue ...")
	q, err := ch.QueueDeclare(
		rabbitMqCfg.Queue.Name,
		rabbitMqCfg.Queue.Durable,
		rabbitMqCfg.Queue.AutoDelete,
		rabbitMqCfg.Queue.Exclusive,
		rabbitMqCfg.Queue.NoWait,
		rabbitMqCfg.Queue.Args,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to declare RabbitMQ Queue. Err: [%s]", err))
	}
	logger.Debugf("Successfully declaring RabbitMQ Queue!")

	return &q, nil
}
