package jobs

type Queue struct {
}

func NewQueue() *Queue {
	return &Queue{}
}

var queuedJobs []*Job

func (q *Queue) GetJobs() ([]*Job, error) {
	return queuedJobs, nil
}

func (q *Queue) AddJob(job *Job) error {
	queuedJobs = append(queuedJobs, job)
	return nil
}
