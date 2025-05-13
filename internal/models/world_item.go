package models

import "go.mongodb.org/mongo-driver/v2/bson"

type WorldItem struct {
	ID          bson.ObjectID `bson:"_id" json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	WorldSeed   string        `bson:"world_seed" json:"worldSeed"`
	GameMode    string        `bson:"game_mode" json:"gameMode"`
	CreateDate  string        `bson:"create_date" json:"createDate"`
}
