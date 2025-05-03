-- +goose Up
-- +goose StatementBegin
ALTER TABLE Wg_peers ADD COLUMN blocked BOOL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Wg_peers DROP COLUMN blocked;
-- +goose StatementEnd
