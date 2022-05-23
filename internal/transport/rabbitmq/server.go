package rabbitmq

import (
	"context"
	"fmt"

	"github.com/ernur-eskermes/crud-audit-log/internal/config"
	amqpHandler "github.com/ernur-eskermes/crud-audit-log/internal/transport/rabbitmq/handlers"
	"github.com/streadway/amqp"
)

type Server struct {
	conn     *amqp.Connection
	handlers *amqpHandler.Handler
}

func New(cfg *config.Config, amqpHandlers *amqpHandler.Handler) (*Server, error) {
	conn, err := amqp.Dial(cfg.AMQP.URI)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Failed to connect to RabbitMQ", err)
	}

	return &Server{conn: conn, handlers: amqpHandlers}, nil
}

func (s *Server) ListenAndServe() error {
	ch, err := s.conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to open a channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"logs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to declare a queue", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", "Failed to register a consumer", err)
	}

	for d := range msgs {
		s.handlers.Audit.Log(context.TODO(), d)
	}

	return nil
}

func (s *Server) Stop() error {
	return s.conn.Close()
}
