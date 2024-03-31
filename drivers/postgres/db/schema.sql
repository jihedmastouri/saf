CREATE TYPE job_status AS ENUM ('pending', 'failed', 'completed', 'active');

CREATE TABLE IF NOT EXISTS Queues (
	id serial PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

	-- Monitoring fields
	active_jobs INT NOT NULL DEFAULT 0,
	completed_jobs INT NOT NULL DEFAULT 0,
	failed_jobs INT NOT NULL DEFAULT 0,
	all_jobs INT NOT NULL DEFAULT 0,

	-- Organizational fields
	max_concurrency INT NOT NULL DEFAULT 1,
	max_jobs INT NOT NULL DEFAULT 0,
	max_completed_jobs INT NOT NULL DEFAULT 0,
	max_failed_jobs INT NOT NULL DEFAULT 0,

	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE IF NOT EXISTS Jobs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	queue_id serial NOT NULL,
	parent_id UUID,
	name TEXT NOT NULL,
	status job_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	priority INT NOT NULL DEFAULT 0,
	payload JSON,
	ChildrenData JSON,

	-- Organizational fields
	remove_on_complete BOOLEAN NOT NULL DEFAULT FALSE,
	remove_on_fail BOOLEAN NOT NULL DEFAULT FALSE,
	retry BOOLEAN NOT NULL DEFAULT FALSE,
	max_attempts INT NOT NULL DEFAULT 3,
	attempts INT NOT NULL DEFAULT 0,

	FOREIGN KEY (queue_id) REFERENCES Queues(id)
)

CREATE INDEX time_inserted_idx ON jobs (created_at ASC, priority ASC);

-- DELETE JOBS TRIGGER

CREATE OR REPLACE FUNCTION job_finished()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'completed' AND NEW.remove_on_complete THEN
        DELETE FROM Jobs WHERE id = NEW.id;
    ELSIF NEW.status = 'failed' AND NEW.remove_on_fail THEN
        IF NEW.retry AND NEW.attempts <= NEW.max_attempts THEN
			UPDATE Jobs
			SET status = 'pending',
				attempts = attempts + 1
			WHERE id = NEW.id;
		ELSE
        	DELETE FROM Jobs WHERE id = NEW.id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER job_finished_trigger
BEFORE UPDATE OF status ON Jobs
FOR EACH ROW
EXECUTE FUNCTION job_finished();

-- RETRY JOBS TRIGGER

CREATE OR REPLACE FUNCTION queue_trigger()
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.status = 'pending' THEN
		UPDATE Queues
		SET active_jobs = active_jobs + 1
		WHERE id = NEW.queue_id;
	ELSIF NEW.status = 'active' THEN
		UPDATE Queues
		SET active_jobs = active_jobs - 1
		WHERE id = NEW.queue_id;
	ELSIF NEW.status = 'completed' THEN
		UPDATE Queues
		SET completed_jobs = completed_jobs + 1
		WHERE id = NEW.queue_id;
	ELSIF NEW.status = 'failed' THEN
		UPDATE Queues
		SET failed_jobs = failed_jobs + 1
		WHERE id = NEW.queue_id;
	END IF;
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER queue_trigger
BEFORE UPDATE OF status ON jobs
FOR EACH ROW
EXECUTE FUNCTION queue_trigger();
