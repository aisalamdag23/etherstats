-- migrate:up

CREATE TABLE IF NOT EXISTS balances (
    id SERIAL PRIMARY KEY, 
    address VARCHAR(255) NOT NULL,
    balance VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- migrate:down

DROP TABLE IF EXISTS balances;
