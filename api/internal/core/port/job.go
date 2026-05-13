package port

//go:generate mockgen -source=job.go -destination=mock/job.go -package=mock

// JobRepository is an interface for interacting with utility-related business logic
type JobRepository interface {
	// IsJobRunning returns information about whether a job is currently running
	IsJobRunning(name string) (bool, error)
	// UpdateJobStatus updates a single job
	UpdateJobStatus(name string, status string) error
	// UpdateAllJobsStatus sets status of all running jobs
	UpdateAllJobsStatus(status string) error
	// CreateBaseJobs creates the jobs required for a new instance
	CreateBaseJobs() error
	// StartJobIfNotRunning atomically starts a job only if it's not already running
	StartJobIfNotRunning(name string) error
}

// AlbumService is an interface for interacting with Album-related business logic
type JobService interface {
	// IsJobRunning returns information about whether a job is currently running
	IsJobRunning(name string) (bool, error)
	// UpdateAllJobsStatus sets status of all running jobs
	UpdateAllJobsStatus(status string) error
	// JobStart updates a job to RUNNING
	JobStart(name string) error
	// JobComplete updates a job to NOT_RUNNING
	JobComplete(name string) error
	// CreateBaseJobs creates the jobs required for a new instance
	CreateBaseJobs() error
	// StartJobIfNotRunning atomically starts a job only if it's not already running
	StartJobIfNotRunning(name string) error
}
