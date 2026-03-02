package audit_test

import (
	"time"

	"github.com/garnizeh/fingo/app/domain/auditapp"
	"github.com/garnizeh/fingo/business/domain/auditbus"
)

func toAppAudit(bus *auditbus.Audit) auditapp.Audit {
	return auditapp.Audit{
		ID:        bus.ID.String(),
		ObjID:     bus.ObjID.String(),
		ObjDomain: bus.ObjDomain.String(),
		ObjName:   bus.ObjName.String(),
		ActorID:   bus.ActorID.String(),
		Action:    bus.Action,
		Data:      string(bus.Data),
		Message:   bus.Message,
		Timestamp: bus.Timestamp.Format(time.RFC3339),
	}
}

func toAppAudits(audits []auditbus.Audit) []auditapp.Audit {
	app := make([]auditapp.Audit, len(audits))
	for i := range audits {
		app[i] = toAppAudit(&audits[i])
	}

	return app
}
