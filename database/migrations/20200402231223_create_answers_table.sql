-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS answers
(
    id          serial not null,
    user_id     int    not null,
    question_id int    not null,
    answer_id   int    null,
    description text   not null,
    is_accepted boolean     default false,
    created_at  timestamp default current_timestamp,
    updated_at  timestamp default current_timestamp,

    primary key (id),
    foreign key (user_id) references users (id),
    foreign key (question_id) references questions (id),
    foreign key (answer_id) references answers (id)
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS answers;
