package redis

import (
	"context"
	"encoding/json"

	pkg "github.com/jihedmastouri/saf/pkg/queueing"
	"github.com/redis/go-redis/v9"
)

type redisDriver struct {
	client *redis.Client
}

func init() {
	d := newRedisDriver()
	pkg.RegisterDriver("redis", d)
}

func newRedisDriver() *redisDriver {
	return &redisDriver{}
}

func (d *redisDriver) NewQueue() pkg.Queue {
	queue := NewQueue()
	return queue
}

func (d *redisDriver) Connect(addr string) error {
	d.client = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return nil
}

func (d *redisDriver) Disconnect() error {
	return d.client.Close()
}

func (d *redisDriver) NewWorkTrigger(queueName string, jobName string) (chan pkg.Job, error) {
	ctx := context.Background()
	ch := make(chan pkg.Job)

	go func() {
		for {
			jobJSON, err := d.client.LPop(ctx, jobName).Result()
			if err != nil {
				continue
			}

			var job pkg.Job
			json.Unmarshal([]byte(jobJSON), &job)

			ch <- job
		}
	}()

	return ch, nil
}

func (d *redisDriver) NewEventListner(eventName pkg.QueueEvents) (chan pkg.Job, error) {
	ctx := context.Background()
	ch := make(chan pkg.Job)
	pubsub := d.client.Subscribe(ctx, string(eventName))

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				continue
			}

			var job pkg.Job
			json.Unmarshal([]byte(msg.Payload), &job)

			ch <- job
		}
	}()

	return ch, nil
}
