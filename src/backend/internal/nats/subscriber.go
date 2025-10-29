package nats

import (
	"backend/internal/cache"
	"backend/internal/db"
	"backend/internal/validator"
	"context"
	"log"

	stan "github.com/nats-io/stan.go"
)

const (
	channel  = "orders"
	clientID = "backend"
	durable  = "order-durable"
)

func Start(ctx context.Context, natsURL string, db *db.DB, c *cache.Cache) error {
	sc, err := stan.Connect("test-cluster", clientID, stan.NatsURL(natsURL))
	if err != nil {
		return err
	}

	_, err = sc.Subscribe(channel, func(m *stan.Msg) {
		order, err := validator.ValidateOrder(m.Data)
		if err != nil {
			log.Printf("invalid message: %v", err)
			m.Ack() // ack всё равно, чтобы не зациклить
			return
		}

		if err := db.Save(ctx, order); err != nil {
			log.Printf("db save error: %v", err)
			return // не ack – NATS будет повторять
		}

		c.Set(order)
		m.Ack()
	}, stan.DurableName(durable), stan.SetManualAckMode())

	if err != nil {
		return err
	}

	log.Println("NATS subscriber started")
	<-ctx.Done()
	return sc.Close()
}
