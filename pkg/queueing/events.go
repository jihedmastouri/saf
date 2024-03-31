package pkg

type QueueEvents string

const (
	JobEnqueued  QueueEvents = "job:enqueued"
	JobFailed    QueueEvents = "job:failed"
	JobCompleted QueueEvents = "job:completed"
)

type GlobalEvents string

const (
	NewQueueCreated GlobalEvents = "queue:created"
)
