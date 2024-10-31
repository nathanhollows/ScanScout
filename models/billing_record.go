package models

import "time"

type BillingRecord struct {
	baseModel

	ID        string    `bun:"id,unique,pk,type:varchar(36)"`
	UserID    string    `bun:"user_id,type:varchar(36)"`
	Tier      string    `bun:"tier,type:varchar(50)"`
	StartDate time.Time `bun:"start_date"`
	EndDate   time.Time `bun:"end_date,nullzero"`
	Notes     string    `bun:"notes,type:varchar(255)"`
}
