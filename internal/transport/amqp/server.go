package amqpSrv

import (
	"context"
	amqpHandler "github.com/ernur-eskermes/crud-audit-log/internal/transport/amqp/handlers"
	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type Server struct {
	conn     *amqp.Connection
	logger   *logging.Logger
	handlers *amqpHandler.Handler
}

func New(addr string, logger *logging.Logger, handlers *amqpHandler.Handler) (*Server, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:     conn,
		logger:   logger,
		handlers: handlers,
	}, nil
}

func (s *Server) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		var f func(ctx context.Context, message []byte) error

		switch delivery.Type {
		case "log":
			//f = amqpHandler.Log
		}

		if err := f(ctx, delivery.Body); err != nil {
			if err = delivery.Reject(false); err != nil {
				s.logger.Errorf("Err delivery.Reject: %v", err)
			}
			s.logger.Errorf("Failed to process delivery: %v", err)
		} else {
			err = delivery.Ack(false)
			if err != nil {
				s.logger.Errorf("Failed to acknowledge delivery: %v", err)
			}
		}
	}
}
func (s *Server) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	s.logger.Infof("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	s.logger.Infof("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	s.logger.Infof("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}
func (s *Server) ListenAndServe(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := s.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.Wrap(err, "CreateChannel")
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	for i := 0; i < workerPoolSize; i++ {
		go s.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	s.logger.Errorf("ch.NotifyClose: %v", chanErr)
	return chanErr
}

func (s *Server) Stop() {

}
