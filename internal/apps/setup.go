package core

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/mrKrabsmr/commerce-edu-api/internal/configs"
	dbConnection "github.com/mrKrabsmr/commerce-edu-api/pkg/db_connection"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	config     *configs.Config
	configOnce = new(sync.Once)

	database     *sqlx.DB
	databaseOnce = new(sync.Once)

	logger     *logrus.Logger
	loggerOnce = new(sync.Once)

	v     *validator.Validate
	vOnce = new(sync.Once)
)

func GetConfig() *configs.Config {
	configOnce.Do(func() {
		config = &configs.Config{
			LogLevel:  os.Getenv("LOG_LEVEL"),
			Address:   os.Getenv("ADDRESS"),
			DBDialect: os.Getenv("DBDialect"),
			DBAddress: os.Getenv("DBAddress"),
		}
	})

	return config
}

func GetDB() *sqlx.DB {
	databaseOnce.Do(func() {
		c := GetConfig()
		conn := dbConnection.NewPGConnection(c)
		db, err := conn.PostgreSQLConnection()
		if err != nil {
			panic(err)
		}

		database = db
	})

	return database
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(func() {
		c := GetConfig()
		logLevel, err := logrus.ParseLevel(c.LogLevel)
		if err != nil {
			panic(err)
		}

		logger = logrus.New()
		logger.SetLevel(logLevel)
	})

	return logger
}

func GetValidator() *validator.Validate {
	vOnce.Do(func() {
		v = validator.New()
	})

	return v
}
