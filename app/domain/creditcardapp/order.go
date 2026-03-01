package creditcardapp

import (
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
)

var orderByFields = map[string]string{
	"credit_card_id": creditcardbus.OrderByID,
	"user_id":        creditcardbus.OrderByUserID,
	"name":           creditcardbus.OrderByName,
	"enabled":        creditcardbus.OrderByEnabled,
	"date_created":   creditcardbus.OrderByDateCreated,
}
