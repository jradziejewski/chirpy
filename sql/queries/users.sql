-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
	$1,
	$2,
	$3,
	$4,
	$5
)
returning *;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;
