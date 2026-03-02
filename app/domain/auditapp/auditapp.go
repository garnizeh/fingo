// Package auditapp maintains the app layer api for the audit domain.
package auditapp

import (
	"context"
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/app/sdk/query"
	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/foundation/web"
)

type app struct {
	auditBus auditbus.ExtBusiness
}

func newApp(auditBus auditbus.ExtBusiness) *app {
	return &app{
		auditBus: auditBus,
	}
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := qp.parseFilter()
	if err != nil {
		if enc, ok := err.(web.Encoder); ok {
			return enc
		}
		return errs.New(errs.Internal, err)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, userbus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	adts, err := a.auditBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	total, err := a.auditBus.Count(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppAudits(adts), total, page)
}
