package creditcarddb

import (
	"fmt"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/sdk/order"
)

var orderByFields = map[string]string{
	creditcardbus.OrderByID:          "credit_card_id",
	creditcardbus.OrderByUserID:      "user_id",
	creditcardbus.OrderByName:        "name",
	creditcardbus.OrderByEnabled:     "enabled",
	creditcardbus.OrderByDateCreated: "date_created",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
