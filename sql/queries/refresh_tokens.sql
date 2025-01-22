-- name: CreateRefreshToken :one
insert into refresh_tokens (
	token,
	created_at,
	updated_at,
	expires_at,
	revoked_at,
	user_id
	)
values (
	$1,
	$2,
	$3,
	$4,
	null,
	$5
)
returning *;

-- name: GetRefreshToken :one
select * from refresh_tokens
where token = $1;

-- name: RevokeToken :exec
update refresh_tokens
set revoked_at = $1, updated_at = $1
where token = $2;
