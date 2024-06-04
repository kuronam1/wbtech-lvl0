CREATE TABLE IF NOT EXISTS payment (
    payment_id SERIAL PRIMARY KEY,
    transact VARCHAR NOT NULL,
    request_id VARCHAR NOT NULL,
    currency VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt INTEGER NOT NULL,
    bank VARCHAR NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery (
    delivery_id SERIAL PRIMARY KEY,
    delivery_name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
     zip VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    region VARCHAR NOT NULL,
    email VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    order_uid VARCHAR UNIQUE NOT NULL,
    track_number VARCHAR NOT NULL,
    entry VARCHAR NOT NULL,
    delivery_id INTEGER REFERENCES delivery(delivery_id),
    payment_id INTEGER REFERENCES payment(payment_id),
    locale VARCHAR NOT NULL,
    internal_signature VARCHAR,
    customer_id VARCHAR NOT NULL,
    delivery_service VARCHAR NOT NULL,
    shard_key VARCHAR NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created DATE NOT NULL,
    oof_shard VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
    item_id SERIAL PRIMARY KEY,
    order_uid VARCHAR REFERENCES orders(order_uid),
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR NOT NULL,
    item_name VARCHAR NOT NULL,
    sale INTEGER NOT NULL,
    item_Size VARCHAR NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR NOT NULL,
    status INTEGER NOT NULL
);