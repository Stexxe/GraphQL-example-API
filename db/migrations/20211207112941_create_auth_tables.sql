-- migrate:up
CREATE TABLE codes (
    user_id INTEGER,
    code INTEGER,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX idx_code ON codes(code);

CREATE TABLE tokens (
    user_id INTEGER,
    token VARCHAR (32),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE INDEX idx_token ON tokens(token);

-- migrate:down
DROP INDEX idx_code;
DROP TABLE codes;

DROP INDEX idx_token;
DROP TABLE tokens;

