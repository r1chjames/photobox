package components

import (
	"log/slog"

	"github.com/robfig/cron/v3"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
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
		slog.Warn("Unable to get photo index cron expression from database, index scheduling will not be enabled")
		return
	}

	_, err = s.cron.AddFunc(setting.Value, func() {
		s.photoSvc.PerformPhotoIndex()
	})
	if err != nil {
		slog.Error("Unable to add job schedule, check CRON expression in settings", "key", setting.Key, "cron", setting.Value, "error", err)
	}

	// Schedule AI analysis job to run every hour
	_, err = s.cron.AddFunc("0 * * * *", func() {
		if err := s.photoSvc.AnalyzeExistingPhotos(); err != nil {
			slog.Error("AI analysis job failed", "error", err)
		}
	})
	if err != nil {
		slog.Error("Unable to add AI analysis job schedule", "error", err)
	}

	slog.Info("Scheduled jobs loaded", "entries", len(s.cron.Entries()))
}

func UpdateJobSchedule() {

}

func (s *Scheduler) StopAllRunningJobs() {
	res := s.cron.Stop()
	if res.Err() != nil {
		slog.Error("Unable to stop scheduler", "error", res.Err())
	}
	s.jobSvc.UpdateAllJobsStatus("NOT_RUNNING")
}
