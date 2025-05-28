package models

type PaginationResult[V any] struct {
	Total   int64  `json:"total"`
	Results []*V   `json:"results"`
	NextID  string `json:"nextId"`
	PrevID  string `json:"prevId"`
	HasMore bool   `json:"hasMore"`
}
