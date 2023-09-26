package utils

import "net/http"

func GetHTTPCode(err error) int {
	e := GetError(err)
	if e == nil {
		return http.StatusInternalServerError
	}

	return e.HTTPCode
}
