-- +goose Up
-- +goose StatementBegin
ALTER TABLE Wg_peers
ALTER COLUMN server_ip TYPE VARCHAR(32) USING server_ip::VARCHAR(32),
ALTER COLUMN provided_ip TYPE VARCHAR(32) USING provided_ip::VARCHAR(32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Wg_peers
ALTER COLUMN server_ip TYPE inet USING server_ip::inet,
ALTER COLUMN provided_ip TYPE inet USING provided_ip::inet;
-- +goose StatementEnd
