package server

import (
	"io"
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/helpers"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"strconv"
)

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
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

	err = database.RemovePost(post.Id)
	if err != nil {
		log.Printf("deletePostHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// redirect back to the last page
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)

}

func editPostGetHandler(w http.ResponseWriter, r *http.Request) {
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
		ParentId:   -1,
		OriginalId: post.Id,
	}
	if post.ParentId != nil {
		addPostView.ParentId = *post.ParentId
	}

	err = templates.ExecuteTemplate(w, "editPost.html", addPostView)
	if err != nil {
		log.Printf("addPostGetHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

func editPostPostHandler(w http.ResponseWriter, r *http.Request) {
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

	var newPostInfo structs.AddPostRequest

	if !parsePostForm(&newPostInfo, r) {
		errorServer(w, r, http.StatusBadRequest)
		http.Error(w, "Invalid form submitted", http.StatusBadRequest)
		return
	}

	if newPostInfo.Title == "" || newPostInfo.Content == "" {
		log.Println("addPostPostHandler: failed validation")
		http.Error(w, "Empty fields are not allowed", http.StatusBadRequest)
		return
	}

	// update the post
	post.Title = newPostInfo.Title
	post.Message = newPostInfo.Content

	haveImage := newPostInfo.Image.Size != 0
	if haveImage && newPostInfo.Image.Size > maxImageSize {
		http.Error(w, "Image size too large", http.StatusBadRequest)
		return
	}

	if haveImage {
		file, err := newPostInfo.Image.Open()
		if err != nil {
			log.Printf("editPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}

		imageBuff, err := io.ReadAll(file)
		if err != nil {
			log.Printf("editPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}

		isImage, _ := helpers.IsDataImage(imageBuff)
		if !isImage {
			errorServer(w, r, http.StatusUnsupportedMediaType)
			http.Error(w, "file type not allowed", http.StatusUnsupportedMediaType)
			return
		}
		imageId, err := database.UploadImage(imageBuff)
		if err != nil {
			log.Printf("editPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}
		post.ImageId = imageId
	}

	database.UpdatePostInfo(post)

	// redirect to the post
	http.Redirect(w, r, "/post/"+strconv.Itoa(post.Id), http.StatusFound)

}
