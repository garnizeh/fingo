package auditdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/types/domain"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type audit struct {
	Timestamp time.Time          `db:"timestamp"`
	ObjDomain string             `db:"obj_domain"`
	ObjName   string             `db:"obj_name"`
	Action    string             `db:"action"`
	Message   string             `db:"message"`
	Data      types.NullJSONText `db:"data"`
	ID        uuid.UUID          `db:"id"`
	ObjID     uuid.UUID          `db:"obj_id"`
	ActorID   uuid.UUID          `db:"actor_id"`
}

func toDBAudit(bus *auditbus.Audit) (audit, error) {
	db := audit{
		ID:        bus.ID,
		ObjID:     bus.ObjID,
		ObjDomain: bus.ObjDomain.String(),
		ObjName:   bus.ObjName.String(),
		ActorID:   bus.ActorID,
		Action:    bus.Action,
		Data:      types.NullJSONText{JSONText: []byte(bus.Data), Valid: true},
		Message:   bus.Message,
		Timestamp: bus.Timestamp.UTC(),
	}

	return db, nil
}

func toBusAudit(db *audit) (auditbus.Audit, error) {
	domain, err := domain.Parse(db.ObjDomain)
	if err != nil {
		return auditbus.Audit{}, fmt.Errorf("parse domain: %w", err)
	}

	name, err := name.Parse(db.ObjName)
	if err != nil {
		return auditbus.Audit{}, fmt.Errorf("parse name: %w", err)
	}

	bus := auditbus.Audit{
		ID:        db.ID,
		ObjID:     db.ObjID,
		ObjDomain: domain,
		ObjName:   name,
		ActorID:   db.ActorID,
		Action:    db.Action,
		Data:      json.RawMessage(db.Data.JSONText),
		Message:   db.Message,
		Timestamp: db.Timestamp.Local(),
	}

	return bus, nil
}

func toBusAudits(dbs []audit) ([]auditbus.Audit, error) {
	audits := make([]auditbus.Audit, len(dbs))

	for i := range dbs {
		a, err := toBusAudit(&dbs[i])
		if err != nil {
			return nil, err
		}

		audits[i] = a
	}

	return audits, nil
}
