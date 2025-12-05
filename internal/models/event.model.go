package models

import (
	"context"
	"log"
	"time"

	"coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/db"
)

type Event struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Target      string              `json:"target"`
	Command     string              `json:"command"`
	Status      constants.EventCode `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	ProcessedAt *time.Time          `json:"processed_at,omitempty"`
}

func (e *Event) Save(c context.Context) error {
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
	tag, err := db.DB.Exec(c, evenCreate,
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

func (e *Event) GetEventsByStatus(c context.Context, status constants.EventCode) ([]Event, error) {
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
	events, err := db.DB.Query(c, getEventByStatus, status)
	if err != nil {
		return nil, err
	}
	defer events.Close()
	var allEvents = []Event{}
	for events.Next() {
		var event Event
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
