package models

type PlayTimerItem struct {
	ID             string `json:"id" bson:"_id"`
	Minutes        int    `bson:"minutes" json:"minutes"`
	EndDate        string `bson:"end_date" json:"endDate"`
	IsActive       bool   `bson:"is_active" json:"isActive"`
	CreateDate     string `bson:"create_date" json:"createDate"`
	LastUpdateDate string `bson:"last_update_date" json:"lastUpdateDate"`
}

type PlayTimerListResult struct {
	ApiResult
	Data []*PlayTimerItem `json:"data"`
}
