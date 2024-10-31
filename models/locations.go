package models

type Location struct {
	baseModel

	ID           string           `bun:"id,pk,notnull"`
	Name         string           `bun:"name,type:varchar(255)"`
	InstanceID   string           `bun:"instance_id,notnull"`
	MarkerID     string           `bun:"marker_id,notnull"`
	ContentID    string           `bun:"content_id,notnull"`
	Criteria     string           `bun:"criteria,type:varchar(255)"`
	Order        int              `bun:"order,type:int"`
	TotalVisits  int              `bun:"total_visits,type:int"`
	CurrentCount int              `bun:"current_count,type:int"`
	AvgDuration  float64          `bun:"avg_duration,type:float"`
	Completion   CompletionMethod `bun:"completion,type:int"`
	Points       int              `bun:"points,"`

	Clues    []Clue   `bun:"rel:has-many,join:id=location_id"`
	Instance Instance `bun:"rel:has-one,join:instance_id=id"`
	Marker   Marker   `bun:"rel:has-one,join:marker_id=code"`
	Blocks   []Block  `bun:"rel:has-many,join:id=location_id"`
}
