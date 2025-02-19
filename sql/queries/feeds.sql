-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: ListFeedsWithUsers :many
SELECT feeds.*, users.name as user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
WHERE last_fetched_at IS NULL
   OR last_fetched_at = (
      SELECT MIN(last_fetched_at)
      FROM feeds
   )
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
