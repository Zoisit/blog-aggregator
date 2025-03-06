-- name: CreateFeedFollow :one
WITH new_row AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
) SELECT
    new_row.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM new_row 
INNER JOIN users ON new_row.user_id = users.id 
INNER JOIN feeds ON new_row.feed_id = feeds.id; 

-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name, users.name AS user_name FROM feeds INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id INNER JOIN users ON users.id = feed_follows.user_id WHERE users.name = $1;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows A USING feeds B WHERE A.user_id = $1 AND B.url = $2;
