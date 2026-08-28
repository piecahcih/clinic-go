package httperr

type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string { return e.Message }

func NotFound(msg string) *APIError      { return &APIError{404, "not_found", msg} }
func BadRequest(msg string) *APIError    { return &APIError{400, "bad_request", msg} }
func Conflict(msg string) *APIError      { return &APIError{409, "conflict", msg} }
func Unprocessable(msg string) *APIError { return &APIError{422, "unprocessable", msg} }
