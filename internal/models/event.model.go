package models

import (
	"log"
	"time"

	"coolbreez.lk/moderator/internal/constants"
)

type Event struct {
	ID          string
	Name        string
	Description string
	Target      string
	Command     string
	Status      constants.EventCode
	CreatedAt   time.Time
	ProcessedAt time.Time
}

// func New(name string, description string, target string, command string) *Event {
// 	return &Event{
// 		ID:          "0000000-0000-0000-0000-0000000",
// 		Name:        name,
// 		Description: description,
// 		Target:      target,
// 		Command:     command,
// 		CreatedAt:   time.Now(),
// 	}
// }

func (e *Event) Save() {
	log.Println("Event Saved...")
}
