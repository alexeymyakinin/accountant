-- +goose Up
-- +goose StatementBegin
create table users
(
    id       serial primary key,
    email    varchar not null unique,
    password varchar not null
);
create unique index users_email_uq_ix on users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index users_email_uq_ix;
drop table users;
-- +goose StatementEnd
