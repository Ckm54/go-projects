-- name: CreateFeedFollow :one
WITH inserted AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  ON CONFLICT (user_id, feed_id) DO NOTHING
  RETURNING id, created_at, updated_at, user_id, feed_id
),
selected AS (
  SELECT * FROM inserted
  UNION ALL
  SELECT id, created_at, updated_at, user_id, feed_id
  FROM feed_follows
  WHERE user_id = $4 AND feed_id = $5
  LIMIT 1
)
SELECT 
  selected.id, selected.created_at, selected.updated_at, selected.user_id, selected.feed_id,
  users.name AS user_name,
  feeds.name AS feed_name
FROM selected
JOIN users ON selected.user_id = users.id
JOIN feeds ON selected.feed_id = feeds.id;


-- name: GetFeedFollowsForUser :many
SELECT 
  ff.id,
  ff.created_at,
  ff.updated_at,
  u.name AS user_name,
  f.name AS feed_name
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2;