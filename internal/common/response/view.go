package response

import (
	"encoding/json"
)

// ViewResponse - структура ответа сервиса
type ViewResponse struct {
	Result    interface{} `json:"result"`
	Error     string      `json:"error"`
	ErrorType string      `json:"errorType"`
	Code      int         `json:"code"`
}

// Marshal - маршалинг внутренней структуры для корректного json ответа
func (r *ViewResponse) Marshal() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return bytes
}
