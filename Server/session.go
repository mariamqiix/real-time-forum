package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
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

const SESSION_EXPIRY = 30 * 24 * 60 * 60

// Get a user Struct from a request, if user us a Guest.
// response will only have the type set to structs.UserTypeIdGuest
func GetUser(r *http.Request) *structs.User {
	sessionCookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return nil
	}
	sessionExists, _ := database.CheckExistance("UserSession", "token", sessionCookie.Value)
	if !sessionExists {
		return nil
	}
	session, err := database.GetSession(sessionCookie.Value)
	if err != nil {
		log.Printf("error getting session: %s\n", err.Error())
		return nil
	}
	// a nill foreign key means the user is a guest
	if session.UserId == nil {
		return &structs.User{
			Type: structs.UserTypeIdGuest,
		}
	}
	user, err := database.GetUserById(*session.UserId)
	if err != nil {
		log.Printf("error getting user by user id: %s\n", err.Error())
		return nil
	}
	return user
}

// converts a session struct to a http.Cookie
func SessionToCookie(session *structs.Session) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Expires:  time.Unix(session.CreationTime+SESSION_EXPIRY, 0),
		HttpOnly: true, // to prevent XSS
		Path:     "/",
	}
}

// generates a new session token based on a UUIDv7, returns "" if an error occurs
func GenerateSession() string {
	uuid, err := uuid.NewV7()
	if err != nil {
		return ""
	}
	buff, err := uuid.MarshalText()
	if err != nil {
		return ""
	}
	return string(buff)
}
