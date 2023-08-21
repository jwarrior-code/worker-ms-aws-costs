package pkg

import (
	"github.com/jmoiron/sqlx"
	"worker-ms-aws-costs/pkg/ec2_pricing"
	"worker-ms-aws-costs/pkg/rds_pricing"
)

type Server struct {
	SrvEc2Pricing ec2_pricing.PortsServerEc2Pricing
	SrvRDSPricing rds_pricing.PortsServerRDSPricing
}

func NewServerWorkerAwsCosts(db *sqlx.DB) *Server {

	return &Server{
		SrvEc2Pricing: ec2_pricing.NewEc2PricingService(
			ec2_pricing.FactoryStorage(db)),
		SrvRDSPricing: rds_pricing.NewRDSPricingService(
			rds_pricing.FactoryStorage(db)),
	}
}
