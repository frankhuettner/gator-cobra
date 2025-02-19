-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
    )
    ON CONFLICT (url) DO UPDATE
    SET
        updated_at = EXCLUDED.updated_at,
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        published_at = EXCLUDED.published_at
    RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
JOIN feed_follows ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
