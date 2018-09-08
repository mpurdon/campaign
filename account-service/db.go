package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"math"
	"time"
)

var DATABASE = map[string]string{
	"host":     "localhost",
	"port":     "5432",
	"name":     "fc",
	"user":     "fc",
	"password": "fcpass",
}

func getDbConfig() map[string]string {

	Logger.Info("Getting database configuration.")

	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		dbHost = DATABASE["host"]
		Logger.Infof("DB_HOST not found in environment, using default value '%s'", dbHost)
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		dbPort = DATABASE["port"]
		Logger.Infof("DB_PORT not found in environment, using default value '%s'", dbPort)
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = DATABASE["name"]
		Logger.Infof("DB_NAME not found in environment, using default value '%s'", dbName)
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		dbUser = DATABASE["user"]
		Logger.Infof("DB_USER not found in environment, using default value '%s'", dbUser)
	}

	dbPassword, ok := os.LookupEnv("DB_PASS")
	if !ok {
		dbPassword = DATABASE["password"]
		Logger.Infof("DB_PASS not found in environment, using default value '%s'", dbPassword)
	}

	config := make(map[string]string)
	config["host"] = dbHost
	config["port"] = dbPort
	config["name"] = dbName
	config["user"] = dbUser
	config["password"] = dbPassword

	return config
}

/**
 * Create a GORM database connection
 */
func createConnection() *gorm.DB {

	config := getDbConfig()
	//psqlInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config["user"], config["password"], config["host"], config["name"])
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", config["host"], config["port"], config["name"], config["user"], config["password"])

	// @Security remove this log line.
	Logger.Infof("Connecting to the database with ", psqlInfo)

	var (
		attempt    = 1
		db         *gorm.DB
		err        error
		retryDelay time.Duration = 0
	)

	for attempt <= 4 {
		Logger.Infof("Making connection attempt: #%d", attempt)

		db, err = gorm.Open(
			"postgres",
			psqlInfo,
		)

		if err == nil {
			break
		}

		retryDelay = time.Duration(math.Exp2(float64(attempt)))
		Logger.Warnf("Error opening connection to database, retrying after %ds: %v", retryDelay, err)
		//time.Sleep(retryDelay)
		time.Sleep(retryDelay * time.Second)

		attempt++
	}

	if attempt == 3 && err != nil {
		Logger.Fatalf("Error opening connection to database: ", err)
	}

	err = db.DB().Ping()
	if err != nil {
		Logger.Fatalf("Error pinging database: ", err)
	}

	Logger.Infof("Connected to database.")
	return db
}
