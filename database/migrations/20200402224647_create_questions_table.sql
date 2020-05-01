-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS questions
(
    id           serial              not null,
    user_id      int                 not null,
    title        varchar(255) unique not null,
    description  text                not null,
    published_at timestamp         null,
    created_at   timestamp default current_timestamp,
    updated_at   timestamp default current_timestamp,

    primary key (id),
    foreign key (user_id) references users (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS questions;