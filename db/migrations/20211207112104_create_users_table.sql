-- migrate:up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR (16)
);

CREATE INDEX idx_phone ON users(phone);

-- migrate:down
DROP INDEX idx_phone;
DROP TABLE users;
