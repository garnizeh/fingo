// Package creditcardapp maintains the app layer api for the credit card domain.
package creditcardapp

import (
	"context"
	"net/http"

	"github.com/garnizeh/fingo/app/sdk/errs"
	"github.com/garnizeh/fingo/app/sdk/mid"
	"github.com/garnizeh/fingo/app/sdk/query"
	"github.com/garnizeh/fingo/business/domain/creditcardbus"
	"github.com/garnizeh/fingo/business/sdk/order"
	"github.com/garnizeh/fingo/business/sdk/page"
	"github.com/garnizeh/fingo/foundation/web"
	"github.com/google/uuid"
)

type app struct {
	creditCardBus creditcardbus.ExtBusiness
}

func newApp(creditCardBus creditcardbus.ExtBusiness) *app {
	return &app{
		creditCardBus: creditCardBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var ncc NewCreditCard
	if err := web.Decode(r, &ncc); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	busNCC, err := toBusNewCreditCard(ctx, ncc)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	cc, err := a.creditCardBus.Create(ctx, actorID, busNCC)
	if err != nil {
		return errs.Errorf(errs.Internal, "create: cc[%+v]: %s", busNCC, err)
	}

	return toAppCreditCard(&cc)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var ucc UpdateCreditCard
	if err := web.Decode(r, &ucc); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	busUCC, err := toBusUpdateCreditCard(ucc)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ccID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	cc, err := a.creditCardBus.QueryByID(ctx, actorID, ccID)
	if err != nil {
		return errs.Errorf(errs.Internal, "querybyid: ccID[%s]: %s", ccID, err)
	}

	updated, err := a.creditCardBus.Update(ctx, actorID, &cc, busUCC)
	if err != nil {
		return errs.Errorf(errs.Internal, "update: ccID[%s] ucc[%+v]: %s", ccID, ucc, err)
	}

	return toAppCreditCard(&updated)
}

func (a *app) delete(ctx context.Context, r *http.Request) web.Encoder {
	ccID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	cc, err := a.creditCardBus.QueryByID(ctx, actorID, ccID)
	if err != nil {
		return errs.Errorf(errs.Internal, "querybyid: ccID[%s]: %s", ccID, err)
	}

	if err := a.creditCardBus.Delete(ctx, actorID, &cc); err != nil {
		return errs.Errorf(errs.Internal, "delete: ccID[%s]: %s", ccID, err)
	}

	return nil
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseQueryParams(r)

	filter, err := qp.parseFilter()
	if err != nil {
		if enc, ok := err.(web.Encoder); ok {
			return enc
		}
		return errs.New(errs.Internal, err)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, creditcardbus.DefaultOrderBy)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	pg, err := page.Parse(qp.Page, qp.Rows)

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	ccs, err := a.creditCardBus.Query(ctx, actorID, filter, orderBy, pg)
	if err != nil {
		return errs.Errorf(errs.Internal, "query: %s", err)
	}

	total, err := a.creditCardBus.Count(ctx, actorID, filter)
	if err != nil {
		return errs.Errorf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppCreditCards(ccs), total, pg)
}

func (a *app) queryByID(ctx context.Context, r *http.Request) web.Encoder {
	ccID, err := uuid.Parse(web.Param(r, "id"))
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	actorID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	cc, err := a.creditCardBus.QueryByID(ctx, actorID, ccID)
	if err != nil {
		return errs.Errorf(errs.Internal, "querybyid: ccID[%s]: %s", ccID, err)
	}

	return toAppCreditCard(&cc)
}
