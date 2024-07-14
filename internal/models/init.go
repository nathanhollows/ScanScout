package models

import (
	"context"
	"log"
	"time"

	"github.com/nathanhollows/Rapua/internal/contextkeys"
	"github.com/nathanhollows/Rapua/pkg/db"
)

func CreateTables() {
	var models = []interface{}{
		(*CompletionCriteria)(nil),
		(*Location)(nil),
		(*LocationContent)(nil),
		(*Team)(nil),
		(*Marker)(nil),
		(*Scan)(nil),
		(*Instance)(nil),
		(*User)(nil),
	}

	for _, model := range models {
		_, err := db.DB.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
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

// GetUserFromContext gets the user from the context
func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(contextkeys.UserIDKey).(*User)
	if !ok {
		return nil
	}
	return user
}
