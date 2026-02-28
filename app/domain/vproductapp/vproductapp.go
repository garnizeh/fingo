// Package vproductapp maintains the app layer api for the vproduct domain.
package vproductapp

import (
	"context"
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/app/sdk/query"
	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/foundation/web"
)

type app struct {
	vproductBus vproductbus.ExtBusiness
}

func newApp(vproductBus vproductbus.ExtBusiness) *app {
	return &app{
		vproductBus: vproductBus,
	}
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, vproductbus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	prds, err := a.vproductBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	total, err := a.vproductBus.Count(ctx, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppProducts(prds), total, page)
}
