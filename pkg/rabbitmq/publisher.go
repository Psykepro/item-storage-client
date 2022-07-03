package rabbitmq

import (
	"fmt"

	domain "github.com/Psykepro/item-storage-client/_domain"
	"github.com/Psykepro/item-storage-client/config"
	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	cfg        *config.Publisher
	logger     domain.Logger
}

func NewPublisher(amqpConn *amqp.Connection, amqpChannel *amqp.Channel, amqpQueue *amqp.Queue, cfg *config.Publisher, logger domain.Logger) *Publisher {
	return &Publisher{
		connection: amqpConn,
		channel:    amqpChannel,
		queue:      amqpQueue,
		cfg:        cfg,
		logger:     logger,
	}
}

func (p *Publisher) Publish(body []byte) error {

	err := p.channel.Publish(
		p.cfg.Exchange,
		p.queue.Name,
		p.cfg.Mandatory,
		p.cfg.Immediate,
		amqp.Publishing{
			ContentType: p.cfg.ContentType,
			Body:        body,
		})

	return err
}

func (p *Publisher) PublishAllRequests(requests []*pb.ItemRequest) {
	for _, req := range requests {
		body, err := proto.Marshal(req)
		if err != nil {
			p.logger.Errorf("Failed to marshal item request. Err: [%s]", err)
			continue
		}
		p.logger.Infof("Publishing request - %s", getPrintMsgByCommandType(req))
		err = p.Publish(body)
		if err != nil {
			p.logger.Errorf("Failed to publish item request. Err: [%s]", err)
			continue
		}
	}
}

func getPrintMsgByCommandType(request *pb.ItemRequest) string {
	switch request.Command {
	case pb.Command_CREATE:
		return fmt.Sprintf("[%s]: uuid: %s, data: %s", request.Command, request.Item.Uuid, request.Item.Data)
	case pb.Command_GET, pb.Command_DELETE:
		return fmt.Sprintf("[%s]: uuid: %s", request.Command, request.Item.Uuid)
	case pb.Command_LIST:
		return fmt.Sprintf("[%s]", request.Command)
	}

	return ""
}
