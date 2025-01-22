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
and expires_at > NOW()
and revoked_at is null;

-- name: UpdateEmailAndPassword :one
update users
set email = $1, hashed_password = $2, updated_at = NOW()
where id = $3
returning *;

-- name: UpdateIsChirpyRed :one
update users
set is_chirpy_red = $1, updated_at = NOW()
where id = $2
returning *;

-- name: DeleteUsers :exec
DELETE FROM users;
