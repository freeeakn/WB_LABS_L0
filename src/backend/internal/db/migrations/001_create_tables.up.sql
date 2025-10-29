CREATE TABLE delivery (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    phone         TEXT NOT NULL,
    zip           TEXT NOT NULL,
    city          TEXT NOT NULL,
    address       TEXT NOT NULL,
    region        TEXT NOT NULL,
    email         TEXT NOT NULL
);

CREATE TABLE payment (
    id            SERIAL PRIMARY KEY,
    transaction   TEXT NOT NULL,
    request_id    TEXT,
    currency      TEXT NOT NULL,
    provider      TEXT NOT NULL,
    amount        INTEGER NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          TEXT NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total   INTEGER NOT NULL,
    custom_fee    INTEGER NOT NULL
);

CREATE TABLE item (
    id            SERIAL PRIMARY KEY,
    chrt_id       BIGINT NOT NULL,
    track_number  TEXT NOT NULL,
    price         INTEGER NOT NULL,
    rid           TEXT NOT NULL,
    name          TEXT NOT NULL,
    sale          INTEGER NOT NULL,
    size          TEXT NOT NULL,
    total_price   INTEGER NOT NULL,
    nm_id         BIGINT NOT NULL,
    brand         TEXT NOT NULL,
    status        INTEGER NOT NULL
);

CREATE TABLE orders (
    order_uid          TEXT PRIMARY KEY,
    track_number       TEXT NOT NULL,
    entry              TEXT NOT NULL,
    delivery_id        INTEGER REFERENCES delivery(id),
    payment_id         INTEGER REFERENCES payment(id),
    locale             TEXT NOT NULL,
    internal_signature TEXT,
    customer_id        TEXT NOT NULL,
    delivery_service   TEXT NOT NULL,
    shardkey           TEXT NOT NULL,
    sm_id              INTEGER NOT NULL,
    date_created       TIMESTAMPTZ NOT NULL,
    oof_shard          TEXT NOT NULL
);

CREATE TABLE order_items (
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    item_id   INTEGER REFERENCES item(id),
    PRIMARY KEY (order_uid, item_id)
);