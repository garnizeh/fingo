package userbus

import "github.com/garnizeh/fingo/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByUserID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByUserID  = "user_id"
	OrderByName    = "name"
	OrderByEmail   = "email"
	OrderByRoles   = "roles"
	OrderByEnabled = "enabled"
)
