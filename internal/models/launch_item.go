package models

const (
	PageDirectionNext = "next"
	PageDirectionPrev = "prev"
)

type LaunchItem struct {
	ID         string `bson:"_id" json:"id"`
	WorldID    string `bson:"world_id" json:"worldId"`
	Name       string `json:"name"`
	WorldSeed  string `bson:"world_seed" json:"worldSeed"`
	GameMode   string `bson:"game_mode" json:"gameMode"`
	LaunchDate string `bson:"launch_date" json:"launchDate"`
	Status     string `bson:"status" json:"status"`
	Message    string `bson:"message" json:"message"`
}

type LaunchItemListRequest struct {
	Direction string
	CursorID  string
}

type LaunchItemListResult struct {
	ApiResult
	PaginationResult[LaunchItem]
}
type CreateLaunchItem struct {
	WorldItemID string `json:"worldItemId"`
	GameMode    string `json:"gameMode"`
}

// ToLaunchItem creates a new LaunchItem from a WorldItem
func ToLaunchItem(w *WorldItem, launchDt string, status string) *LaunchItem {
	return &LaunchItem{
		WorldID:    w.ID,
		Name:       w.Name,
		WorldSeed:  w.WorldSeed,
		GameMode:   w.GameMode,
		LaunchDate: launchDt,
		Status:     status,
		Message:    "",
	}
}
