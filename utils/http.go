package utils

func GetHTTPCode(err error) int {
	e := GetError(err)
	if e == nil {
		return 0
	}

	return e.HTTPCode
}
