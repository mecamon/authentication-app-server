package helpers

type ErrorsMap struct {
	Success bool              `json:"success"`
	Message map[string]string `json:"message"`
}
