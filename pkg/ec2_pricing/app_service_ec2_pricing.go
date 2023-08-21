package ec2_pricing

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"strconv"
	"worker-ms-aws-costs/utils/logger"
)

const (
	ec2ApiUrl = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
	batchSize = 5000
)

type PortsServerEc2Pricing interface {
	GetPricing() error
}

type service struct {
	repository ServicesEc2PricingRepository
}

func NewEc2PricingService(repository ServicesEc2PricingRepository) PortsServerEc2Pricing {
	return &service{repository: repository}
}

func (s service) GetPricing() error {

	color.Blue("Iniciando worker para obtener datos de EC2 - EBS Pricing...")

	// Obtener datos de precios de AWS
	ec2Products, ebsProducts, err := s.fetchPricingData()
	if err != nil {
		logger.Error.Println("Error al obtener datos de precios de EC2:", err)
		return err
	}

	// Insertar datos en la base de datos
	err = s.repository.insertBulkEBSProducts(ebsProducts)
	if err != nil {
		return err
	}
	err = s.repository.insertBulkEC2Products(ec2Products)
	if err != nil {
		return err
	}

	color.Green("Información almacenada exitosamente.")

	err = s.repository.updateValorProductoProveedorNube()
	if err != nil {
		logger.Error.Println("Error al obtener actualizar valores:", err)
		return err
	}
	color.Green("Información actualizada exitosamente.")
	return nil
}

func (s service) fetchPricingData() ([]ProductAttributes, []ProductAttributes, error) {

	resp, err := http.Get(ec2ApiUrl)
	if err != nil {
		logger.Error.Println("Error fetch:", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to fetch EC2 pricing data: %s", resp.Status)
	}

	var pricingResponse PricingResponse
	err = json.NewDecoder(resp.Body).Decode(&pricingResponse)
	if err != nil {
		logger.Error.Println("Error json Decode:", err)
		return nil, nil, err
	}

	var productsEC2 []ProductAttributes
	var productsEBS []ProductAttributes
	for sku, product := range pricingResponse.Products {
		if product.ProductFamily == "Storage" || product.ProductFamily == "Compute Instance" {
			onDemandDetails, exists := pricingResponse.Terms.OnDemand[sku]
			if !exists {
				continue
			}

			for _, details := range onDemandDetails {
				for _, priceDimension := range details.PriceDimensions {
					product.Attributes.Sku = product.Sku
					product.Attributes.ProductFamily = product.ProductFamily
					price, _ := strconv.ParseFloat(priceDimension.PricePerUnit.USD, 64)
					product.Attributes.TermType = "OnDemand"
					product.Attributes.PricePerHour = price
					product.Attributes.BillingItem = priceDimension.Unit
					if product.ProductFamily == "Storage" {
						productsEBS = append(productsEBS, product.Attributes)
					}
					if product.ProductFamily == "Compute Instance" {
						productsEC2 = append(productsEC2, product.Attributes)
					}
				}
			}
		}
	}

	return productsEC2, productsEBS, nil
}
