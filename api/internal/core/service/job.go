package service

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

/**
 * JobService implements port.JobService interface
 * and provides an access to the utility repository
 */
type JobService struct {
	repo port.JobRepository
}

// NewJobService creates a new job service instance
func NewJobService(repo port.JobRepository) *JobService {
	return &JobService{
		repo,
	}
}

func (js *JobService) IsJobRunning(name string) (bool, error) {
	return js.repo.IsJobRunning(name)
}

func (js *JobService) GetAllJobs() ([]domain.Job, error) {
	return js.repo.GetAllJobs()
}

func (js *JobService) UpdateAllJobsStatus(status string) error {
	return js.repo.UpdateAllJobsStatus(status)
}

func (js *JobService) JobStart(name string) error {
	return js.repo.UpdateJobStatus(name, "RUNNING")
}

func (js *JobService) JobComplete(name string) error {
	return js.repo.UpdateJobStatus(name, "NOT_RUNNING")
}

func (js *JobService) CreateBaseJobs() error {
	return js.repo.CreateBaseJobs()
}

func (js *JobService) StartJobIfNotRunning(name string) error {
	return js.repo.StartJobIfNotRunning(name)
}
