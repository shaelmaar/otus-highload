package handlers

import "github.com/shaelmaar/otus-highload/social-network/counters/gen/serverhttp"

func Simple500JSONResponse(message string) serverhttp.N5xxJSONResponse {
	if message == "" {
		message = "Внутренняя ошибка сервера"
	}

	return serverhttp.N5xxJSONResponse{
		Body: struct {
			Code      *int    `json:"code,omitempty"`
			Message   string  `json:"message"`
			RequestId *string `json:"request_id,omitempty"`
		}{
			Code:      nil,
			Message:   message,
			RequestId: nil,
		},
		Headers: serverhttp.N5xxResponseHeaders{
			RetryAfter: 0,
		},
	}
}
