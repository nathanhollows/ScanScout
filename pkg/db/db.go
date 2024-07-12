package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func Connect() {
	var sqldb *sql.DB
	var err error

	dataSourceName := os.Getenv("DB_CONNECTION")
	if dataSourceName == "" {
		log.Fatal("DB_CONNECTION not set")
	}

	driverName := os.Getenv("DB_TYPE")
	switch driverName {
	case "mysql":
		sqldb, err = sql.Open(driverName, dataSourceName)
		DB = bun.NewDB(sqldb, mysqldialect.New())
	case "sqlite3":
		sqldb, err = sql.Open(sqliteshim.ShimName, dataSourceName)
		DB = bun.NewDB(sqldb, sqlitedialect.New())
	default:
		log.Fatalf("unsupported DB_TYPE: %s", driverName)
	}

	if err != nil {
		log.Fatal(err)
	}

	DB.AddQueryHook(bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(false),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	))

}
