package models

type ApiResult struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}
type ApiIDResult struct {
	ApiResult
	ID string `json:"id"`
}

func NewApiIDResult(id string) *ApiIDResult {
	return &ApiIDResult{
		ApiResult: ApiResult{
			Success: true,
		},
		ID: id,
	}
}
