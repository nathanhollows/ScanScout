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

func MustOpen() *bun.DB {
	var sqldb *sql.DB
	var err error
	var db *bun.DB

	dataSourceName := os.Getenv("DB_CONNECTION")
	if dataSourceName == "" {
		log.Fatal("DB_CONNECTION not set. Please set DB_CONNECTION in the environment")
	}

	driverName := os.Getenv("DB_TYPE")
	switch driverName {
	case "mysql":
		sqldb, err = sql.Open(driverName, dataSourceName)
		db = bun.NewDB(sqldb, mysqldialect.New())
	case "sqlite3":
		sqldb, err = sql.Open(sqliteshim.ShimName, dataSourceName)
		db = bun.NewDB(sqldb, sqlitedialect.New())
	default:
		panic("unsupported DB_TYPE: " + driverName + ". Supported types are mysql and sqlite3")
	}

	if err != nil {
		panic(err)
	}

	debugEnabled := os.Getenv("BUNDEBUG") == "1" || os.Getenv("BUNDEBUG") == "2"
	db.AddQueryHook(bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(debugEnabled),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}
