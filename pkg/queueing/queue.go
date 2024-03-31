package pkg

type Driver interface {
	Connect(string) error
	Disconnect() error
	NewQueue() Queue
	NewWorkTrigger(string, string) (chan Job, error)
	NewEventListner(QueueEvents) (chan Job, error)
}

var drivers = make(map[string]Driver)

type QueueManager struct {
	driver Driver
}

type QueueOptions struct {
	Name             string
	Description      string
	MaxActiveJobs    int
	MaxCompletedJobs int
	MaxFailedJobs    int
}

type Queue interface {
	// Get the name of the Queue
	Name() string
	// Enqueue a job
	Schedule(JobInfo) (Job, error)
	// Listen to events from the queue
	NewEventListner(QueueEvents) (chan QueueEvents, error)
	// Return the number of all jobs in the queue (all status: pending, completed...).
	Len() (int, error)
	// Return the number of jobs in the queue with the given status
	Count(JobStatus) (int, error)
	// Drain the queue, Removing all jobs
	Drain() error
}

// Register a driver for use as a backend for `saf`.
func RegisterDriver(name string, driver Driver) {
	drivers[name] = driver
}

// NewQueueManager creates a new QueueManager and connects to the backend using the driver.
func NewQueueManager(driver string, addr string) (*QueueManager, error) {
	drv, ok := drivers[driver]
	if !ok {
		panic("driver not found")
	}
	if err := drv.Connect(addr); err != nil {
		return nil, err
	}
	return &QueueManager{
		driver: drv,
	}, nil
}

type NewQueueOptions func(*QueueOptions)

func (qm *QueueManager) NewQueue(name string, opts ...NewQueueOptions) Queue {
	return qm.driver.NewQueue()
}

func WithMaxActiveJobs(n int) func(*QueueOptions) {
	return func(queueOptions *QueueOptions) {
		queueOptions.MaxActiveJobs = n
	}
}

func WithMaxCompletedJobs(n int) func(*QueueOptions) {
	return func(queueOptions *QueueOptions) {
		queueOptions.MaxCompletedJobs = n
	}
}

func WithMaxFailedJobs(n int) func(*QueueOptions) {
	return func(queueOptions *QueueOptions) {
		queueOptions.MaxFailedJobs = n
	}
}
