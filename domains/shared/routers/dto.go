package routers

// JsonErr structures a standard error to return
type JsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// JWTError is a struct that is used to contain a json encoded error message for any JWT related errors
type JWTError struct {
	Message string `json:"message"`
}
