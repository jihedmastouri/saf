package redis

import pkg "github.com/jihedmastouri/saf/pkg/queueing"

// Define Queues struct
type Q struct{}

func NewQueue() *Q {
	return &Q{}
}

func (q *Q) Name() string {
	return ""
}

func (q *Q) Schedule(job pkg.JobInfo) (pkg.Job, error) {
	return nil, nil
}

func (q *Q) NewEventListner(eventName pkg.QueueEvents) (chan pkg.QueueEvents, error) {
	ch := make(chan pkg.QueueEvents)
	return ch, nil
}

func (q *Q) Len() (int, error) {
	return 0, nil
}

func (q *Q) Count(status pkg.JobStatus) (int, error) {
	return 0, nil
}

func (q *Q) Drain() error {
	return nil
}
