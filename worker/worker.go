package worker

import (
	"worker-ms-aws-costs/pkg"
	"worker-ms-aws-costs/utils/logger"
)

type Worker struct {
	srv *pkg.Server
}

func NewWorker(srv *pkg.Server) IWorker {
	return &Worker{srv: srv}
}

func (w Worker) Execute() {

	err := w.srv.SrvRDSPricing.GetPricing()
	if err != nil {
		logger.Error.Println("no se puede actualizar costos de rds", err)
		return
	}
	err = w.srv.SrvEc2Pricing.GetPricing()
	if err != nil {
		logger.Error.Println("no se puede actualizar costos de ec2", err)
		return
	}

}
