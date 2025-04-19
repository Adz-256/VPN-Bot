-- +goose Up
-- +goose StatementBegin
INSERT INTO Plans (country, duration_days, price, description)
VALUES
('NDR', 30, 149, '1 месяц подписки'),
('NDR', 90,  399, '3 месяца подписки'),
('NDR', 180, 799, '6 месяцев подписки'),
('NDR', 365, 1500, '12 месяцев подписки');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM Plans
-- +goose StatementEnd
