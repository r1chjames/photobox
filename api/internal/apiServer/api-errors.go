package apiServer

import "fmt"

func notFoundError(item string) string {
	return fmt.Sprintf("Requested %s not found", item)
}

func missingQueryParam(param string) string {
	return fmt.Sprintf("Query parameter %s missing", param)
}