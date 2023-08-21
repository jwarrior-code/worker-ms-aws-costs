package rds_pricing

import (
	"github.com/jmoiron/sqlx"
	"worker-ms-aws-costs/utils/logger"
)

const (
	Mysql = "mysql"
)

type ServicesRDSPricingRepository interface {
	insertBulkProducts(products []ProductAttributes) error
	insertBatch(baseQuery string, products []ProductAttributes) error
}

func FactoryStorage(db *sqlx.DB) ServicesRDSPricingRepository {
	var s ServicesRDSPricingRepository
	engine := db.DriverName()
	switch engine {
	case Mysql:
		return newRDSPricingMysqlRepository(db)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
