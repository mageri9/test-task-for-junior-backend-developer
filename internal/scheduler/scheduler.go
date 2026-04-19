package scheduler

import (
	"context"
	"log/slog"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	ListRecurring(ctx context.Context) ([]taskdomain.Task, error)
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
}

type Scheduler struct {
	repo   Repository
	logger *slog.Logger
	ticker *time.Ticker
	interval time.Duration
}

func New(repo Repository, logger *slog.Logger, interval time.Duration) *Scheduler {
	return &Scheduler{
		repo:   repo,
		logger: logger,
		ticker: time.NewTicker(interval),
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
    s.logger.Info("scheduler started", "interval", s.interval)

	go func() {
		defer s.ticker.Stop()

		for {
			select {
			case <-s.ticker.C:
				s.processTasks(ctx)

			case <-ctx.Done():
				s.logger.Info("scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) processTasks(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("scheduler panic recovered", "error", r)
		}
	}()

	s.logger.Info("scheduler tick started")

	templates, err := s.repo.ListRecurring(ctx)
	if err != nil {
		s.logger.Error("failed to list recurring tasks", "error", err)
		return
	}

	today := time.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	dayOfMonth := today.Day()

	created := 0
	for _, template := range templates {
		if template.Recurrence == nil {
			continue
		}

		if !s.shouldCreateToday(template.Recurrence, todayDate, dayOfMonth) {
			continue
		}

		newTask := &taskdomain.Task{
			Title:         template.Title,
			Description:   template.Description,
			Status:        taskdomain.StatusNew,
			CreatedAt:     today,
			UpdatedAt:     today,
			SourceTaskID:  &template.ID,
			ScheduledDate: &todayDate,
			Recurrence:    nil,
		}

		_, err := s.repo.Create(ctx, newTask)
		if err != nil {
			s.logger.Warn("failed to create scheduled task",
				"source_id", template.ID,
				"date", todayDate.Format("2006-01-02"),
				"error", err)
			continue
		}

		created++
		s.logger.Info("scheduled task created",
			"source_id", template.ID,
			"title", template.Title,
			"date", todayDate.Format("2006-01-02"))
	}

	s.logger.Info("scheduler tick finished", "tasks_created", created)
}

func (s *Scheduler) shouldCreateToday(rule *taskdomain.RecurrenceRule, today time.Time, dayOfMonth int) bool {
	switch rule.Type {
	case taskdomain.RecurrenceDaily:
		dayOfYear := today.YearDay()
		return dayOfYear%rule.Interval == 0

	case taskdomain.RecurrenceMonthly:
		return dayOfMonth == rule.DayOfMonth

	case taskdomain.RecurrenceSpecificDates:
		todayStr := today.Format("2006-01-02")
		for _, d := range rule.Dates {
			if d == todayStr {
				return true
			}
		}
		return false

	case taskdomain.RecurrenceParity:
		if rule.Parity == "even" {
			return dayOfMonth%2 == 0
		}
		return dayOfMonth%2 != 0

	default:
		return false
	}
}