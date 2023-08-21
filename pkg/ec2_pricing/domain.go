package ec2_pricing

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
	Sku                string  `json:"sku" db:"sku"`
	ProductFamily      string  `json:"productFamily" db:"product_family"`
	Servicecode        string  `json:"servicecode" db:"servicecode"`
	Location           string  `json:"location" db:"location"`
	VolumeType         string  `json:"volumeType,omitempty" db:"volume_type"`
	InstanceType       string  `json:"instanceType,omitempty" db:"instance_type"`
	Storage            string  `json:"Storage,omitempty" db:"storage"`
	PricePerGbMonth    float64 `json:"PricePerGbMonth,omitempty" db:"price_per_gb_month"`
	Vcpu               string  `json:"vcpu,omitempty" db:"vcpu"`
	Memory             string  `json:"memory,omitempty" db:"memory"`
	NetworkPerformance string  `json:"networkPerformance,omitempty" db:"network_performance"`
	TermType           string  `db:"term_type" db:"term_type"`
	PricePerHour       float64 `db:"price_per_hour" db:"price_per_hour"`
	BillingItem        string  `db:"billing_item" db:"billing_item"`
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
