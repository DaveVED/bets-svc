package modals

type SuccessResponse struct {
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []Error  `json:"errors"`
}

type Error struct {
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
