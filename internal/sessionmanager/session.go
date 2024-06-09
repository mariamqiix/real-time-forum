package sessionmanager

import (
	"errors"
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/structs"
	"time"
)

// Get a user Struct from a request, if user us a Guest.
func LoggedOrNot(w http.ResponseWriter, r *http.Request) (*structs.Session, bool) {
	cookies, err := r.Cookie("SessionToken")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return nil, false
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}
	// Get cookie value
	session, err := database.GetSession(cookies.Value)
	if err != nil {
		log.Printf("error getting session: %s\n", err.Error())
		return nil, false
	}

	err1 := database.AddSession(session)
	if err1 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	return &session, true // return the session
}

// Create a new session and set the cookie
// This function is used when a user logs in
func CreateSessionAndSetCookie(token string, w http.ResponseWriter, user *structs.User) error {
	sID := ""
	if token != "" {
		sID = token
	} else {
		sID = GenerateSession()

	}
	if sID == "" {
		log.Println("error generating session id")
		return errors.New("error generating session id")
	}

	newSession := &structs.Session{
		Token:        sID,
		UserId:       &user.Id,
		CreationTime: time.Now().Unix(),
	}
	// Add the session to the database
	err := database.AddSession(*newSession)
	if err != nil {
		log.Printf("error adding session: %s\n", err.Error())
		return errors.New("error adding session")
	}

	cookie := SessionToCookie(newSession)
	if cookie == nil {
		log.Println("error creating cookie")
		return errors.New("error creating cookie")
	}
	http.SetCookie(w, cookie)

	return nil
}
