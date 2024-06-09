package sessionmanager

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/structs"
	"time"

	"github.com/google/uuid"
)

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
