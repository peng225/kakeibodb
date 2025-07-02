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

-- name: ListTags :many
SELECT * FROM tag
ORDER BY tag.id;

-- name: CreatePattern :execresult
INSERT INTO pattern (key_string) VALUES (?);

-- name: GetPattern :one
SELECT * FROM pattern WHERE key_string = ?;

-- name: DeletePatternByID :exec
DELETE FROM pattern
WHERE id = ?;

-- name: ListPatterns :many
SELECT pattern.*, tag.name AS tags FROM pattern
LEFT OUTER JOIN pattern_to_tag ON pattern.id = pattern_to_tag.pattern_id
LEFT OUTER JOIN tag ON tag.id = pattern_to_tag.tag_id
ORDER BY pattern.id;

-- name: MapEventToTag :execresult
INSERT INTO event_to_tag (event_id, tag_id) VALUES (?, ?);

-- name: UnmapEventFromTag :execresult
DELETE FROM event_to_tag
where event_id = ? AND tag_id = ?;

-- name: GetEventToTagMap :one
SELECT * FROM event_to_tag
WHERE event_id = ? AND tag_id = ?;

-- name: MapPatternToTag :execresult
INSERT INTO pattern_to_tag (pattern_id, tag_id) VALUES (?, ?);

-- name: UnmapPatternFromTag :execresult
DELETE FROM pattern_to_tag
where pattern_id = ? AND tag_id = ?;

-- name: GetPatternToTagMap :one
SELECT * FROM pattern_to_tag
WHERE pattern_id = ? AND tag_id = ?;
