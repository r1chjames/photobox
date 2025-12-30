package repository

import (
	. "gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/gorm/clause"
	"time"
)

type JobRepository struct {
	dbEnv *Env
}

func NewJobRepository(dbEnv *Env) *JobRepository {
	return &JobRepository{
		dbEnv,
	}
}

func (jr *JobRepository) IsJobRunning(name string) (bool, error) {
	var job domain.Job
	job.Name = name
	result := jr.dbEnv.Db.First(&job)
	err := HandleError(result)
	if err != nil {
		return false, err
	}
	if job.Status == "RUNNING" {
		return true, nil
	}
	return false, nil
}

func (jr *JobRepository) UpdateJobStatus(name string, status string) error {
	var job domain.Job
	job.Name = name
	result := jr.dbEnv.Db.Model(&job).Where("name = ?", name).Updates(domain.Job{Status: status, LastRun: time.Now()})
	return result.Error
}

// StartJobIfNotRunning atomically starts a job only if it's not already running
// Returns an error if the job is already running
func (jr *JobRepository) StartJobIfNotRunning(name string) error {
	// Atomic update: only set to RUNNING if current status is NOT_RUNNING
	result := jr.dbEnv.Db.Model(&domain.Job{}).
		Where("name = ? AND status != ?", name, "RUNNING").
		Updates(domain.Job{Status: "RUNNING", LastRun: time.Now()})

	if result.Error != nil {
		return result.Error
	}

	// If no rows were affected, the job is already running
	if result.RowsAffected == 0 {
		return domain.ErrJobAlreadyRunning
	}

	return nil
}

func (jr *JobRepository) UpdateAllJobsStatus(status string) error {
	result := jr.dbEnv.Db.Model(domain.Job{}).Where("status = ?", "RUNNING").Update("Status", status)
	return result.Error
}

func (jr *JobRepository) CreateBaseJobs() error {
	result := jr.dbEnv.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&domain.Job{Name: "Photo_index", Status: "NOT_RUNNING", LastRun: time.Now()})
	return result.Error
}
