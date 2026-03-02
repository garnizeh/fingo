package auditbus

import (
	"encoding/json"
	"time"

	"github.com/garnizeh/fingo/business/types/domain"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// Audit represents information about an individual audit record.
type Audit struct {
	Timestamp time.Time
	ObjDomain domain.Domain
	ObjName   name.Name
	Action    string
	Message   string
	Data      json.RawMessage
	ID        uuid.UUID
	ObjID     uuid.UUID
	ActorID   uuid.UUID
}

// NewAudit represents the information needed to create a new audit record.
type NewAudit struct {
	Data      any
	ObjDomain domain.Domain
	ObjName   name.Name
	Action    string
	Message   string
	ObjID     uuid.UUID
	ActorID   uuid.UUID
}
