package server

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"strconv"
	"time"
)

// same logic as addPostGetHandler
func addCommentGetHandler(w http.ResponseWriter, r *http.Request) {

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

	postId, err := strconv.Atoi(r.PathValue("post_id"))
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	post, err := database.GetPost(postId)
	if err != nil || post == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	// fill the view data
	addPostView := addPostView{
		User: &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		},
		Categories: nil,
		ParentId:   postId,
	}

	categories, err := database.GetCategories()
	if err != nil {
		log.Printf("addPostGetHandler: %s\n", err.Error())
		err = nil
	} else {
		addPostView.Categories = mapCategories(categories)
	}

	err = templates.ExecuteTemplate(w, "newPost.html", addPostView)
	if err != nil {
		log.Printf("addPostGetHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

}

// same logic as addPostPostHandler
func addCommentPostHandler(w http.ResponseWriter, r *http.Request) {
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

	// check post actually exists
	postId, err := strconv.Atoi(r.PathValue("post_id"))
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	post, err := database.GetPost(postId)
	if err != nil || post == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	var pstReq structs.AddPostRequest
	if !parsePostForm(&pstReq, r) {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	if pstReq.Title == "" || pstReq.Content == "" {
		log.Println("addCommentGetHandler: failed validation")
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	dbPostAdd := structs.Post{
		ParentId: &postId,
		UserId:   sessionUser.Id,
		Title:    pstReq.Title,
		Message:  pstReq.Content,
		ImageId:  -1,
		Time:     time.Now().UTC(),
	}

	commentID, err := database.AddPost(dbPostAdd)
	if err != nil {
		log.Printf("addCommentPostHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// Create a new notification for the comment
	notification := structs.UserNotification{
		CommentID: commentID,
	}

	// Add the notification to the database
	_, err = database.AddNotification(notification)
	if err != nil {
		log.Printf("Failed to add notification: %v", err)
		// Handle the error appropriately
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postId), http.StatusSeeOther)
}
