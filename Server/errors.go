package Server

import (
	"RealTimeForum/structs"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string                `json:"message"`
	User    *structs.UserResponse `json:"user,omitempty"`
}

func errorServer(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	view := ErrorResponse{}
	sessionUser := GetUser(r)
	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Id:          sessionUser.Id,
			Username:    sessionUser.Username,
			FirstName:   sessionUser.FirstName,
			LastName:    sessionUser.LastName,
			DateOfBirth: sessionUser.DateOfBirth,
			Location:    sessionUser.Country,
			ImageURL:    GetImageData(sessionUser.ImageId),
			Type:        userTypeToResponse(sessionUser.Type),
		}
	}

	switch code {
	case http.StatusNotFound:
		view.Message = "Resource Not Found"
	case http.StatusInternalServerError:
		view.Message = "Internal Server Error"
	default:
		// as a fallback get a default text for the status code
		log.Printf("errorServer: %d is not implemented\n", code)
		view.Message = http.StatusText(code)
	}

	err := json.NewEncoder(w).Encode(view)
	if err != nil {
		log.Printf("errorServer: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
