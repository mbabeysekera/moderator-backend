package constants

type EventCode int

const (
	EventNew EventCode = iota
	EventPending
	EventCompleted
	EventFailed
	EventStale
)

var EventStatus = map[EventCode]string{
	EventNew:       "NEW",
	EventPending:   "PENDING",
	EventCompleted: "COMPLTED",
	EventFailed:    "FAILED",
	EventStale:     "STALE",
}
