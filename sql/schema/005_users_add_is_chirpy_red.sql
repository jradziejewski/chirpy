-- +goose Up
alter table users
add is_chirpy_red boolean default false;

-- +goose Down
alter table users
drop is_chirpy_red;
