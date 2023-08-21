package rds_pricing

type PricingResponse struct {
	Products map[string]Product `json:"products"`
	Terms    Terms              `json:"terms"`
}

type Product struct {
	Sku           string            `json:"sku" db:"sku"`
	ProductFamily string            `json:"productFamily" db:"product_family"`
	Attributes    ProductAttributes `json:"attributes"`
}

type ProductAttributes struct {
	Sku              string  `json:"sku" db:"sku"`
	ProductFamily    string  `json:"productFamily" db:"product_family"`
	Servicecode      string  `json:"servicecode" db:"servicecode"`
	Location         string  `json:"location" db:"location"`
	DatabaseEngine   string  `json:"databaseEngine,omitempty" db:"database_engine"`
	DeploymentOption string  `json:"deploymentOption,omitempty" db:"deployment_option"`
	LicenseModel     string  `json:"licenseModel,omitempty" db:"license_model"`
	Vcpu             string  `json:"vcpu,omitempty" db:"vcpu"`
	Memory           string  `json:"memory,omitempty" db:"memory"`
	InstanceClass    string  `json:"instanceType,omitempty" db:"instance_class"`
	TermType         string  `db:"term_type" db:"term_type"`
	PricePerHour     float64 `db:"price_per_hour" db:"price_per_hour"`
	BillingItem      string  `db:"billing_item" db:"billing_item"`
}

type Terms struct {
	OnDemand map[string]map[string]OnDemandDetails `json:"OnDemand"`
}

type OnDemandDetails struct {
	PriceDimensions map[string]PriceDimension `json:"priceDimensions"`
}

type PriceDimension struct {
	PricePerUnit PricePerUnit `json:"pricePerUnit"`
	Unit         string       `json:"unit"`
}

type PricePerUnit struct {
	USD string `json:"USD"`
}
