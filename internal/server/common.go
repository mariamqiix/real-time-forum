package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// converts `v`	to json and writes it to the response writer
func writeToJson(v any, w http.ResponseWriter) error {
	buff, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(buff)
	return err
}

func isAcceptedMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		errorServer(w, r, http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func imageIdToUrl(imageId int) string {
	return "/uploads/" + strconv.Itoa(imageId)
}
