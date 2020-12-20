package components

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	"gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

var c *cron.Cron

func InitScheduler(appConfig types.AppConfig) {
	c = cron.New(cron.WithLocation(appConfig.Timezone))
	c.Start()
}

func AddScheduledJobs(appConfig types.AppConfig, dbEnv *database.Env) {
	setting, err := dbEnv.GetSetting("index_frequency_cron")
	if err != nil {
		log.Print("unable to get photo index cron expression from database. Index scheduling will not be enabled")
	}

	_, err = c.AddFunc(setting.Value, func() { PerformPhotoIndex(appConfig, dbEnv) })
	if err != nil {
		log.Printf("unable to add job schedule for %s. Parsed CRON expression: %s. Check CRON expression in settings", setting.Key, setting.Value)
	}
	log.Print(c.Entries())
}

func UpdateJobSchedule(dbEnv *database.Env) {

}

func StopAllRunningJobs(dbEnv *database.Env) {
	res := c.Stop()
	if res.Err() != nil {
		log.Print(fmt.Sprintf("unable to stop scheduler: %s", res.Err()))
	}
	dbEnv.UpdateAllJobsStatus("NOT_RUNNING")
}