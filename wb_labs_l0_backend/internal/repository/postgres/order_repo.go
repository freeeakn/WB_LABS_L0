package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"wb_labs_l0_backend/internal/domain/order"
	"wb_labs_l0_backend/internal/logger"
)

type PostgresOrderRepo struct {
	db  *sql.DB
	log logger.Logger
}

func New(db *sql.DB, log logger.Logger) *PostgresOrderRepo {
	return &PostgresOrderRepo{db: db, log: log}
}

func (r *PostgresOrderRepo) Save(ctx context.Context, o *order.Order) error {
	payload, err := json.Marshal(o)
	if err != nil {
		return err
	}

	var dt sql.NullTime
	if o.DateCreated != "" {
		if t, err := time.Parse(time.RFC3339, o.DateCreated); err == nil {
			dt = sql.NullTime{Time: t, Valid: true}
		}
	}

	q := `
INSERT INTO orders (order_uid, track_number, customer_id, date_created, payload)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (order_uid) DO UPDATE
  SET track_number = EXCLUDED.track_number,
      customer_id = EXCLUDED.customer_id,
      date_created = EXCLUDED.date_created,
      payload = EXCLUDED.payload
`
	_, err = r.db.ExecContext(ctx, q, o.OrderUID, o.TrackNumber, o.CustomerID, dt, payload)
	if err != nil {
		r.log.Errorf("repo save exec failed: %v", err)
		return err
	}
	return nil
}

func (r *PostgresOrderRepo) FindAll(ctx context.Context) ([]*order.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT payload FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*order.Order
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			r.log.Warnf("scan payload failed: %v", err)
			continue
		}
		var o order.Order
		if err := json.Unmarshal(raw, &o); err != nil {
			r.log.Warnf("unmarshal payload failed: %v", err)
			continue
		}
		res = append(res, &o)
	}
	return res, rows.Err()
}

func (r *PostgresOrderRepo) FindByID(ctx context.Context, id string) (*order.Order, error) {
	row := r.db.QueryRowContext(ctx, `SELECT payload FROM orders WHERE order_uid = $1`, id)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		return nil, err
	}
	var o order.Order
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}
	return &o, nil
}
