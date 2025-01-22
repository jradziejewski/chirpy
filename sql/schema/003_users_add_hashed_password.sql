-- +goose Up
alter table users
add hashed_password text default 'unset';

-- +goose Down
alter table users
drop hashed_password;
