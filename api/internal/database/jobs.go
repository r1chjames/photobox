package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/gorm/clause"
	"time"
)

func (dbEnv *Env) IsJobRunning(jobName string) (bool, error) {
	var job Job
	job.Name = jobName
	result := dbEnv.db.First(&job)
	if job.Status == "RUNNING" {
		return true, result.Error
	}
	return false, result.Error
}

func (dbEnv *Env) updateJobStatus(jobName string, status string) error {
	var job Job
	job.Name = jobName
	result := dbEnv.db.Model(&job).Where("name = ?", jobName).Updates(Job{Status: status, LastRun: time.Now()})
	return result.Error
}

func (dbEnv *Env) StopAllRunningJobs(status string) error {
	result := dbEnv.db.Model(Job{}).Where("status = ?", "RUNNING").Update("Status", status)
	return result.Error
}

func (dbEnv *Env) JobStarting(jobName string) error {
	return dbEnv.updateJobStatus(jobName, "RUNNING")
}

func (dbEnv *Env) JobCompleted(jobName string) error {
	return dbEnv.updateJobStatus(jobName, "NOT_RUNNING")
}

func (dbEnv *Env) createBaseJobs() {
	dbEnv.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&Job{Name: "Photo_index", Status: "NOT_RUNNING", LastRun: time.Now()})
}
