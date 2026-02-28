package homebus

import "github.com/garnizeh/fingo/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByHomeID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByHomeID = "home_id"
	OrderByType   = "type"
	OrderByUserID = "user_id"
)
