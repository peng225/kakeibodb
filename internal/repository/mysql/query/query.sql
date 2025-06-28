-- name: CreateEvent :execresult
INSERT INTO event (
   dt, money, description
) VALUES (
  ?, ?, ?
);

-- name: GetEvent :one
SELECT * FROM event
where dt = ? AND money = ? AND description = ?;

-- name: GetEventByID :one
SELECT * FROM event
where id = ?;

-- name: DeleteEventByID :exec
DELETE FROM event
where id = ?;
