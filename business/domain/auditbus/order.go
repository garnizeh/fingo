package auditbus

import "github.com/garnizeh/fingo/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByObjID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByObjID     = "obj_id"
	OrderByObjDomain = "obj_domain"
	OrderByObjName   = "obj_name"
	OrderByActorID   = "actor_id"
	OrderByAction    = "action"
)
