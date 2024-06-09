package server

import (
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"strconv"
	"strings"
)

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	// Get the notifications for the user
	userNotifications, err := GetNotifications(sessionUser.Id)
	if err != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// Return the notifications as JSON response
	writeToJson(userNotifications, w)
}

func GetNotifications(userID int) ([]structs.NotificationResponse, error) {
	// Get the user notifications from the database
	userNotifications, err := database.GetUserNotifications(userID)
	if err != nil {
		return nil, err
	}

	// Create a slice to store the notification responses
	var notificationResponses []structs.NotificationResponse

	// Iterate over the user notifications and create notification responses
	for _, notification := range userNotifications {
		var notificationResponse structs.NotificationResponse

		// Set the notification ID
		notificationResponse.Id = notification.ID

		// Set the notification title and body based on the notification type
		if notification.CommentID != 0 {
			// Notification is for a comment
			post, err := database.GetPost(notification.CommentID)
			if err != nil {
				// Handle the error, e.g., log it and continue to the next notification
				continue
			}
			notificationResponse.Title = "New Comment"
			notificationResponse.Body = "A new comment was added to your post: " + post.Title
			notificationResponse.Link = "/post/" + strconv.Itoa(notification.CommentID)
		} else if notification.PostReactionID != 0 {
			// Notification is for a post reaction
			post, err := database.GetPost(notification.PostReactionID)
			if err != nil {
				// Handle the error, e.g., log it and continue to the next notification
				continue
			}
			notificationResponse.Title = "New Reaction"
			notificationResponse.Body = "Someone reacted to your post: " + post.Title
			notificationResponse.Link = "/post/" + strconv.Itoa(notification.PostReactionID)
		}

		// Append the notification response to the slice
		notificationResponses = append(notificationResponses, notificationResponse)
	}

	return notificationResponses, nil
}

func markNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the notification ID from the URL
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	notificationID := urlParts[2]
	notificationIDInt, err := strconv.Atoi(notificationID)
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	sessionUser := sessionmanager.GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	err = database.MarkNotificationAsRead(notificationIDInt)
	if err != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	writeToJson(map[string]string{"message": "Notification marked as read"}, w)
}
