package pkg

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCompleted JobStatus = "completed"
	JobStatusRunning   JobStatus = "running"
	JobStatusCanceled  JobStatus = "canceled"
	JobStatusUnknown   JobStatus = "unknown"
)

type JobInfo struct {
	// The name of the job. It should be unique for each job as it is used to identify the job.
	Name string
	// The data to be passed to the job
	Data map[string]any
	// The minimum delay in seconds before the job is executed
	Delay int
	// The maximum number of retries before the job is marked as failed
	MaxRetries int
	// The maximum number of seconds the job is allowed to run
	Timeout int
	// The priority of the job (0 is the lowest priority)
	Priority int
}

type Job interface {
	// Mark the job as done
	Done() error
	// Mark the job as failed
	Fail() error
	// Remove the job from the queue
	Remove() error
	// Get the job ID
	GetID() string
	// Get the status of the job
	GetStatus() JobStatus
	// Get the job details
	GetInfo() JobInfo
}

func NewJob(name string, opts ...JobOptions) JobInfo {
	job := JobInfo{Name: name}
	for _, opt := range opts {
		opt(&job)
	}
	return job
}

type JobOptions func(*JobInfo)

func JobWithTimeout(timeout int) func(*JobInfo) {
	return func(ji *JobInfo) {
		ji.Timeout = timeout
	}
}

func JobWithDelay(delay int) func(*JobInfo) {
	return func(ji *JobInfo) {
		ji.Timeout = delay
	}
}

func JobWithMaxRetries(maxRetries int) func(*JobInfo) {
	return func(ji *JobInfo) {
		ji.MaxRetries = maxRetries
	}
}

func JobWithPriority(priority int) func(*JobInfo) {
	return func(ji *JobInfo) {
		ji.Priority = priority
	}
}

func JobWithData(data map[string]any) func(*JobInfo) {
	return func(ji *JobInfo) {
		ji.Data = data
	}
}

type Worker struct {
	ch       chan Job
	pause    chan struct{}
	callback func(Job)
}

type WorkerOptions struct {
	queueName  string
	maxWorkers int
	jobName    string
}

// Initialize a new worker
func (qm *QueueManager) NewWorker(queueName string, jobName string, fn func(job Job)) (*Worker, error) {
	c, err := qm.driver.NewWorkTrigger(queueName, jobName)
	if err != nil {
		return nil, err
	}
	return &Worker{ch: c, callback: fn, pause: make(chan struct{})}, nil
}

// Start/Resume the worker
func (w *Worker) Start() {
	for {
		select {
		case <-w.pause:
			return
		case job, ok := <-w.ch:
			if ok {
				go w.callback(job)
			}
		}
	}
}

// Stop the worker completely
func (w *Worker) Stop() {
	close(w.ch)
}

// Pause the worker temporarily
func (w *Worker) Pause() {
	w.pause <- struct{}{}
}

type Flow struct {
	Children *[]Flow
	Details  JobInfo
}

func addFlow(queue Queue, flow Flow) error {
	var jobs []Job

	job, err := queue.Schedule(flow.Details)
	if err != nil {
		for _, j := range jobs {
			j.Remove()
		}
		return err
	}
	jobs = append(jobs, job)

	for _, child := range *flow.Children {
		job, err := queue.Schedule(child.Details)
		if err != nil {
			for _, j := range jobs {
				j.Remove()
			}
			return err
		}
		jobs = append(jobs, job)
	}

	return nil
}
