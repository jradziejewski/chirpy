-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values (
	$1,
	$2,
	$3,
	$4,
	$5
)
returning *;

-- name: GetChirps :many
select * from chirps
where coalesce(sqlc.narg('user_id'), user_id) is null or user_id =coalesce(sqlc.narg('user_id'), user_id) 
order by created_at;

-- name: GetChirp :one
select * from chirps
where id = $1;

-- name: DeleteChirps :exec
delete from chirps;

-- name: DeleteChirp :exec
delete from chirps
where id = $1;
