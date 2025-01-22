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

-- name: DeleteChirps :exec
DELETE FROM chirps;
