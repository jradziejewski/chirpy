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

-- name: GetUserFromRefreshToken :one
select * from users u
inner join refresh_tokens r
on r.user_id = u.id
where r.token = $1
and expires_at > current_timestamp
and revoked_at is null;

-- name: DeleteUsers :exec
DELETE FROM users;
