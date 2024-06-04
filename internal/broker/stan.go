package broker

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"io"
	"log/slog"
	"os"
	"time"
	"wbLvL0/internal/appErrors"
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

	return broker
}

func (s *Stan) Publish(ctx context.Context, subject string) {
	go func() {
		data, err := s.GetData()
		if err != nil {
			s.Logger.Error("[Pub] error while getting data: ", appErrors.WrapLogErr(err))
			s.Notify <- err
			close(s.Notify)
			return
		}

		ticker := time.NewTicker(5 * time.Second)
		counter := 0
		for {
			select {
			case <-ticker.C:
				msg := data[counter]
				pubMsg, err := json.Marshal(msg)
				if err != nil {
					s.Logger.Error("[Pub] error while marshaling data: ", appErrors.WrapLogErr(err))
					s.Notify <- err
					close(s.Notify)
					return
				}

				if err := s.Conn.Publish(subject, pubMsg); err != nil {
					s.Logger.Error("[Pub] error while publishing msg: ", appErrors.WrapLogErr(err))
					s.Notify <- err
					close(s.Notify)
					return
				}
			case <-ctx.Done():
				s.Logger.Info("[Pub] exiting")
				return
			}
		}
	}()
}

func (s *Stan) Subscribe(ctx context.Context, subject string, save func(data models.Order) error) {
	go func() {
		sub, err := s.Conn.Subscribe(subject, func(msg *stan.Msg) {
			order := models.Order{}
			if err := json.Unmarshal(msg.Data, &order); err != nil {
				s.Logger.Error("[Stan] error while parsing order", appErrors.WrapLogErr(err))
				return
			}
			validate := validator.New(validator.WithRequiredStructEnabled())
			err := validate.Struct(order)
			if err != nil {
				/*if _, ok := err.(*validator.InvalidValidationError); !ok {
					s.logger.Error("error while validation: [%s]", appErrors.WrapLogErr(err))
					return
				}*/
				s.Logger.Error("[Stan] error while validation: [%s]", appErrors.WrapLogErr(err))
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

				return
			}

			if err := save(order); err != nil {
				s.Logger.Error("[Stan] error while saving massage", appErrors.WrapLogErr(err))
				return
			}
		}, stan.DeliverAllAvailable())
		if err != nil {
			s.Logger.Error("[Stan] error while subscribe", appErrors.WrapLogErr(err))
			s.Notify <- err
			close(s.Notify)
		}

		<-ctx.Done()

		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				s.Logger.Error("[Stan] error while unsubscribe", appErrors.WrapLogErr(err))
			}
		}()
	}()
}

func (s *Stan) GetData() ([]models.Order, error) {
	file, err := os.Open("dataForPub.json")
	if err != nil {
		s.Logger.Error("[Pub.GetData] error while reading file: ", appErrors.WrapLogErr(err))
		s.Notify <- err
		close(s.Notify)
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		s.Logger.Error("[Pub.GetData] error while reading file: ", appErrors.WrapLogErr(err))
		s.Notify <- err
		close(s.Notify)
		return nil, err
	}

	var slice []models.Order
	err = json.Unmarshal(data, &slice)

	return slice, err
}
