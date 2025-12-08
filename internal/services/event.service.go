package services

import (
	"context"

	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
)

type EventServiceImpl struct {
	eventRepo *repositories.EventReposiroty
}

func NewEventService(repo *repositories.EventReposiroty) *EventServiceImpl {
	return &EventServiceImpl{
		eventRepo: repo,
	}
}

func (es *EventServiceImpl) CreateEvent(rc context.Context,
	event dto.EventRequest) error {

	err := es.eventRepo.Create(rc, &models.Event{
		Name:        event.Name,
		Description: event.Description,
		Target:      event.Target,
		Command:     event.Command,
		Status:      event.Status,
	})
	if err != nil {
		return err
	}
	return nil
}
