package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func InitDB() {
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
		db = bun.NewDB(sqldb, mysqldialect.New())
	case "sqlite3":
		sqldb, err = sql.Open(sqliteshim.ShimName, dataSourceName)
		db = bun.NewDB(sqldb, sqlitedialect.New())
	default:
		log.Fatalf("unsupported DB_TYPE: %s", driverName)
	}

	if err != nil {
		log.Fatal(err)
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(false),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	))

	var models = []interface{}{
		(*Team)(nil),
		(*Location)(nil),
		(*Scan)(nil),
		(*Instance)(nil),
		(*User)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

}

type baseModel struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type belongsToInstance struct {
	InstanceID string `bun:",notnull" json:"instance_id"`
}
