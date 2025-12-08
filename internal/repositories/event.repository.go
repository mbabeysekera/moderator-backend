package repositories

import (
	"context"
	"log"

	"coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventReposiroty struct {
	pool *pgxpool.Pool
}

func NewEventRepository(dbPool *pgxpool.Pool) *EventReposiroty {
	return &EventReposiroty{pool: dbPool}
}

func (eventRepo *EventReposiroty) Create(c context.Context, e *models.Event) error {

	evenCreate := `
		INSERT INTO events (
			name, 
			description, 
			target, 
			command, 
			status
		) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	tag, err := eventRepo.pool.Exec(c, evenCreate,
		e.Name,
		e.Description,
		e.Target,
		e.Command,
		e.Status,
	)
	if err != nil {
		return err
	}
	log.Printf("Insert data into event table. No. of rows created: %d", tag.RowsAffected())
	return nil
}

func (eventRepo *EventReposiroty) GetEventsByStatus(c context.Context,
	status constants.EventCode) ([]models.Event, error) {

	getEventByStatus := `SELECT 
		id,
		name,
		description,
		target,
		command,
		status,
		created_at,
		processed_at
	FROM events WHERE status = $1`
	events, err := eventRepo.pool.Query(c, getEventByStatus, status)
	if err != nil {
		return nil, err
	}
	defer events.Close()
	var allEvents = []models.Event{}
	for events.Next() {
		var event models.Event
		err = events.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Target,
			&event.Command,
			&event.Status,
			&event.CreatedAt,
			&event.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, event)
	}
	err = events.Err()
	if err != nil {
		return nil, err
	}
	return allEvents, nil
}
