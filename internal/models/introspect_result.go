package models

import "time"

type IntrospectResult struct {
	Active    bool      `json:"active"`
	Username  string    `json:"username"`
	IssuedUtc time.Time `json:"issuedUtc"`
}
