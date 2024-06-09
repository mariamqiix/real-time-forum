package helpers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[\w]+@[\w]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsDataImage(buff []byte) (bool, string) {
	// the function that actually does the trick
	t := http.DetectContentType(buff)
	return strings.HasPrefix(t, "image"), t
}

// Returns false if an error happens
func ParseBody(data any, r *http.Request) bool {
	return json.NewDecoder(r.Body).Decode(&data) == nil
}
