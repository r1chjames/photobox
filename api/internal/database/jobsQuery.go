package database

import (
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"time"
)

func IsJobRunning(jobName string) (bool, error) {
	var job Job
	job.Name = jobName
	result := dbConn.First(&job)
	if job.Status == "RUNNING" {
		return true, result.Error
	}
	return false, result.Error
}

func updateJobStatus(jobName string, status string) error {
	var job Job
	job.Name = jobName
	result := dbConn.Model(&job).Where("name = ?", jobName).Updates(Job{Status: status, LastRun: time.Now()})
	return result.Error
}

func StopAllRunningJobs(status string) error {
	result := dbConn.Model(Job{}).Where("status = ?", "RUNNING").Update("Status", status)
	return result.Error
}

func JobStarting(jobName string) error {
	return updateJobStatus(jobName, "RUNNING")
}

func JobCompleted(jobName string) error {
	return updateJobStatus(jobName, "NOT_RUNNING")
}


