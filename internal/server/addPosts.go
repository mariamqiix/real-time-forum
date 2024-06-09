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
	"time"
)

func parsePostForm(data *structs.AddPostRequest, r *http.Request) bool {
	err := r.ParseMultipartForm(22 << 10)
	if err != nil {
		log.Printf("parseForm: %s\n", err.Error())
		return false
	}
	data.Title = r.FormValue("title")
	data.Content = r.FormValue("content")
	data.Categories = r.Form["categories"]

	_, fh, err := r.FormFile("image")
	if err == nil {
		data.Image = *fh
	}
	return true
}

func addPostGetHandler(w http.ResponseWriter, r *http.Request) {
	// this endpoint will just serve the addPost.html template
	// make sure use is logged in
	// if not, redirect to login

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

func addPostPostHandler(w http.ResponseWriter, r *http.Request) {
	// takes a post request and adds a post to the database
	// post is in structs.AddPostRequest
	// make sure contentLength is (maxImageSize + idk maybe 2 MB)
	// maxPostSize := maxImageSize + 2*1024*1024
	// on success, redirect to the post

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
	var pstReq structs.AddPostRequest

	if !parsePostForm(&pstReq, r) {
		errorServer(w, r, http.StatusBadRequest)
		http.Error(w, "Invalid form submitted", http.StatusBadRequest)
		return
	}

	if pstReq.Title == "" || pstReq.Content == "" || len(pstReq.Categories) == 0 {
		log.Println("addPostPostHandler: failed validation")
		http.Error(w, "Empty fields are not allowed", http.StatusBadRequest)
		return
	}

	categoriesIds := make([]int, len(pstReq.Categories))
	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("addPostPostHandler: %s\n", err.Error())
		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		return
	}

	for i, catName := range pstReq.Categories {
		id := getCategoryIdFromArr(dbCategories, catName)
		if id == -1 {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}
		categoriesIds[i] = id
	}

	dbPostAdd := structs.Post{
		UserId:        sessionUser.Id,
		Title:         pstReq.Title,
		Message:       pstReq.Content,
		ImageId:       -1,
		Time:          time.Now().UTC(),
		CategoriesIDs: categoriesIds,
		ParentId:      nil,
	}

	haveImage := pstReq.Image.Size != 0
	if haveImage && pstReq.Image.Size > maxImageSize {
		http.Error(w, "Image size too large", http.StatusBadRequest)
		return
	}

	if haveImage {
		file, err := pstReq.Image.Open()
		if err != nil {
			log.Printf("addPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}

		imageBuff, err := io.ReadAll(file)
		if err != nil {
			log.Printf("addPostPostHandler: %s\n", err.Error())
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
			log.Printf("addPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}
		log.Printf("addPostPostHandler: image uploaded with id %d\n", imageId)
		dbPostAdd.ImageId = imageId
	}

	newPostId, err := database.AddPost(dbPostAdd)
	if err != nil {
		log.Printf("addPostPostHandler: %s\n", err.Error())
		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), http.StatusSeeOther)
}
