package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *taskusecase.RecurrenceRuleDTO `json:"recurrence,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Recurrence  *taskusecase.RecurrenceRuleDTO `json:"recurrence,omitempty"`
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
		dto.Recurrence = &taskusecase.RecurrenceRuleDTO{
			Type:       string(task.Recurrence.Type),
			Interval:   task.Recurrence.Interval,
			DayOfMonth: task.Recurrence.DayOfMonth,
			Dates:      task.Recurrence.Dates,
			Parity:     task.Recurrence.Parity,
		}
	}

	return dto
}
