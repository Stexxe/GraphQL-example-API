-- migrate:up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR (256)
)

-- migrate:down
DROP TABLE products
