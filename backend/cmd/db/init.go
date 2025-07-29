package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	sqlc "sensicore/internal/sqlc/generate"

	_ "github.com/go-sql-driver/mysql"
)

type DataStore struct {
	SQLDB   *sql.DB
	Queries *sqlc.Queries
}

func NewDataStore() (*DataStore, error) {

	ds := new(DataStore)
	var err error
	ds.SQLDB, ds.Queries, err = InitDB()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func InitDB() (*sql.DB, *sqlc.Queries, error) {

	fmt.Println("Connecting to Databases and Cache...")

	ctx := context.Background()

	dbConnStr, exists := os.LookupEnv("MYSQL_DB_CONN_STR")
	if !exists {
		return nil, nil, errors.New("mysql db conn str not found in env")
	}

	sqlDB, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating database pool: %s", err)
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, nil, errors.New("database connection failed : " + err.Error())
	} else {
		fmt.Println("Database connection is alive!")
	}

	queries := sqlc.New(sqlDB)

	return sqlDB, queries, nil
}

// Close Data store connections
func Close(ds *DataStore) error {
	fmt.Println("Closing connections of Data stores...")

	if ds.SQLDB != nil {
		ds.SQLDB.Close()
	}

	return nil
}
