package creditcarddb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/garnizeh/fingo/business/domain/creditcardbus"
)

func applyFilter(filter creditcardbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["credit_card_id"] = *filter.ID
		wc = append(wc, "credit_card_id = :credit_card_id")
	}

	if filter.UserID != nil {
		data["user_id"] = *filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
