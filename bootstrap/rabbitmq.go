package bootstrap

import (
	domain "github.com/Psykepro/item-storage-client/_domain"
	"github.com/Psykepro/item-storage-client/config"
	"github.com/Psykepro/item-storage-client/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

func RabbitMq(cfg *config.Config, logger domain.Logger) (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	mqConn, err := rabbitmq.NewRabbitMQConn(cfg.RabbitMQ)
	if err != nil {
		logger.Fatalf("Failed to init RabbitMQ Connection")
	}
	mqChannel, mqQueue, err := rabbitmq.InitRabbitMQConn(mqConn, cfg.RabbitMQ, logger)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	return mqConn, mqChannel, mqQueue
}
