package ec2_pricing

import (
	"github.com/jmoiron/sqlx"
	"worker-ms-aws-costs/utils/logger"
)

const (
	Mysql = "mysql"
)

type ServicesEc2PricingRepository interface {
	insertBulkEC2Products(products []ProductAttributes) error
	insertEC2Batch(baseQuery string, batch []ProductAttributes) error
	insertBulkEBSProducts(products []ProductAttributes) error
	insertEBSBatch(baseQuery string, batch []ProductAttributes) error
	updateValorProductoProveedorNube() error
}

func FactoryStorage(db *sqlx.DB) ServicesEc2PricingRepository {
	var s ServicesEc2PricingRepository
	engine := db.DriverName()
	switch engine {
	case Mysql:
		return newEc2PricingMysqlRepository(db)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
