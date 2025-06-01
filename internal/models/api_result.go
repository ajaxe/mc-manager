package models

type ApiResult struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}
type ApiIDResult struct {
	ApiResult
	ID string `json:"id"`
}

func SuccessApiResult() *ApiResult {
	return &ApiResult{
		Success: true,
	}
}

func NewApiIDResult(id string) *ApiIDResult {
	return &ApiIDResult{
		ApiResult: ApiResult{
			Success: true,
		},
		ID: id,
	}
}

type ApiAuthResult struct {
	ApiResult
	RedirectURL string `json:"redirectUrl"`
}

func NewApiAuthResult(url ...string) *ApiAuthResult {
	u := ""
	if len(url) > 0 {
		u = url[0]
	}
	return &ApiAuthResult{
		ApiResult: ApiResult{
			Success: true,
		},
		RedirectURL: u,
	}
}
