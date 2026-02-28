package userapp

import (
	"github.com/garnizeh/fingo/business/domain/userbus"
)

var orderByFields = map[string]string{
	"user_id": userbus.OrderByUserID,
	"name":    userbus.OrderByName,
	"email":   userbus.OrderByEmail,
	"roles":   userbus.OrderByRoles,
	"enabled": userbus.OrderByEnabled,
}
