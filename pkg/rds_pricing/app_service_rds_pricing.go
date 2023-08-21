package rds_pricing

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strconv"
	"worker-ms-aws-costs/utils/logger"
)

type PortsServerRDSPricing interface {
	GetPricing() error
}

type service struct {
	repository ServicesRDSPricingRepository
}

func NewRDSPricingService(repository ServicesRDSPricingRepository) PortsServerRDSPricing {
	return &service{repository: repository}
}

const (
	rdsApiUrl = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRDS/current/index.json"
	batchSize = 5000 // Número de registros por lote
)

func (s service) GetPricing() error {

	color.Blue("Iniciando worker para obtener datos de RDS Pricing...")

	response, err := http.Get(rdsApiUrl)
	if err != nil {
		logger.Error.Println("Error al obtener los datos:", err)
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error.Println("Error al leer la respuesta:", err)
		return err
	}

	var pricingData PricingResponse
	err = json.Unmarshal(body, &pricingData)
	if err != nil {
		logger.Error.Println("Error al deserializar JSON:", err)
		return err
	}

	var productsToInsert []ProductAttributes
	for sku, product := range pricingData.Products {
		if product.ProductFamily == "Database Instance" {
			onDemandDetails, exists := pricingData.Terms.OnDemand[sku]
			if !exists {
				continue
			}

			for _, details := range onDemandDetails {
				for _, priceDimension := range details.PriceDimensions {
					price, _ := strconv.ParseFloat(priceDimension.PricePerUnit.USD, 64)
					product.Attributes.Sku = product.Sku
					product.Attributes.ProductFamily = product.ProductFamily
					product.Attributes.TermType = "OnDemand"
					product.Attributes.PricePerHour = price
					product.Attributes.BillingItem = priceDimension.Unit
					productsToInsert = append(productsToInsert, product.Attributes)

				}
			}
		}
	}
	err = s.repository.insertBulkProducts(productsToInsert)
	if err != nil {
		return err
	}
	fmt.Println("Información RDS almacenada exitosamente.")
	return nil
}
