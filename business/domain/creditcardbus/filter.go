package creditcardbus

import (
	"fmt"
	"time"

	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID               *uuid.UUID
	UserID           *uuid.UUID
	Name             *name.Name
	Enabled          *bool
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}

// Validate checks the data within the filter.
func (f *QueryFilter) Validate() error {
	if f.StartCreatedDate != nil && f.EndCreatedDate != nil {
		if f.StartCreatedDate.After(*f.EndCreatedDate) {
			return fmt.Errorf("start_created_date cannot be after end_created_date")
		}
	}
	return nil
}
