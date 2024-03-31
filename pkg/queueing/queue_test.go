package pkg_test

import (
	"github.com/jihedmastouri/saf/pkg/queueing"
	"testing"

	_ "github.com/jihedmastouri/saf/drivers/redis"
)

func TestQueue(t *testing.T) {
	qm, err := pkg.NewQueueManager("redis", "localhost:6379")
	if err != nil {
		t.Error(err)
	}

	jobInfo := pkg.NewJob("test")

	q := qm.NewQueue("test")
	_, err = q.Schedule(jobInfo)
	if err != nil {
		t.Error(err)
	}

	qm.NewWorker("test", "test", func(job pkg.Job) {
		job.Done()
	})
}
