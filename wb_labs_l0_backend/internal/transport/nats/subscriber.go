package natstransport

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"wb_labs_l0_backend/internal/config"
	"wb_labs_l0_backend/internal/domain/order"
	"wb_labs_l0_backend/internal/logger"
	"wb_labs_l0_backend/internal/service"

	stan "github.com/nats-io/stan.go"
)

type Subscriber struct {
	sc      stan.Conn
	sub     stan.Subscription
	cfg     *config.Config
	svc     *service.OrderService
	log     logger.Logger
	subject string
}

func NewSubscriber(cfg *config.Config, svc *service.OrderService, log logger.Logger) (stan.Conn, *Subscriber, error) {
	sc, err := stan.Connect(cfg.NatsClusterID, cfg.NatsClientID, stan.NatsURL(cfg.NatsURL))
	if err != nil {
		return nil, nil, err
	}
	s := &Subscriber{
		sc:      sc,
		cfg:     cfg,
		svc:     svc,
		log:     log,
		subject: cfg.NatsSubject,
	}
	return sc, s, nil
}

func (s *Subscriber) Start() error {
	if s.sc == nil {
		return errors.New("stan conn nil")
	}
	var err error
	s.sub, err = s.sc.Subscribe(s.subject, s.handleMsg,
		stan.SetManualAckMode(),
		stan.DurableName("wb_labs_l0_backend_durable"),
		stan.AckWait(30*time.Second),
		stan.MaxInflight(20),
	)
	if err != nil {
		return err
	}
	s.log.Infof("nats streaming subscribed to %s", s.subject)
	return nil
}

func (s *Subscriber) Stop() {
	if s.sub != nil {
		_ = s.sub.Unsubscribe()
	}
}

func (s *Subscriber) handleMsg(m *stan.Msg) {
	if len(m.Data) == 0 {
		s.log.Warn("empty message, ack and skip")
		_ = m.Ack()
		return
	}
	if len(m.Data) > 1024*1024 {
		s.log.Warn("message too large, ack and skip")
		_ = m.Ack()
		return
	}

	var o order.Order
	if err := json.Unmarshal(m.Data, &o); err != nil {
		s.log.Warnf("invalid json payload: %v", err)
		_ = m.Ack()
		return
	}

	if err := o.Validate(); err != nil {
		s.log.Warnf("validation failed: %v", err)
		_ = m.Ack()
		return
	}

	ctx := context.Background()
	if err := s.svc.SaveOrder(ctx, &o); err != nil {
		s.log.Errorf("save order failed: %v (will not ack so it can be redelivered)", err)
		return
	}
	if err := m.Ack(); err != nil {
		s.log.Warnf("ack failed: %v", err)
		return
	}
	s.log.Infof("processed order %s (seq=%d)", o.OrderUID, m.Sequence)
}
