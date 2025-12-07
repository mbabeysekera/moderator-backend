package services

import (
	"context"

	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventService struct {
	eventRepo *repositories.EventReposiroty
}

func NewEventService(pool *pgxpool.Pool) *EventService {
	return &EventService{
		eventRepo: repositories.NewEventRepository(pool),
	}
}

func (es *EventService) CreateEvent(rc context.Context,
	event dto.EventRequest) (*dto.EventResponse, error) {

	newEvent := models.Event{
		Name:        event.Name,
		Description: event.Description,
		Target:      event.Target,
		Command:     event.Command,
		Status:      event.Status,
	}
	err := es.eventRepo.Create(rc, newEvent)
	if err != nil {
		return nil, err
	}
	return &dto.EventResponse{}, nil
}
