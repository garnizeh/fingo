package dbtest

import (
	"time"

	"github.com/garnizeh/fingo/business/domain/auditbus"
	"github.com/garnizeh/fingo/business/domain/auditbus/extensions/auditotel"
	"github.com/garnizeh/fingo/business/domain/auditbus/stores/auditdb"
	"github.com/garnizeh/fingo/business/domain/homebus"
	"github.com/garnizeh/fingo/business/domain/homebus/extensions/homeotel"
	"github.com/garnizeh/fingo/business/domain/homebus/stores/homedb"
	"github.com/garnizeh/fingo/business/domain/productbus"
	"github.com/garnizeh/fingo/business/domain/productbus/extensions/productotel"
	"github.com/garnizeh/fingo/business/domain/productbus/stores/productdb"
	"github.com/garnizeh/fingo/business/domain/userbus"
	"github.com/garnizeh/fingo/business/domain/userbus/extensions/useraudit"
	"github.com/garnizeh/fingo/business/domain/userbus/extensions/userotel"
	"github.com/garnizeh/fingo/business/domain/userbus/stores/usercache"
	"github.com/garnizeh/fingo/business/domain/userbus/stores/userdb"
	"github.com/garnizeh/fingo/business/domain/vproductbus"
	"github.com/garnizeh/fingo/business/domain/vproductbus/extensions/vproductotel"
	"github.com/garnizeh/fingo/business/domain/vproductbus/stores/vproductdb"
	"github.com/garnizeh/fingo/business/sdk/delegate"
	"github.com/garnizeh/fingo/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Audit    auditbus.ExtBusiness
	Home     homebus.ExtBusiness
	Product  productbus.ExtBusiness
	User     userbus.ExtBusiness
	VProduct vproductbus.ExtBusiness
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)

	auditOtelExt := auditotel.NewExtension()
	auditStorage := auditdb.NewStore(log, db)
	auditBus := auditbus.NewBusiness(log, auditStorage, auditOtelExt)

	userOtelExt := userotel.NewExtension()
	userAuditExt := useraudit.NewExtension(auditBus)
	userStorage := usercache.NewStore(log, userdb.NewStore(log, db), time.Hour)
	userBus := userbus.NewBusiness(log, delegate, userStorage, userOtelExt, userAuditExt)

	productOtelExt := productotel.NewExtension()
	productStorage := productdb.NewStore(log, db)
	productBus := productbus.NewBusiness(log, userBus, delegate, productStorage, productOtelExt)

	homeOtelExt := homeotel.NewExtension()
	homeStorage := homedb.NewStore(log, db)
	homeBus := homebus.NewBusiness(log, userBus, delegate, homeStorage, homeOtelExt)

	vproductOtelExt := vproductotel.NewExtension()
	vproductStorage := vproductdb.NewStore(log, db)
	vproductBus := vproductbus.NewBusiness(vproductStorage, vproductOtelExt)

	return BusDomain{
		Delegate: delegate,
		Audit:    auditBus,
		Home:     homeBus,
		Product:  productBus,
		User:     userBus,
		VProduct: vproductBus,
	}
}
