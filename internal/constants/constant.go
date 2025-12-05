package constants

type EventCode string

const (
	EventNew       EventCode = "NEW"
	EventPending   EventCode = "PENDING"
	EventCompleted EventCode = "COMPLTED"
	EventFailed    EventCode = "FAILED"
	EventStale     EventCode = "STALE"
)
