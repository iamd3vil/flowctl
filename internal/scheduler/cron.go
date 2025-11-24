package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cvhariharan/flowctl/internal/repo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// syncScheduledFlows syncs scheduled flows from the database into the in-memory cache
func (s *Scheduler) syncScheduledFlows(ctx context.Context) error {
	scheduledFlows, err := s.store.GetScheduledFlows(ctx)
	if err != nil {
		return err
	}

	s.scheduledMu.Lock()
	defer s.scheduledMu.Unlock()

	s.scheduledFlows = make(map[string]repo.GetScheduledFlowsRow)
	for _, flow := range scheduledFlows {
		key := fmt.Sprintf("%s_%s_%s", flow.Slug, flow.Cron, flow.Timezone)
		s.scheduledFlows[key] = flow
	}

	s.logger.Debug("synced scheduled flows to cache", "count", len(s.scheduledFlows))
	return nil
}

// checkPeriodicTasks checks for flows with cron schedules that should run now
func (s *Scheduler) checkPeriodicTasks(ctx context.Context) error {
	s.scheduledMu.RLock()
	scheduledFlows := make([]repo.GetScheduledFlowsRow, 0, len(s.scheduledFlows))
	for _, flow := range s.scheduledFlows {
		scheduledFlows = append(scheduledFlows, flow)
	}
	s.scheduledMu.RUnlock()

	for _, flow := range scheduledFlows {
		// Check if this specific schedule should run now
		if flow.Cron != "" && s.shouldRunNow(flow.Cron, flow.Timezone) {
			if err := s.createImmediateTaskFromFlow(ctx, flow); err != nil {
				s.logger.Error("failed to create immediate task from scheduled flow", "flow", flow.Name, "error", err)
			}
		}
	}

	return nil
}

// shouldRunNow evaluates if a cron expression should execute in the current minute
func (s *Scheduler) shouldRunNow(cronExpr string, timezone string) bool {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		log.Printf("Failed to parse cron expression '%s': %v", cronExpr, err)
		return false
	}

	// Load the timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		s.logger.Error("failed to load timezone, using UTC", "timezone", timezone, "error", err)
		loc = time.UTC
	}

	// Convert current time to the schedules's timezone
	nowInTz := time.Now().In(loc)
	currentMinute := nowInTz.Truncate(time.Minute)

	lastMinute := currentMinute.Add(-time.Minute)
	nextRun := schedule.Next(lastMinute)

	// Task should run if the next scheduled time falls within the current minute
	return nextRun.Equal(currentMinute) || (nextRun.After(currentMinute) && nextRun.Before(currentMinute.Add(time.Minute)))
}

// createImmediateTaskFromFlow creates an immediate task from a scheduled flow
func (s *Scheduler) createImmediateTaskFromFlow(ctx context.Context, flow repo.GetScheduledFlowsRow) error {
	namespace, err := s.store.GetNamespaceByUUID(ctx, flow.NamespaceUuid)
	if err != nil {
		return fmt.Errorf("could not create periodic task %s: %w", flow.Name, err)
	}

	if s.flowLoader == nil {
		return fmt.Errorf("flow loader not configured")
	}

	schedulerFlow, err := s.flowLoader(ctx, flow.Slug, flow.NamespaceUuid.String())
	if err != nil {
		log.Printf("Failed to load flow %s: %v", flow.Slug, err)
		return err
	}

	input := applyDefaultInputValues(schedulerFlow.Inputs)

	payload := FlowExecutionPayload{
		Workflow:          schedulerFlow,
		Input:             input,
		StartingActionIdx: 0,
		ExecID:            uuid.NewString(),
		NamespaceID:       namespace.Uuid.String(),
		TriggerType:       TriggerTypeScheduled,
		UserUUID:          "00000000-0000-0000-0000-000000000000", // System user
	}

	_, err = s.QueueTask(ctx, payload)
	if err != nil {
		return err
	}

	s.logger.Info("created immediate task from scheduled flow", "flow", flow.Slug, "id", flow.ID, "cron", flow.Cron)
	return nil
}

// applyDefaultInputValues creates an input map using default values from flow inputs
func applyDefaultInputValues(inputs []Input) map[string]interface{} {
	result := make(map[string]interface{})
	for _, input := range inputs {
		if input.Default != "" {
			result[input.Name] = input.Default
		}
	}
	return result
}
