-- name: CreateFeedFollow :one
WITH ff AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, f.name, u.name FROM
ff INNER JOIN feeds f ON ff.feed_id = f.id INNER JOIN users u ON f.user_id = u.id;

-- name: GetFeedFollows :many
WITH ff AS (
    SELECT * FROM feed_follows WHERE feed_follows.user_id = $1
    )
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, f.name, u.name FROM
ff INNER JOIN feeds f ON ff.feed_id = f.id INNER JOIN users u ON f.user_id = u.id;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2 RETURNING user_id, feed_id;
