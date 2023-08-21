package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"worker-ms-aws-costs/pkg"
	"worker-ms-aws-costs/utils/dbx"
	"worker-ms-aws-costs/utils/logger"
	"worker-ms-aws-costs/worker"
)

func main() {
	fmt.Println("test")
	color.Blue("worker-ms-aws-costs v1.0.0")
	db := dbx.GetConnection()
	srv := pkg.NewServerWorkerAwsCosts(db)
	wk := worker.NewWorker(srv)
	wk.Execute()

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logger.Error.Println("error desconectando de la base de datos", err)
		}
	}(db)
}
