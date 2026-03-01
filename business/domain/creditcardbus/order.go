package creditcardbus

import "github.com/garnizeh/fingo/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID          = "credit_card_id"
	OrderByUserID      = "user_id"
	OrderByName        = "name"
	OrderByEnabled     = "enabled"
	OrderByDateCreated = "date_created"
)
