package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func convertRecurrenceRuleDTO(dto *RecurrenceRuleDTO) (*taskdomain.RecurrenceRule, error) {
	if dto == nil {
		return nil, nil
	}

	rule := &taskdomain.RecurrenceRule{
		Type:       taskdomain.RecurrenceType(dto.Type),
		Interval:   dto.Interval,
		DayOfMonth: dto.DayOfMonth,
		Dates:      dto.Dates,
		Parity:     dto.Parity,
	}

	switch rule.Type {
	case taskdomain.RecurrenceDaily:
		if rule.Interval <= 0 {
			return nil, fmt.Errorf("%w: daily interval must be > 0", ErrInvalidInput)
		}

	case taskdomain.RecurrenceMonthly:
		if rule.DayOfMonth < 1 || rule.DayOfMonth > 30 {
			return nil, fmt.Errorf("%w: monthly day must be between 1 and 30", ErrInvalidInput)
		}

	case taskdomain.RecurrenceSpecificDates:
		if len(rule.Dates) == 0 {
			return nil, fmt.Errorf("%w: specific dates cannot be empty", ErrInvalidInput)
		}

		for _, date := range rule.Dates {
			if _, err := time.Parse("2006-01-02", date); err != nil {
				return nil, fmt.Errorf("%w: invalid date format '%s', expected YYYY-MM-DD", ErrInvalidInput, date)
			}
		}

	case taskdomain.RecurrenceParity:
		if rule.Parity != "even" && rule.Parity != "odd" {
			return nil, fmt.Errorf("%w: parity must be 'even' or 'odd'", ErrInvalidInput)
		}

	default:
		return nil, fmt.Errorf("%w: unknown recurrence type '%s'", ErrInvalidInput, rule.Type)
	}

	return rule, nil
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	recurrenceRule, err := convertRecurrenceRuleDTO(input.Recurrence)
	if err != nil {
		return nil, err
	}

	now := s.now()
	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
		Recurrence:  recurrenceRule,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	recurrenceRule, err := convertRecurrenceRuleDTO(input.Recurrence)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
		Recurrence:  recurrenceRule,
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}
