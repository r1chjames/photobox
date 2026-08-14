package components

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
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
	// Schedule AI analysis job to run every hour (independent of photo index setting)
	_, err := s.cron.AddFunc("0 * * * *", func() {
		if err := s.jobSvc.StartJobIfNotRunning("AI_analysis"); err != nil {
			if errors.Is(err, domain.ErrJobAlreadyRunning) {
				return
			}
			slog.Error("Failed to start AI analysis job", "error", err)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				slog.Error("AI analysis job panicked", "recover", r)
			}
			if err := s.jobSvc.JobComplete("AI_analysis"); err != nil {
				slog.Error("Unable to complete AI analysis job", "error", err)
			}
		}()
		if err := s.photoSvc.AnalyzeExistingPhotos(); err != nil {
			slog.Error("AI analysis job failed", "error", err)
		}
	})
	if err != nil {
		slog.Error("Unable to add AI analysis job schedule", "error", err)
	}

	// Schedule quality backfill (daily, after the index run). Scores any
	// photos missed by indexing (pre-existing libraries) and keeps scores
	// fresh as files are touched.
	_, err = s.cron.AddFunc("0 3 * * *", func() {
		scored, err := s.photoSvc.ScoreAllPhotoQuality(200)
		if err != nil {
			slog.Error("Quality backfill failed", "error", err)
			return
		}
		if scored > 0 {
			slog.Info("Quality backfill complete", "scored", scored)
		}
	})
	if err != nil {
		slog.Error("Unable to add quality backfill job schedule", "error", err)
	}

	// Schedule trash retention cleanup to run daily. The retention period is
	// read from settings each run so admin changes apply without a restart.
	_, err = s.cron.AddFunc("0 2 * * *", func() {
		retentionDays, err := s.trashRetentionDays()
		if err != nil {
			slog.Error("Trash retention cleanup skipped: unable to read retention setting", "error", err)
			return
		}
		if retentionDays <= 0 {
			slog.Info("Trash retention cleanup skipped: retention disabled (0 days)")
			return
		}
		cutoff := time.Now().AddDate(0, 0, -retentionDays)
		deleted, err := s.photoSvc.PurgeExpiredTrash(cutoff)
		if err != nil {
			slog.Error("Trash retention cleanup failed", "error", err)
			return
		}
		slog.Info("Trash retention cleanup complete", "deleted", deleted)
	})
	if err != nil {
		slog.Error("Unable to add trash retention job schedule", "error", err)
	}

	setting, err := s.utilSvc.GetSetting("index_frequency_cron")
	if err != nil {
		slog.Warn("Unable to get photo index cron expression from database, index scheduling will not be enabled")
		return
	}

	_, err = s.cron.AddFunc(setting.Value, func() {
		s.photoSvc.PerformPhotoIndex(context.Background())
	})
	if err != nil {
		slog.Error("Unable to add job schedule, check CRON expression in settings", "key", setting.Key, "cron", setting.Value, "error", err)
	}

	slog.Info("Scheduled jobs loaded", "entries", len(s.cron.Entries()))
}

func UpdateJobSchedule() {

}

// trashRetentionDays reads the trash retention period from settings, falling
// back to the configured default when the setting is missing or unparseable.
func (s *Scheduler) trashRetentionDays() (int, error) {
	setting, err := s.utilSvc.GetSetting("trash_retention_days")
	if err == nil && setting != nil && setting.Value != "" {
		days, err := strconv.Atoi(setting.Value)
		if err == nil {
			return days, nil
		}
		slog.Warn("Invalid trash_retention_days setting, using configured default", "value", setting.Value)
	}
	return s.config.TrashRetentionDays, nil
}

func (s *Scheduler) StopAllRunningJobs() {
	res := s.cron.Stop()
	if res.Err() != nil {
		slog.Error("Unable to stop scheduler", "error", res.Err())
	}
	if err := s.jobSvc.UpdateAllJobsStatus("NOT_RUNNING"); err != nil {
		slog.Error("Unable to reset job statuses", "error", err)
	}
}
