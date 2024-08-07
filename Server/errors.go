package Server

import (
	"RealTimeForum/structs"
	"log"
	"net/http"
)

func errorServer(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)

	view := errorView{}
	sessionUser := GetUser(r)
	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			DateOfBirth: sessionUser.DateOfBirth,
			Location:  sessionUser.Country,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
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

	err := templates.ExecuteTemplate(w, "error.html", view)
	if err != nil {
		log.Printf("errorServer: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
