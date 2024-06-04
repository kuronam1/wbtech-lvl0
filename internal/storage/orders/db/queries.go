package db

const (
	CreateOrderQuery = `
INSERT INTO orders(order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`
	SelectOrderQuery = `
SELECT order_id, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
FROM orders
WHERE order_uid = $1
`
	SelectOrdersQuery = `
SELECT order_id, order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
FROM orders
`

	CreatePaymentQuery = `
INSERT INTO payment(transact, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING payment_id
`
	SelectPaymentQuery = `
SELECT transact, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
FROM payment
WHERE payment_id = $1
`
	SelectPaymentsQuery = `
SELECT payment_id, transact, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
FROM payment
`
	CreateDeliveryQuery = `
INSERT INTO delivery(delivery_name, phone, zip, city, address, region, email)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING delivery_id
`
	SelectDeliveryQuery = `
SELECT delivery_name, phone, zip, city, address, region, email
FROM delivery
WHERE delivery_id = $1
`
	SelectDeliveriesQuery = `
SELECT delivery_id, delivery_name, phone, zip, city, address, region, email
FROM delivery
`
	CreateItemQuery = `
INSERT INTO item(order_uid, chrt_id, track_number, price, rid, item_name, sale, item_size, total_price, nm_id, brand, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`
	SelectItemsQuery = `
SELECT item_id, order_uid, chrt_id, track_number, price, rid, item_name, sale, item_size, total_price, nm_id, brand, status
FROM item
WHERE order_uid = $1
`

	CheckCacheQuery = `SELECT EXISTS(SELECT 1 FROM orders)`

	SelectOrderQueryV2 = `
SELECT order_id, order_uid, track_number, entry, locale, internal_signature,
       customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard,
       o.delivery_id, delivery_name, phone, zip, city, address, region, email,
       p.payment_id, transact, request_id, currency, provider, amount, payment_dt,
       bank, delivery_cost, goods_total, custom_fee
FROM orders o
         JOIN delivery d on d.delivery_id = o.delivery_id
         JOIN payment p on p.payment_id = o.payment_id
WHERE order_id = $1
`
)
