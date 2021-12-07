-- migrate:up
CREATE TABLE products (
    id INTEGER PRIMARY KEY,
    name VARCHAR (256)
)

-- migrate:down
DROP TABLE products
