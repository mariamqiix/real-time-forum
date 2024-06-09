package server

import (
	"io"
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/helpers"
	"sandbox/internal/sessionmanager"
	"strconv"
)

const maxImageSize = 20971520 // 20 MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	if r.ContentLength <= 0 || r.ContentLength > maxImageSize {
		errorServer(w, r, http.StatusRequestEntityTooLarge)
		return
	}

	// Parse the form data, including the uploaded file
	err := r.ParseMultipartForm(maxImageSize)
	if err != nil {
		log.Printf("uploadHandler: unable to parse form: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// Get a reference to the file
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("uploadHandler: unable to get file: %s\n", err.Error())
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file data into a buffer
	buff, err := io.ReadAll(file)
	if err != nil {
		log.Printf("uploadHandler: unable to read file: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// Check if the content is an image
	isImage, typeOfImage := helpers.IsDataImage(buff)
	if !isImage {
		log.Printf("uploadHandler: not an image: %s\n", typeOfImage)
		errorServer(w, r, http.StatusUnsupportedMediaType)
		return
	}

	// Upload the image data to the database
	_, err = database.UploadImage(buff)
	if err != nil {
		log.Printf("uploadHandler: unable to upload image to the database: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
}
func uploadedContentServerHandler(w http.ResponseWriter, r *http.Request) {
	imageID, err := strconv.Atoi(r.PathValue("image_id"))
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	imageData, err := database.GetImage(imageID)
	if err != nil {
		log.Printf("uploadedContentServerHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	if imageData == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	w.Write(imageData)
}
