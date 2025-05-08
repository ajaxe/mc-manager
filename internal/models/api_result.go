package models

type ApiResult struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}
