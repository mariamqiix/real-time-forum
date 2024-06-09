package server

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/helpers"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
)

func reportPostHandler(w http.ResponseWriter, r *http.Request) {
	// payload should be of type structs.PostReportRequest  // ✅
	// if the user is not logged in, return 401             // ✅
	// if the post is not found, return 404                 // ✅

	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	var reportRequest structs.PostReportRequest

	if !helpers.ParseBody(&reportRequest, r) {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// Add error handling for user not logged in
	session, ok := sessionmanager.LoggedOrNot(w, r)
	if !ok || session == nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	// Add error handling for post not found
	report := structs.Report{
		PostId: reportRequest.PostID,
		Reason: reportRequest.Reason,
	}

	err2 := database.AddReport(report)
	if err2 != nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	err2 = writeToJson(report, w)
	if err2 != nil {
		log.Printf("reportPostHandler: %s\n", err2.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
	writeToJson(map[string]string{"message": "Operation performed successfully"}, w)
}

func reportUserHandler(w http.ResponseWriter, r *http.Request) {
	// payload should be of type structs.PostReportRequest  // ✅
	// if the user is not logged in, return 401             // ✅
	// if the user is not found, return 404                 // ✅

	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	var reportRequest structs.UserReportRequest

	if !helpers.ParseBody(&reportRequest, r) {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// Add error handling for user not logged in
	session, ok := sessionmanager.LoggedOrNot(w, r)
	if !ok || session == nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	// Add error handling for user not found
	user, err := database.GetUserById(*session.UserId) // Dereference the pointer value
	if err != nil || user == nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusNotFound)
		return
	}

	// Add error handling for post not found
	report := structs.Report{
		ReportedId: reportRequest.Username,
		Reason:     reportRequest.Reason,
	}

	err2 := database.AddReport(report)
	if err2 != nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	err2 = writeToJson(report, w)
	if err2 != nil {
		log.Printf("reportPostHandler: %s\n", err2.Error())
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	writeToJson(map[string]string{"message": "Operation performed successfully"}, w)
}
