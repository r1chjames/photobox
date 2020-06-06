package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)


func initDbConnection(appConfig AppConfig) *sql.DB {
	db, err := sql.Open("mysql", appConfig.DbUrl)
	if err != nil {
		log.Fatal(fmt.Sprintf("An error occurred connecting to MySQL database: %s", err))
	}

	return db
}

func executeDbInsert(appConfig AppConfig, query string) error {
	db := initDbConnection(appConfig)
	defer db.Close()
	results, err := db.Query(query)
	if err != nil {
		log.Fatal(fmt.Sprintf("An error occurred executing insert: %s", err))
		return err
	}
	results.Close()
	return nil
}

func executeDbQuery(appConfig AppConfig, query string) *sql.Rows {
	db := initDbConnection(appConfig)
	defer db.Close()
	results, err := db.Query(query)
	if err != nil {
		log.Fatal(fmt.Sprintf("An error occurred executing query: %s", err))
	}
	return results
}

func executeTransformDbQuery(appConfig AppConfig, query string, transform func(*sql.Rows) interface{}) interface{} {
	return transform(executeDbQuery(appConfig, query))
}