-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users (
    id serial PRIMARY KEY,
    chat_id BIGINT UNIQUE,
    username TEXT,
    is_admin BOOL DEFAULT FALSE,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Plans (
    id serial primary key,
    name text NOT NULL UNIQUE,
    duration_days integer NOT NULL CHECK ( duration_days > 0 ),
    price DECIMAL(10,2),
    description TEXT
    );

CREATE TABLE IF NOT EXISTS Payments (
    id serial PRIMARY KEY,
    transaction_id TEXT NOT NULL UNIQUE,
    user_id INT NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    plan_id INT NOT NULL REFERENCES Plans(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    status TEXT NOT NULL CHECK (status IN ('pending', 'canceled', 'paid')) DEFAULT 'pending',
    method TEXT NOT NULL,
    created_at timestamptz NOT NULL default CURRENT_TIMESTAMP,
    paid_at timestamptz
);


CREATE TABLE IF NOT EXISTS Wg_peers (
    id SERIAL PRIMARY KEY, -- уникальный ID внутри базы
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- владелец peer-а
    public_key TEXT NOT NULL UNIQUE, -- публичный ключ WireGuard
    config_file TEXT NOT NULL,
    server_ip inet NOT NULL,
    provided_ip inet NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users CASCADE;
DROP TABLE Plans CASCADE ;
DROP TABLE Payments CASCADE;
DROP TABLE Wg_peers CASCADE
-- +goose StatementEnd
