-- name: CreateEvent :execresult
INSERT INTO event (
   dt, money, description
) VALUES (
  ?, ?, ?
);

-- name: GetEvent :one
SELECT * FROM event
WHERE dt = ? AND money = ? AND description = ?;

-- name: GetEventByID :one
SELECT * FROM event
WHERE id = ?;

-- name: DeleteEventByID :exec
DELETE FROM event
WHERE id = ?;

-- name: ListOutcomeEvents :many
SELECT event.*, tag.name AS tags FROM event
LEFT OUTER JOIN event_to_tag ON event.id = event_to_tag.event_id
LEFT OUTER JOIN tag ON tag.id = event_to_tag.tag_id
WHERE
  (event.dt BETWEEN ? AND ?) AND
  (event.money < 0)
ORDER BY event.dt;

-- name: ListOutcomeEventsWithTags :many
SELECT event.*, tag.name AS tags FROM event
LEFT OUTER JOIN event_to_tag ON event.id = event_to_tag.event_id
LEFT OUTER JOIN tag ON tag.id = event_to_tag.tag_id
WHERE
  (event.dt BETWEEN ? AND ?) AND
  (tag.name IN (sqlc.slice(tags))) AND
  (event.money < 0)
ORDER BY event.dt;

-- name: ListEvents :many
SELECT event.*, tag.name AS tags FROM event
LEFT OUTER JOIN event_to_tag ON event.id = event_to_tag.event_id
LEFT OUTER JOIN tag ON tag.id = event_to_tag.tag_id
WHERE event.dt BETWEEN ? AND ?
ORDER BY event.dt;

-- name: ListEventsWithTags :many
SELECT event.*, tag.name AS tags FROM event
LEFT OUTER JOIN event_to_tag ON event.id = event_to_tag.event_id
LEFT OUTER JOIN tag ON tag.id = event_to_tag.tag_id
WHERE
  (event.dt BETWEEN ? AND ?) AND
  (tag.name IN (sqlc.slice(tags)))
ORDER BY event.dt;

-- name: CreateTag :execresult
INSERT INTO tag (name) VALUES (?);

-- name: GetTag :one
SELECT * FROM tag WHERE name = ?;

-- name: DeleteTagByID :exec
DELETE FROM tag
WHERE id = ?;
