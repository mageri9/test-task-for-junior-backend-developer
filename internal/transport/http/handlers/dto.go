package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type RecurrenceRuleDTO struct {
	Type       string   `json:"type"`
	Interval   int      `json:"interval,omitempty"`
	DayOfMonth int      `json:"day_of_month,omitempty"`
	Dates      []string `json:"dates,omitempty"`
	Parity     string   `json:"parity,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *RecurrenceRuleDTO `json:"recurrence,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Recurrence  *RecurrenceRuleDTO `json:"recurrence,omitempty"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.Recurrence != nil {
		dto.Recurrence = &RecurrenceRuleDTO{
			Type:       string(task.Recurrence.Type),
			Interval:   task.Recurrence.Interval,
			DayOfMonth: task.Recurrence.DayOfMonth,
			Dates:      task.Recurrence.Dates,
			Parity:     task.Recurrence.Parity,
		}
	}

	return dto
}
