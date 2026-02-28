package homeapp

import (
	"github.com/garnizeh/fingo/business/domain/homebus"
)

var orderByFields = map[string]string{
	"home_id": homebus.OrderByHomeID,
	"type":    homebus.OrderByType,
	"user_id": homebus.OrderByUserID,
}
