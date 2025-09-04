-- +goose Up
ALTER TABLE appointments ADD COLUMN category VARCHAR(255);

-- +goose Down
ALTER TABLE appointments DROP COLUMN category;
