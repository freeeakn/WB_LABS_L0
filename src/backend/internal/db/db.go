package db

import (
	"backend/internal/model"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB — обёртка над пулом соединений
type DB struct {
	pool *pgxpool.Pool
}

// New — создаёт пул с retry-логикой
func New(dsn string) (*DB, error) {
	const maxRetries = 10
	const delay = 2 * time.Second

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		cfg, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			return nil, parseErr
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err == nil {
			if pingErr := pool.Ping(context.Background()); pingErr == nil {
				log.Printf("Connected to PostgreSQL (attempt %d)", i+1)
				return &DB{pool: pool}, nil
			}
		}

		log.Printf("Failed to connect to DB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(delay)
	}

	return nil, err
}

// Save сохраняет Order + все связанные сущности (транзакция)
func (d *DB) Save(ctx context.Context, o *model.Order) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// delivery
	var deliveryID int
	err = tx.QueryRow(ctx,
		`INSERT INTO delivery(name,phone,zip,city,address,region,email)
         VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email,
	).Scan(&deliveryID)
	if err != nil {
		return err
	}

	// payment
	var paymentID int
	err = tx.QueryRow(ctx,
		`INSERT INTO payment(transaction,request_id,currency,provider,amount,payment_dt,bank,
         delivery_cost,goods_total,custom_fee)
         VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank,
		o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee,
	).Scan(&paymentID)
	if err != nil {
		return err
	}

	// order
	_, err = tx.Exec(ctx,
		`INSERT INTO orders(order_uid,track_number,entry,delivery_id,payment_id,locale,
         internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard)
         VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		o.OrderUID, o.TrackNumber, o.Entry, deliveryID, paymentID, o.Locale,
		o.InternalSignature, o.CustomerID, o.DeliveryService,
		o.Shardkey, o.SmID, o.DateCreated, o.OofShard,
	)
	if err != nil {
		return err
	}

	// items + order_items
	for _, it := range o.Items {
		var itemID int
		err = tx.QueryRow(ctx,
			`INSERT INTO item(chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status)
             VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`,
			it.ChrtID, it.TrackNumber, it.Price, it.Rid, it.Name,
			it.Sale, it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status,
		).Scan(&itemID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items(order_uid,item_id) VALUES($1,$2)`,
			o.OrderUID, itemID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// LoadAll заполняет кэш из БД при старте
func (d *DB) LoadAll(ctx context.Context) (map[string]*model.Order, error) {
	rows, err := d.pool.Query(ctx,
		`SELECT o.order_uid, o.track_number, o.entry,
                d.name,d.phone,d.zip,d.city,d.address,d.region,d.email,
                p.transaction,p.request_id,p.currency,p.provider,p.amount,p.payment_dt,
                p.bank,p.delivery_cost,p.goods_total,p.custom_fee,
                o.locale,o.internal_signature,o.customer_id,o.delivery_service,
                o.shardkey,o.sm_id,o.date_created,o.oof_shard
         FROM orders o
         JOIN delivery d ON o.delivery_id=d.id
         JOIN payment p ON o.payment_id=p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cache := make(map[string]*model.Order)

	for rows.Next() {
		var o model.Order
		var del model.Delivery
		var pay model.Payment
		err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry,
			&del.Name, &del.Phone, &del.Zip, &del.City, &del.Address, &del.Region, &del.Email,
			&pay.Transaction, &pay.RequestID, &pay.Currency, &pay.Provider, &pay.Amount, &pay.PaymentDt,
			&pay.Bank, &pay.DeliveryCost, &pay.GoodsTotal, &pay.CustomFee,
			&o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService,
			&o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
		)
		if err != nil {
			return nil, err
		}
		o.Delivery = del
		o.Payment = pay

		// items отдельным запросом
		itemRows, _ := d.pool.Query(ctx,
			`SELECT i.chrt_id,i.track_number,i.price,i.rid,i.name,i.sale,i.size,
                    i.total_price,i.nm_id,i.brand,i.status
             FROM order_items oi
             JOIN item i ON oi.item_id=i.id
             WHERE oi.order_uid=$1`, o.OrderUID)
		for itemRows.Next() {
			var it model.Item
			_ = itemRows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
				&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status)
			o.Items = append(o.Items, it)
		}
		itemRows.Close()

		cache[o.OrderUID] = &o
	}
	return cache, rows.Err()
}
