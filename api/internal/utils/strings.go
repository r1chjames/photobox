package utils

import (
	"fmt"
	"os"
	"strings"
)

func AllIndicesOfChar(str string, charToFind string) []int {
	var indices []int
	for pos, char := range str {
		if string(char) == charToFind {
			indices = append(indices, pos)
		}
	}
	return indices
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if len(fallback) > 0 {
			return fallback
		} else {
			panic(fmt.Sprintf("Required application variable %s not defined", key))
		}
	}
	return value
}

func Exists(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func EscapeInvalidCharacters(field string) string {
	return strings.ReplaceAll(field, "'", "\\'")
}
