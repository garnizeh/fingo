package creditcardapp

import (
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/google/uuid"
)

type queryParams struct {
	Page    string
	Rows    string
	OrderBy string
	ID      string
	UserID  string
	Name    string
	Enabled string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()

	filter := queryParams{
		Page:    values.Get("page"),
		Rows:    values.Get("rows"),
		OrderBy: values.Get("orderBy"),
		ID:      values.Get("credit_card_id"),
		UserID:  values.Get("user_id"),
		Name:    values.Get("name"),
		Enabled: values.Get("enabled"),
	}

	return filter
}

func parseFilter(qp queryParams) (creditcardbus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter creditcardbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("credit_card_id", err)
		}
	}

	if qp.UserID != "" {
		userID, err := uuid.Parse(qp.UserID)
		switch err {
		case nil:
			filter.UserID = &userID
		default:
			fieldErrors.Add("user_id", err)
		}
	}

	if qp.Name != "" {
		name, err := name.Parse(qp.Name)
		switch err {
		case nil:
			filter.Name = &name
		default:
			fieldErrors.Add("name", err)
		}
	}

	if qp.Enabled != "" {
		switch qp.Enabled {
		case "true":
			b := true
			filter.Enabled = &b
		case "false":
			b := false
			filter.Enabled = &b
		default:
			fieldErrors.Add("enabled", errs.Errorf(errs.InvalidArgument, "invalid boolean format"))
		}
	}

	if fieldErrors != nil {
		return creditcardbus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
