package models

import "go.mongodb.org/mongo-driver/v2/bson"

type LaunchItem struct {
	LaunchID   bson.ObjectID `bson:"_id" json:"launchId"`
	ID         bson.ObjectID `bson:"_id" json:"id"`
	Name       string        `json:"name"`
	WorldSeed  string        `bson:"world_seed" json:"worldSeed"`
	GameMode   string        `bson:"game_mode" json:"gameMode"`
	LaunchDate string        `bson:"launch_date" json:"launchDate"`
	Status     string        `bson:"status" json:"status"`
	Message    string        `bson:"message" json:"message"`
}
type LaunchItemListResult struct {
	ApiResult
	Data []*LaunchItem `json:"data"`
}
type CreateLaunchItem struct {
	WorldItemID bson.ObjectID `json:"worldItemId"`
}
