CREATE TABLE orders
(
    id serial PRIMARY KEY NOT NULL,
    user_id integer NOT NULL,
    restaurant_id integer NOT NULL,
    total integer NOT NULL,
    currency_code character varying NOT NULL,
    placed_at timestamp with time zone NOT NULL
)
