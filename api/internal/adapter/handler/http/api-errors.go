package http

import "fmt"

func notFoundError(item string) string {
	return fmt.Sprintf("Requested %s not found", item)
}

func invalidRequest() string {
	return "Invalid Request"
}

func missingQueryParam(param string) string {
	return fmt.Sprintf("Query parameter %s missing", param)
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
