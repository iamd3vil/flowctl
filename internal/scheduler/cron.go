package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// syncScheduledJobs syncs scheduled jobs from the job syncer into the cache
func (s *Scheduler) syncScheduledJobs(ctx context.Context) error {
	if s.jobSyncer == nil {
		return nil
	}

	jobs, err := s.jobSyncer(ctx)
	if err != nil {
		return err
	}

	s.scheduledMu.Lock()
	defer s.scheduledMu.Unlock()

	s.scheduledJobs = make(map[string]ScheduledJob)
	for _, job := range jobs {
		key := job.ID + "_" + job.Cron + "_" + job.Timezone
		s.scheduledJobs[key] = job
	}

	s.logger.Debug("synced scheduled jobs to cache", "count", len(s.scheduledJobs))
	return nil
}

// checkPeriodicTasks checks for scheduled jobs that should run now
func (s *Scheduler) checkPeriodicTasks(ctx context.Context) error {
	if s.jobSyncer == nil {
		return nil
	}

	s.scheduledMu.RLock()
	jobs := make([]ScheduledJob, 0, len(s.scheduledJobs))
	for _, job := range s.scheduledJobs {
		jobs = append(jobs, job)
	}
	s.scheduledMu.RUnlock()

	for _, job := range jobs {
		if job.Cron != "" && s.shouldRunNow(job.Cron, job.Timezone) {
			// Generate a new execID for each execution
			execID := uuid.NewString()

			if _, err := s.QueueTask(ctx, job.PayloadType, execID, job.Payload); err != nil {
				s.logger.Error("failed to queue scheduled job", "job", job.Name, "error", err)
			} else {
				s.logger.Info("queued scheduled job", "job", job.Name, "id", job.ID, "execID", execID, "cron", job.Cron)
			}
		}
	}

	return nil
}

// shouldRunNow evaluates if a cron expression should execute in the current minute
func (s *Scheduler) shouldRunNow(cronExpr string, timezone string) bool {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		s.logger.Error("failed to parse cron expression", "cron", cronExpr, "error", err)
		return false
	}

	// Load the timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		s.logger.Error("failed to load timezone, using UTC", "timezone", timezone, "error", err)
		loc = time.UTC
	}

	// Convert current time to the schedule's timezone
	nowInTz := time.Now().In(loc)
	currentMinute := nowInTz.Truncate(time.Minute)

	lastMinute := currentMinute.Add(-time.Minute)
	nextRun := schedule.Next(lastMinute)

	// Task should run if the next scheduled time falls within the current minute
	return nextRun.Equal(currentMinute) || (nextRun.After(currentMinute) && nextRun.Before(currentMinute.Add(time.Minute)))
}
