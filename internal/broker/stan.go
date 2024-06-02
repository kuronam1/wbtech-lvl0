package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log/slog"
	"time"
	"wbLvL0/internal/models"
)

type Stan struct {
	conn   stan.Conn
	logger *slog.Logger
}

func NewBroker(connection stan.Conn, log *slog.Logger) *Stan {
	return &Stan{
		conn:   connection,
		logger: log,
	}
}

func (s *Stan) Publish(subject string, msg []byte) error {
	return s.conn.Publish(subject, msg)
}

func (s *Stan) Subscribe(subject string, save func(ctx context.Context, data models.Order)) error {
	sub, err := s.conn.Subscribe(subject, func(msg *stan.Msg) {
		order := models.Order{}
		if err := json.Unmarshal(msg.Data, &order); err != nil {

		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(order)
		if err != nil {
			/*if _, ok := err.(*validator.InvalidValidationError); !ok {
				s.logger.Error("error while validation: [%s]", errors.WrapLogErr(err))
				return
			}*/

			for _, err := range err.(validator.ValidationErrors) {
				s.logger.Debug("%s", err.Namespace())
				s.logger.Debug("%s", err.Field())
				s.logger.Debug("%s", err.StructNamespace())
				s.logger.Debug("%s", err.StructField())
				s.logger.Debug("%s", err.Tag())
				s.logger.Debug("%s", err.ActualTag())
				s.logger.Debug("%s", err.Kind())
				s.logger.Debug("%s", err.Type())
				s.logger.Debug("%s", err.Value())
				s.logger.Debug("%s", err.Param())
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		save(ctx, order)
		cancel()

	}, stan.DeliverAllAvailable())
	if err != nil {
		return fmt.Errorf("[broker] failed to subscribe to the channel: %w", err)
	}
	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			return
		}
	}()

	return err
}
