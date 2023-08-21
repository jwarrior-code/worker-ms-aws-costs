package ec2_pricing

import (
	"github.com/jmoiron/sqlx"
	"log"
	"worker-ms-aws-costs/utils/logger"
)

// sqlServer estructura de conexión a la BD de mssql
type mysql struct {
	DB *sqlx.DB
}

func newEc2PricingMysqlRepository(db *sqlx.DB) *mysql {
	return &mysql{
		DB: db,
	}
}

func (s mysql) insertBulkEC2Products(products []ProductAttributes) error {
	baseQuery := `INSERT INTO aws_ec2_products (sku, product_family, servicecode, location, instance_type, vcpu, memory, storage, network_performance, term_type, price_per_hour, billing_item) VALUES `

	for start := 0; start < len(products); start += batchSize {
		end := start + batchSize
		if end > len(products) {
			end = len(products)
		}
		err := s.insertEC2Batch(baseQuery, products[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s mysql) insertEC2Batch(baseQuery string, batch []ProductAttributes) error {
	values := []interface{}{}
	insertPlaceholders := ""

	for _, product := range batch {
		insertPlaceholders += `(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`
		values = append(values, product.Sku, product.ProductFamily, product.Servicecode, product.Location, product.InstanceType,
			product.Vcpu, product.Memory, product.Storage, product.NetworkPerformance, product.TermType, product.PricePerHour,
			product.BillingItem)
	}

	insertPlaceholders = insertPlaceholders[:len(insertPlaceholders)-1] // Eliminar la última coma

	query := baseQuery + insertPlaceholders + ` ON DUPLICATE KEY UPDATE 
					product_family = VALUES(product_family), 
					servicecode = VALUES(servicecode), 
					location = VALUES(location), 
					volume_type = VALUES(volume_type), 
					instance_type = VALUES(instance_type), 
					vcpu = VALUES(vcpu), 
					memory = VALUES(memory), 
					storage = VALUES(storage),
					network_performance = VALUES(network_performance),
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

func (s mysql) insertBulkEBSProducts(products []ProductAttributes) error {
	baseQuery := `INSERT INTO aws_storage_products (sku, product_family, servicecode, location, volume_type, price_per_gb_month, term_type, billing_item, price_per_hour) VALUES `

	for start := 0; start < len(products); start += batchSize {
		end := start + batchSize
		if end > len(products) {
			end = len(products)
		}
		err := s.insertEBSBatch(baseQuery, products[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s mysql) insertEBSBatch(baseQuery string, batch []ProductAttributes) error {
	values := []interface{}{}
	insertPlaceholders := ""

	for _, product := range batch {
		insertPlaceholders += `(?, ?, ?, ?, ?, ?, ?, ?, ?),`
		values = append(values, product.Sku, product.ProductFamily, product.Servicecode, product.Location, product.VolumeType, product.PricePerGbMonth, product.TermType, product.BillingItem, product.PricePerHour)
	}

	insertPlaceholders = insertPlaceholders[:len(insertPlaceholders)-1] // Eliminar la última coma

	query := baseQuery + insertPlaceholders + ` ON DUPLICATE KEY UPDATE  
					product_family = VALUES(product_family), 
					servicecode = VALUES(servicecode), 
					location = VALUES(location), 
					volume_type = VALUES(volume_type),
					price_per_gb_month = VALUES(price_per_gb_month),
					term_type = VALUES(term_type),
					billing_item = VALUES(billing_item),
					price_per_hour = VALUES(price_per_hour) ;`

	// Convertir los marcadores de posición
	query = s.DB.Rebind(query)

	_, err := s.DB.Exec(query, values...)
	if err != nil {
		logger.Error.Println("Error al insertar productos en bulk:", err)
		return err
	}
	return nil
}

func (s mysql) updateValorProductoProveedorNube() error {

	query := `UPDATE producto_proveedor_nube AS ppn
		LEFT JOIN aws_rds_products AS arp ON ppn.SKU = arp.sku
		LEFT JOIN aws_ec2_products AS aep ON ppn.SKU = aep.sku
		LEFT JOIN aws_storage_products AS asp ON ppn.SKU = asp.sku
		SET ppn.valor = COALESCE(arp.price_per_hour, aep.price_per_hour, asp.price_per_hour)
		WHERE ppn.proveedores_nube_idProveedor = 1
		AND (arp.sku IS NOT NULL OR aep.sku IS NOT NULL OR asp.sku IS NOT NULL);
	`

	// Ejecuta el query
	result, err := s.DB.Exec(query)
	if err != nil {
		logger.Error.Println("Error al ejecutar la actualización de valor aws:", err)
		return err
	}

	// Puedes también obtener la cantidad de filas afectadas si lo necesitas
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error.Println("Error al obtener filas afectadas:", err)
		return err
	}

	log.Printf("Se actualizó %d fila(s)", rowsAffected)
	return nil
}
