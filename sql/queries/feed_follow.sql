-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        gen_random_uuid(),
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        $1,
        $2
    )
    RETURNING *
)
SELECT
    feed_follows.*,
    users.name as user_name,
    feeds.name as feed_name
FROM inserted_feed_follow AS feed_follows
JOIN users ON feed_follows.user_id = users.id
JOIN feeds ON feed_follows.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.id,
    ff.created_at,
    ff.updated_at,
    ff.user_id,
    ff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1
AND ff.feed_id = $2;
