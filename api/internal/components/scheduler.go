package components

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"log"
)

type Scheduler struct {
	cron     *cron.Cron
	jobSvc   port.JobService
	utilSvc  port.UtilityService
	photoSvc port.PhotoService
	config   appconfig.AppConfig
}

func NewScheduler(utilityService port.UtilityService, jobService port.JobService, photoSvc port.PhotoService, appConfig appconfig.AppConfig) *Scheduler {
	var scheduler = &Scheduler{
		cron.New(cron.WithLocation(appConfig.Timezone)),
		jobService,
		utilityService,
		photoSvc,
		appConfig,
	}

	scheduler.cron.Start()
	return scheduler
}

func (s *Scheduler) AddScheduledJobs() {
	setting, err := s.utilSvc.GetSetting("index_frequency_cron")
	if err != nil {
		log.Print("unable to get photo index cron expression from database. Index scheduling will not be enabled")
	}

	_, err = s.cron.AddFunc(setting.Value, func() {
		s.photoSvc.PerformPhotoIndex()
	})
	if err != nil {
		log.Printf("unable to add job schedule for %s. Parsed CRON expression: %s. Check CRON expression in settings", setting.Key, setting.Value)
	}
	log.Print(s.cron.Entries())
}

func UpdateJobSchedule() {

}

func (s *Scheduler) StopAllRunningJobs() {
	res := s.cron.Stop()
	if res.Err() != nil {
		log.Print(fmt.Sprintf("unable to stop scheduler: %s", res.Err()))
	}
	s.jobSvc.UpdateAllJobsStatus("NOT_RUNNING")
}
