package dbx

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"worker-ms-aws-costs/utils/env"
	"worker-ms-aws-costs/utils/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	dbx      *sqlx.DB
	once     sync.Once
	DBEngine string
)

func init() {
	once.Do(func() {
		setConnection()
	})
}

func setConnection() {
	var err error
	c := env.NewConfiguration()
	DBEngine = c.DB.Engine

	// Check the connection
	dbx, err = sqlx.Open(DBEngine, connectionString("data"))
	if err != nil {
		logger.Error.Printf("no se puede conectar a la base de datos: %v", err)
		panic(err)
	}
	err = dbx.Ping()
	if err != nil {
		logger.Error.Printf("couldn't connect to database: %v", err)
		panic(err)
	}
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(2 * time.Minute)
	dbx.SetMaxOpenConns(95)
}

func GetConnection() *sqlx.DB {
	if dbx == nil {
		setConnection()
	}
	return dbx
}

func connectionString(t string) string {
	c := env.NewConfiguration()

	var host, database, username, password string
	var port int
	switch t {
	case "data":
		host = c.DB.Server
		database = c.DB.Name
		username = c.DB.User
		password = c.DB.Password
		port = c.DB.Port
	default:
		logger.Error.Print("El tipo de conexión no correspondea data/logs")
		return ""
	}
	switch strings.ToLower(DBEngine) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	}
	logger.Error.Print("el motor de bases de datos solicitado no está configurado aún")

	return ""
}
