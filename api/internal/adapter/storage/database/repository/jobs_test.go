package repository

import (
	"testing"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TestIsJobRunning_Running tests when a job is running
func TestIsJobRunning_Running(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	jobName := "Photo_index"
	expectedJob := sqlmock.NewRows([]string{"name", "status", "last_run"}).
		AddRow(jobName, "RUNNING", time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "jobs" WHERE .+ ORDER BY .+ LIMIT \$2`, expectedJob)

	repo := NewJobRepository(env)
	running, err := repo.IsJobRunning(jobName)

	if err != nil {
		t.Errorf("error was not expected while checking job status: %s", err)
	}

	if !running {
		t.Error("expected job to be running, got false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestIsJobRunning_NotRunning tests when a job is not running
func TestIsJobRunning_NotRunning(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	jobName := "Photo_index"
	expectedJob := sqlmock.NewRows([]string{"name", "status", "last_run"}).
		AddRow(jobName, "IDLE", time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "jobs" WHERE .+ ORDER BY .+ LIMIT \$2`, expectedJob)

	repo := NewJobRepository(env)
	running, err := repo.IsJobRunning(jobName)

	if err != nil {
		t.Errorf("error was not expected while checking job status: %s", err)
	}

	if running {
		t.Error("expected job to not be running, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestIsJobRunning_NotFound tests when a job doesn't exist
func TestIsJobRunning_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	jobName := "NonexistentJob"

	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT \* FROM "jobs" WHERE .+ ORDER BY .+ LIMIT \$2`)

	repo := NewJobRepository(env)
	running, err := repo.IsJobRunning(jobName)

	if err == nil {
		t.Error("expected error when job not found, got nil")
	}

	if running {
		t.Error("expected job to not be running when not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateJobStatus_Success tests successful job status update
func TestUpdateJobStatus_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	jobName := "Photo_index"
	newStatus := "RUNNING"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "jobs" SET .+ WHERE name = \$\d+`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.UpdateJobStatus(jobName, newStatus)

	if err != nil {
		t.Errorf("error was not expected while updating job status: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateJobStatus_NotFound tests updating non-existent job
func TestUpdateJobStatus_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	jobName := "NonexistentJob"
	newStatus := "RUNNING"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "jobs" SET .+ WHERE name = \$\d+`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.UpdateJobStatus(jobName, newStatus)

	if err != nil {
		t.Errorf("error was not expected for non-existent job: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateAllJobsStatus_Success tests updating all running jobs
func TestUpdateAllJobsStatus_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	newStatus := "IDLE"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "jobs" SET .+ WHERE status = \$\d+`).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.UpdateAllJobsStatus(newStatus)

	if err != nil {
		t.Errorf("error was not expected while updating all jobs: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateAllJobsStatus_NoRunningJobs tests when no jobs are running
func TestUpdateAllJobsStatus_NoRunningJobs(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	newStatus := "IDLE"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "jobs" SET .+ WHERE status = \$\d+`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.UpdateAllJobsStatus(newStatus)

	if err != nil {
		t.Errorf("error was not expected when no running jobs: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateBaseJobs_Success tests creating base jobs
func TestCreateBaseJobs_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "jobs"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.CreateBaseJobs()

	if err != nil {
		t.Errorf("error was not expected while creating base jobs: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateBaseJobs_AlreadyExists tests when job already exists (upsert)
func TestCreateBaseJobs_AlreadyExists(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	// ON CONFLICT UPDATE ALL behavior - updates existing job
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "jobs"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewJobRepository(env)
	err := repo.CreateBaseJobs()

	if err != nil {
		t.Errorf("error was not expected for upsert: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
