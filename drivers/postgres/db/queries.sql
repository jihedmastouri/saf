-- name: addQueue
INSERT INTO queues (name)
VALUES ($1)
RETURNING *;

-- name: enqueue
INSERT INTO jobs (name, queue_name, payload)
VALUES ((SELECT id FROM Queues WHERE name = $1), $2, $3)
RETURNING *;

-- name: getJobById :one
SELECT *
FROM jobs
WHERE id = $1;

-- name: countJobsByQId :one
SELECT count(*)
FROM Jobs
WHERE queue_name = $1;

-- name: dequeue
UPDATE jobs
SET status = 'active'
WHERE id = (
	SELECT id
	FROM jobs
	WHERE queue_name = $1
	AND status = 'pending'
	ORDER BY priority ASC, created_at ASC
	LIMIT 1
)
RETURNING *;

-- name: updateJobStatus
UPDATE jobs
SET status = $2
WHERE id = $1;

-- name: deleteJobById
DELETE
FROM jobs
RETURNING id;
