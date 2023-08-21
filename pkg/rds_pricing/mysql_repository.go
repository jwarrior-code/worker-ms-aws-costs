package rds_pricing

import (
	"github.com/jmoiron/sqlx"
	"worker-ms-aws-costs/utils/logger"
)

// sqlServer estructura de conexión a la BD de mssql
type mysql struct {
	DB *sqlx.DB
}

func newRDSPricingMysqlRepository(db *sqlx.DB) *mysql {
	return &mysql{
		DB: db,
	}
}

func (s mysql) insertBulkProducts(products []ProductAttributes) error {
	baseQuery := `INSERT INTO aws_rds_products (sku, product_family, servicecode, location, database_engine, deployment_option, license_model, vcpu, memory, instance_class, term_type, price_per_hour, billing_item) VALUES `

	for start := 0; start < len(products); start += batchSize {
		end := start + batchSize
		if end > len(products) {
			end = len(products)
		}

		err := s.insertBatch(baseQuery, products[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s mysql) insertBatch(baseQuery string, products []ProductAttributes) error {

	values := []interface{}{}
	insertPlaceholders := ""

	for _, product := range products {
		insertPlaceholders += `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`
		values = append(values, product.Sku, product.ProductFamily, product.Servicecode, product.Location, product.DatabaseEngine, product.DeploymentOption, product.LicenseModel, product.Vcpu, product.Memory, product.InstanceClass, product.TermType, product.PricePerHour, product.BillingItem)
	}

	insertPlaceholders = insertPlaceholders[:len(insertPlaceholders)-1] // Eliminar la última coma

	query := baseQuery + insertPlaceholders + ` ON DUPLICATE KEY UPDATE
                product_family = VALUES(product_family), 
                servicecode = VALUES(servicecode), 
                location = VALUES(location), 
                database_engine = VALUES(database_engine), 
                deployment_option = VALUES(deployment_option), 
                license_model = VALUES(license_model), 
                vcpu = VALUES(vcpu), 
                memory = VALUES(memory), 
                instance_class = VALUES(instance_class),
                term_type = VALUES(term_type), 
                price_per_hour = VALUES(price_per_hour),
                billing_item = VALUES(billing_item);`

	// Convertir los marcadores de posición
	query = s.DB.Rebind(query)

	_, err := s.DB.Exec(query, values...)
	if err != nil {
		logger.Error.Println("Error al insertar productos en bulk:", err)
		return err
	}
	return nil
}
