package enums

type EventCode string

const (
	EventNew       EventCode = "NEW"
	EventPending   EventCode = "PENDING"
	EventCompleted EventCode = "COMPLTED"
	EventFailed    EventCode = "FAILED"
	EventStale     EventCode = "STALE"
)

type RequestStatus string

const (
	RequestSuccess RequestStatus = "SUCCESS"
	RequestFailed  RequestStatus = "FAILED"
)

type UserRole string

const (
	RoleUser      UserRole = "USER"
	RoleAdmin     UserRole = "ADMIN"
	RoleModerator UserRole = "MODERATOR"
	RoleSystem    UserRole = "SYSTEM"
)

type ProductCategory string

const (
	Earrings  ProductCategory = "EARRING"
	Rings     ProductCategory = "RING"
	Necklace  ProductCategory = "NECKLACE"
	Bracelets ProductCategory = "BRACELET"
	Pendants  ProductCategory = "PENDANT"
)
