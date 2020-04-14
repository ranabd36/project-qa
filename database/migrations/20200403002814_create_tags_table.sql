-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS tags
(
    questions_id int not null,
    name         varchar(40)
);

CREATE INDEX idx_tags_name ON tags (name);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS tags;