CREATE TABLE IF NOT EXISTS payment (
    payment_id SERIAL PRIMARY KEY
    transact VARCHAR NOT NULL,
    requestId VARCHAR NOT NULL,
    currency VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    amount INTEGER NOT NULL,
    paymentDt INTEGER NOT NULL,
    bank VARCHAR NOT NULL,
    deliveryCost INTEGER NOT NULL,
    goodsTotal INTEGER NOT NULL,
    customFee INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery (
    delivery_id SERIAL PRIMARY KEY,
    delivery_name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
     zip VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    region VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
);

CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    order_uid VARCHAR UNIQUE NOT NULL,
    track_number VARCHAR NOT NULL,
    entry VARCHAR NOT NULL,
    delivery_id INTEGER REFERENCES delivery(delivery_id),
    payment_id INTEGER REFERENCES payment(payment_id),
    Locale VARCHAR NOT NULL,
    InternalSignature VARCHAR,
    CustomerId VARCHAR NOT NULL,
    DeliveryService VARCHAR NOT NULL,
    Shardkey VARCHAR NOT NULL,
    SmId INTEGER NOT NULL,
    DateCreated DATE NOT NULL,
    OofShard VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
    item_id SERIAL PRIMARY KEY,
    order_uid REFERENCES orders(order_uid),
    ChrtId INTEGER NOT NULL,
    TrackNumber VARCHAR NOT NULL,
    Price INTEGER NOT NULL,
    Rid VARCHAR NOT NULL,
    item_name VARCHAR NOT NULL,
    Sale INTEGER NOT NULL,
    item_Size VARCHAR NOT NULL,
    TotalPrice INTEGER NOT NULL,
    NmId INTEGER NOT NULL,
    Brand VARCHAR NOT NULL,
    Status INTEGER NOT NULL
);