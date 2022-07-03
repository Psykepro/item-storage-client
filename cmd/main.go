package main

import (
	"github.com/Psykepro/item-storage-client/bootstrap"
	"github.com/Psykepro/item-storage-client/internal/request"

	"github.com/Psykepro/item-storage-client/pkg/rabbitmq"
)

func main() {
	cfg := bootstrap.Config()

	logger := bootstrap.Logger(cfg)
	requestLoader := request.NewLoader(logger)
	requests := requestLoader.LoadRequestsFromCsv(request.DataPath)
	mqConn, mqChannel, mqQueue := bootstrap.RabbitMq(cfg, logger)
	defer mqChannel.Close()
	defer mqConn.Close()
	publisher := rabbitmq.NewPublisher(mqConn, mqChannel, mqQueue, &cfg.RabbitMQ.Publisher, logger)
	publisher.PublishAllRequests(requests)
}
