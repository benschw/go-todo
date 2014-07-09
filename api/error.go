package api

type Error struct {
	Error string `json:"error"`
}

func NewError(msg string) *Error {
	return &Error{Error: msg}
}
