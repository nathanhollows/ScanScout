package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type baseModel struct {
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero"`
}

type m20241209135319_NavigationMode int
type m20241209135319_NavigationMethod int
type m20241209135319_CompletionMethod int
type m20241209135319_GameStatus int

type m20241209135319_Notification struct {
	bun.BaseModel `bun:"table:notifications"`
	baseModel

	ID        string `bun:"id,pk,notnull"`
	Content   string `bun:"content,type:varchar(255)"`
	Type      string `bun:"type,type:varchar(255)"`
	TeamCode  string `bun:"team_code,type:varchar(36)"`
	Dismissed bool   `bun:"dismissed,type:bool"`
}

type m20241209135319_InstanceSettings struct {
	bun.BaseModel `bun:"table:instance_settings"`
	baseModel

	InstanceID        string                           `bun:"instance_id,pk,type:varchar(36)"`
	NavigationMode    m20241209135319_NavigationMode   `bun:"navigation_mode,type:int"`
	NavigationMethod  m20241209135319_NavigationMethod `bun:"navigation_method,type:int"`
	MaxNextLocations  int                              `bun:"max_next_locations,type:int,default:3"`
	CompletionMethod  m20241209135319_CompletionMethod `bun:"completion_method,type:int"`
	ShowTeamCount     bool                             `bun:"show_team_count,type:bool"`
	EnablePoints      bool                             `bun:"enable_points,type:bool"`
	EnableBonusPoints bool                             `bun:"enable_bonus_points,type:bool"`
	ShowLeaderboard   bool                             `bun:"show_leaderboard,type:bool"`
}

type m20241209135319_Block struct {
	bun.BaseModel `bun:"table:blocks"`

	ID                 string          `bun:"id,pk,notnull"`
	LocationID         string          `bun:"location_id,notnull"`
	Type               string          `bun:"type,type:int"`
	Data               json.RawMessage `bun:"data,type:jsonb"`
	Ordering           int             `bun:"ordering,type:int"`
	Points             int             `bun:"points,type:int"`
	ValidationRequired bool            `bun:"validation_required,type:bool"`
}

type m20241209135319_TeamBlockState struct {
	bun.BaseModel `bun:"table:team_block_states"`
	baseModel

	TeamCode      string          `bun:"team_code,pk,notnull"`
	BlockID       string          `bun:"block_id,pk,notnull"`
	IsComplete    bool            `bun:"is_complete,type:bool"`
	PointsAwarded int             `bun:"points_awarded,type:int"`
	PlayerData    json.RawMessage `bun:"player_data,type:jsonb"`
}

type m20241209135319_Location struct {
	bun.BaseModel `bun:"table:locations"`
	baseModel

	ID           string                           `bun:"id,pk,notnull"`
	Name         string                           `bun:"name,type:varchar(255)"`
	InstanceID   string                           `bun:"instance_id,notnull"`
	MarkerID     string                           `bun:"marker_id,notnull"`
	ContentID    string                           `bun:"content_id,notnull"`
	Criteria     string                           `bun:"criteria,type:varchar(255)"`
	Order        int                              `bun:"order,type:int"`
	TotalVisits  int                              `bun:"total_visits,type:int"`
	CurrentCount int                              `bun:"current_count,type:int"`
	AvgDuration  float64                          `bun:"avg_duration,type:float"`
	Completion   m20241209135319_CompletionMethod `bun:"completion,type:int"`
	Points       int                              `bun:"points,"`

	Clues    []m20241209135319_Clue   `bun:"rel:has-many,join:id=location_id"`
	Instance m20241209135319_Instance `bun:"rel:has-one,join:instance_id=id"`
	Marker   m20241209135319_Marker   `bun:"rel:has-one,join:marker_id=code"`
	Blocks   []m20241209135319_Block  `bun:"rel:has-many,join:id=location_id"`
}

type m20241209135319_Clue struct {
	bun.BaseModel `bun:"table:clues"`
	baseModel

	ID         string `bun:"id,pk,type:varchar(36)"`
	InstanceID string `bun:"instance_id,notnull"`
	LocationID string `bun:"location_id,notnull"`
	Content    string `bun:"content,type:text"`
}

type m20241209135319_Marker struct {
	bun.BaseModel `bun:"table:markers"`
	baseModel

	Code         string  `bun:"code,unique,pk"`
	Lat          float64 `bun:"lat,type:float"`
	Lng          float64 `bun:"lng,type:float"`
	Name         string  `bun:"name,type:varchar(255)"`
	TotalVisits  int     `bun:"total_visits,type:int"`
	CurrentCount int     `bun:"current_count,type:int"`
	AvgDuration  float64 `bun:"avg_duration,type:float"`

	Locations []m20241209135319_Location `bun:"rel:has-many,join:code=marker_id"`
}

type m20241209135319_CheckIn struct {
	bun.BaseModel `bun:"table:check_ins"`
	baseModel

	InstanceID      string    `bun:"instance_id,notnull"`
	TeamID          string    `bun:"team_code,pk,type:string"`
	LocationID      string    `bun:"location_id,pk,type:string"`
	TimeIn          time.Time `bun:"time_in,type:datetime"`
	TimeOut         time.Time `bun:"time_out,type:datetime"`
	MustCheckOut    bool      `bun:"must_check_out"`
	Points          int       `bun:"points,"`
	BlocksCompleted bool      `bun:"blocks_completed,type:int"`

	Location m20241209135319_Location `bun:"rel:has-one,join:location_id=id"`
}

type m20241209135319_Instance struct {
	bun.BaseModel `bun:"table:instances"`
	baseModel

	ID                    string                     `bun:"id,pk,type:varchar(36)"`
	Name                  string                     `bun:"name,type:varchar(255)"`
	UserID                string                     `bun:"user_id,type:varchar(36)"`
	StartTime             bun.NullTime               `bun:"start_time,nullzero"`
	EndTime               bun.NullTime               `bun:"end_time,nullzero"`
	Status                m20241209135319_GameStatus `bun:"-"`
	IsQuickStartDismissed bool                       `bun:"is_quick_start_dismissed,type:bool"`

	Teams     []m20241209135319_Team           `bun:"rel:has-many,join:id=instance_id"`
	Locations []m20241209135319_Location       `bun:"rel:has-many,join:id=instance_id"`
	Settings  m20241209135319_InstanceSettings `bun:"rel:has-one,join:id=instance_id"`
}
type m20241209135319_Team struct {
	bun.BaseModel `bun:"table:teams"`
	baseModel

	// ID string `bun:"id,pk"

	Code         string `bun:"code,unique,pk"`
	Name         string `bun:"name,"`
	InstanceID   string `bun:"instance_id,notnull"`
	HasStarted   bool   `bun:"has_started,default:false"`
	MustCheckOut string `bun:"must_scan_out"`
	Points       int    `bun:"points,"`

	Instance         m20241209135319_Instance         `bun:"rel:has-one,join:instance_id=id"`
	CheckIns         []m20241209135319_CheckIn        `bun:"rel:has-many,join:code=team_code"`
	BlockingLocation m20241209135319_Location         `bun:"rel:has-one,join:must_scan_out=marker_id,join:instance_id=instance_id"`
	Messages         []m20241209135319_Notification   `bun:"rel:has-many,join:code=team_code"`
	Blocks           []m20241209135319_TeamBlockState `bun:"rel:has-many,join:code=team_code"`
}

type m20241209135319_User struct {
	bun.BaseModel `bun:"table:users"`
	baseModel

	ID               string       `bun:"id,unique,pk,type:varchar(36)"`
	Name             string       `bun:"name,type:varchar(255)"`
	Email            string       `bun:"email,unique,pk"`
	EmailVerified    bool         `bun:"email_verified,type:boolean"`
	EmailToken       string       `bun:"email_token,type:varchar(36)"`
	EmailTokenExpiry sql.NullTime `bun:"email_token_expiry,nullzero"`
	Password         string       `bun:"password,type:varchar(255)"`
	Provider         string       `bun:"provider,type:varchar(255)"`

	Instances         []m20241209135319_Instance `bun:"rel:has-many,join:id=user_id"`
	CurrentInstanceID string                     `bun:"current_instance_id,type:varchar(36)"`
	CurrentInstance   m20241209135319_Instance   `bun:"rel:has-one,join:current_instance_id=id"`
}

func init() {
	var models = []interface{}{
		(*m20241209135319_Notification)(nil),
		(*m20241209135319_InstanceSettings)(nil),
		(*m20241209135319_Block)(nil),
		(*m20241209135319_TeamBlockState)(nil),
		(*m20241209135319_Location)(nil),
		(*m20241209135319_Clue)(nil),
		(*m20241209135319_Team)(nil),
		(*m20241209135319_Marker)(nil),
		(*m20241209135319_CheckIn)(nil),
		(*m20241209135319_Instance)(nil),
		(*m20241209135319_User)(nil),
	}

	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		for _, model := range models {
			_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
			if err != nil {
				return err
			}
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		for _, model := range models {
			_, err := db.NewDropTable().Model(model).IfExists().Exec(context.Background())
			if err != nil {
				return err
			}
		}
		return nil
	})
}
