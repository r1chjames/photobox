package database

import (
	"database/sql"
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
	"time"
)

func transformRowsIntoBooleanArray(results *sql.Rows) interface{} {
	defer results.Close()
	var records []bool
	for results.Next() {
		var record bool
		err := results.Scan(
			&record)
		if err != nil {
			log.Fatal(fmt.Sprintf("An error parsing settings record: %s", err))
		}
		records = append(records, record)
	}
	return records
}

func IsJobRunning(appConfig AppConfig, jobName string) bool {
	allJobsFields := "RUNNING"
	query := fmt.Sprintf("SELECT %s FROM jobs WHERE job = '%s'", allJobsFields, jobName)

	queryResults := executeTransformDbQuery(appConfig, query, transformRowsIntoBooleanArray).([]bool)
	if len(queryResults) > 0 {
		return queryResults[0]
	}

	return true
}

func updateJobStatus(appConfig AppConfig, jobName string, status bool) {
	insertStatement := fmt.Sprintf("UPDATE jobs SET running = %t, lastRun = '%s' where job = '%s'",
		status,
		time.Now().Format("2006-01-02 15:04:05"),
		jobName)

	executeDbInsert(appConfig, insertStatement)
}

func UpdateAllJobStatus(appConfig AppConfig, status bool) {
	insertStatement := fmt.Sprintf("UPDATE jobs SET running = %t",
		status)
	executeDbInsert(appConfig, insertStatement)
}

func JobStarting(appConfig AppConfig, jobName string) {
	updateJobStatus(appConfig, jobName, true)
}

func jobCompleted(appConfig AppConfig, jobName string) {
	updateJobStatus(appConfig, jobName, false)
}


