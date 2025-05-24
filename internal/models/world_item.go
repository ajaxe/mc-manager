package models

type WorldItem struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WorldSeed   string `bson:"world_seed" json:"worldSeed"`
	GameMode    string `bson:"game_mode" json:"gameMode"`
	CreateDate  string `bson:"create_date" json:"createDate"`
	IsActive    bool   `bson:"is_active" json:"isActive"`
	IsFavorite  bool   `bson:"is_favorite" json:"isFavorite"`
}

type WorldItemListResult struct {
	ApiResult
	Data []*WorldItem `json:"data"`
}
