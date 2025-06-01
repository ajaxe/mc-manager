package client

const (
	ToastStatusSuccess = "success"
	ToastStatusError   = "error"
)

type StatusToastData struct {
	Status  string
	Message string
}
