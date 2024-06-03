package broker

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log/slog"
	"wbLvL0/internal/errors"
	"wbLvL0/internal/storage/orders/models"
)

type Stan struct {
	Conn   stan.Conn
	Logger *slog.Logger
	Notify chan error
}

func New(connection stan.Conn, log *slog.Logger) *Stan {
	broker := &Stan{
		Conn:   connection,
		Logger: log,
		Notify: make(chan error),
	}

	broker.Subscribe()

	return broker
}

func (s *Stan) Publish(subject string, msg []byte) error {
	return s.Conn.Publish(subject, msg)
}

func (s *Stan) Subscribe(ctx context.Context, subject string, save func(data models.Order) error) {
	go func() {
		sub, err := s.Conn.Subscribe(subject, func(msg *stan.Msg) {
			order := models.Order{}
			if err := json.Unmarshal(msg.Data, &order); err != nil {
				s.Logger.Error("[Stan] error while parsing order", errors.WrapLogErr(err))
				s.Notify <- err
				close(s.Notify)
			}
			validate := validator.New(validator.WithRequiredStructEnabled())
			err := validate.Struct(order)
			if err != nil {
				/*if _, ok := err.(*validator.InvalidValidationError); !ok {
					s.logger.Error("error while validation: [%s]", errors.WrapLogErr(err))
					return
				}*/
				s.Logger.Error("[Stan] error while validation: [%s]", errors.WrapLogErr(err))
				for _, err := range err.(validator.ValidationErrors) {
					s.Logger.Debug("%s", err.Namespace())
					s.Logger.Debug("%s", err.Field())
					s.Logger.Debug("%s", err.StructNamespace())
					s.Logger.Debug("%s", err.StructField())
					s.Logger.Debug("%s", err.Tag())
					s.Logger.Debug("%s", err.ActualTag())
					s.Logger.Debug("%s", err.Kind())
					s.Logger.Debug("%s", err.Type())
					s.Logger.Debug("%s", err.Value())
					s.Logger.Debug("%s", err.Param())
				}

				s.Notify <- err
				close(s.Notify)
			}

			if err := save(order); err != nil {
				s.Logger.Error("[Stan] error while saving massage", errors.WrapLogErr(err))
				s.Notify <- err
				close(s.Notify)
			}

			<-ctx.Done()

		}, stan.DeliverAllAvailable())
		if err != nil {
			s.Logger.Error("[Stan] error while subscribe", errors.WrapLogErr(err))
			s.Notify <- err
			close(s.Notify)
		}
		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				s.Logger.Error("[Stan] error while unsubscribe", errors.WrapLogErr(err))
				return
			}
		}()
	}()
}
