package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

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

const maxImageSize = 20971520 // 20 MB

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Get the user from the session
	sessionUser := GetUser(r)

	// Set the limiter username
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}

	// Check if the rate limiter allows the request
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	// Check the request body size
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

	// Get the uploaded file
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
	isImage, typeOfImage := IsDataImage(buff)
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

	// Respond to the client with a 201 Created status
	w.WriteHeader(http.StatusCreated)

}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get the user from the session
	sessionUser := GetUser(r)

	// Set the limiter username
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}

	// Check if the user is allowed to make the request
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	userId, _ := strconv.Atoi(r.FormValue("user_id"))
	var dbUser *structs.User
	if userId != -1 {
		dbUser, err := database.GetUserById(userId)
		if err != nil {
			log.Printf("profileHandler: %s\n", err.Error())
			errorServer(w, r, http.StatusInternalServerError)
			return
		}

		if dbUser == nil {
			errorServer(w, r, http.StatusNotFound)
			return
		}

	} else {
		dbUser = sessionUser
	}
	// Fetch the user from the database

	// Prepare the view data
	view := profileView{}
	view.UserProfile = structs.UserResponse{
		Username:  dbUser.Username,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		ImageURL:  imageIdToUrl(dbUser.ImageId),
		Type:      userTypeToResponse(dbUser.Type),
	}

	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		}
	} else {
		sessionUser = &structs.User{Id: -1}
	}
	casees := r.FormValue("case")
	// Fetch the posts based on the query parameter
	switch casees {
	case "comments":
		posts, err := database.GetPostsByUser(dbUser.Id, -1, 0, true)
		if err == nil {
			view.Posts = mapPosts(posts, sessionUser.Id)
		}
	case "likes":
		posts, err := database.GetUserReactions(dbUser.Id)
		if err == nil {
			filterdPosts := filterPostsByReaction(posts, dbUser.Id, structs.PostReactionTypeLike)
			view.Posts = mapPosts(filterdPosts, sessionUser.Id)
		}
	case "dislikes":
		posts, err := database.GetUserReactions(dbUser.Id)
		if err == nil {
			filterdPosts := filterPostsByReaction(posts, dbUser.Id, structs.PostReactionTypeDislike)
			view.Posts = mapPosts(filterdPosts, sessionUser.Id)
		}
	default:
		posts, err := database.GetPostsByUser(dbUser.Id, -1, 0, false)
		if err == nil {
			view.Posts = mapPosts(posts, sessionUser.Id)
		}
	}
	writeToJson(view, w)
}

func signupPostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Retrieve the form values
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	username := r.FormValue("username")
	country := r.FormValue("country")
	gender := r.FormValue("gender")
	dobString := r.FormValue("dob")
	password := r.FormValue("password")

	// check password length
	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// check if the username exists
	exist, err := database.CheckExistance("User", "username", username)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}

	if exist {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// check if the email is valid and exists
	if !IsValidEmail(email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	exist, err = database.CheckExistance("User", "email", email)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}

	if exist {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	imageID := 0
	// if userData.Image != nil {
	// 	isImage, _ := IsDataImage(userData.Image)
	// 	if isImage {
	// 		imageID, err = database.UploadImage(userData.Image)
	// 		if err != nil {
	// 			log.Printf("SignupHandler: %s\n", err.Error())
	// 		}
	// 	}
	// }

	// structure
	hashedPassword, hashErr := GetHash(password)
	if hashErr != HasherErrorNone {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}
	layout := "2006-01-02"
	dob, _ := time.Parse(layout, dobString)
	cleanedUserData := structs.User{
		Username:       username,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		Country:        country,
		DateOfBirth:    dob,
		HashedPassword: hashedPassword,
		ImageId:        imageID,
		GithubName:     "",
		LinkedinName:   "",
		TwitterName:    "",
		Bio:            "",
		Gender:         gender,
	}
	fmt.Print("cleanedUserData", cleanedUserData.Gender)
	err = database.CreateUser(cleanedUserData)
	if err != nil {
		http.Error(w, "could not create a user, please try again later", http.StatusBadRequest)
		return
	}
	finalUser, err := database.GetUserByUsername(cleanedUserData.Username)
	if err != nil {
		return
	}
	err = CreateSessionAndSetCookie("", w, finalUser)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.ContentLength > 1024 {
		http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}

	var user *structs.User
	var err error

	Username := r.FormValue("username")
	Password := r.FormValue("password")
	// check if the email is valid and exists
	if IsValidEmail(Username) {
		exist, err := database.CheckExistance("User", "email", Username)
		if err != nil {
			log.Printf("loginPostHandler: %s\n", err.Error())
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		} else if !exist {
			http.Error(w, "Invalid username or email or password", http.StatusConflict)
			return
		}

		user, err = database.GetUserByEmail(Username)
		if err != nil {
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		}

	} else {
		user, err = database.GetUserByUsername(Username)
		if err != nil {
			log.Printf("loginPostHandler: %s\n", err.Error())
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		}
	}

	if user == nil {
		http.Error(w, "Invalid username or email", http.StatusConflict)
		return
	}

	if err := CompareHashAndPassword(user.HashedPassword, Password); err != HasherErrorNone {
		http.Error(w, "Invalid password", http.StatusConflict)
		return
	}

	if !user.BannedUntil.IsZero() && user.BannedUntil.After(time.Now()) {
		http.Error(w, "User is blocked", http.StatusForbidden)
		return
	}

	// Create a new session and set the cookie
	err2 := CreateSessionAndSetCookie("", w, user)
	if err2 != nil {
		http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
		return
	}

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get the session token from the cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		writeToJson(err, w)
	}
	// Get the session from the database
	sessionToken := cookie.Value
	session, err := database.GetSession(sessionToken)
	if err != nil {
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		writeToJson(err, w)

	}

	// Remove the session from the database
	if err := database.RemoveSession(session.Id); err != nil {
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		writeToJson(err, w)

	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Expires: time.Unix(0, 0),
	})
}

func categoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	categoryName := r.PathValue("category_name")

	categoryExists, _ := database.CheckExistance("Category", "name", categoryName)
	if !categoryExists {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	posts_count, err := database.GetPostsCountByCategory(categoryName)
	if err != nil {
		log.Printf("error getting posts count by category: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	dbPosts, err := database.GetPostsByCategory(categoryName, posts_count, 0, "latest")
	if err != nil {
		log.Printf("error getting posts by category: %s\n", err.Error())
	}

	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
	}

	mappedPosts := mapPosts(dbPosts, -1)
	mappedCategories := mapCategories(dbCategories)

	view := homeView{
		Posts:       mappedPosts,
		User:        nil,
		Categories:  mappedCategories,
		SortOptions: []string{"latest", "most liked", "least liked", "oldestq"},
	}

	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		}
	}

	writeToJson(view, w)

}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
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
	if post.ParentId != nil {
		newPath := "/post/" + strconv.Itoa(*post.ParentId) + "#" + r.PathValue("post_id")
		http.Redirect(w, r, newPath, http.StatusFound)
		return
	}

	view := discussionView{
		User:     nil,
		Post:     structs.PostResponse{},
		Comments: nil,
	}

	comments, err := database.GetCommentsForPost(post.Id, -1, 0)
	if err != nil {
		log.Printf("postsHandler: %s\n", err.Error())
	}

	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		}
		view.Post = mapPosts([]structs.Post{*post}, sessionUser.Id)[0]
		view.Comments = mapPosts(comments, sessionUser.Id)
	} else {
		view.Post = mapPosts([]structs.Post{*post}, -1)[0]
		view.Comments = mapPosts(comments, -1)

	}

	writeToJson(view, w)

}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
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
}

func editPostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "GET" {
		sessionUser := GetUser(r)
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
		writeToJson(addPostView, w)
	} else if r.Method == "POST" {
		sessionUser := GetUser(r)
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

			isImage, _ := IsDataImage(imageBuff)
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
}

func addPostHandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sessionUser := GetUser(r)
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

		writeToJson(addPostView, w)
	}
}

func addPostHandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		sessionUser := GetUser(r)
		limiterUsername := "[GUESTS]"
		if sessionUser != nil {
			limiterUsername = sessionUser.Username
		}

		title := r.FormValue("title")
		Content := r.FormValue("topic")
		Categories := r.FormValue("selectedCategories")

		if !userLimiter.Allow(limiterUsername) {
			errorServer(w, r, http.StatusTooManyRequests)
			return
		}

		if sessionUser == nil {

			errorServer(w, r, http.StatusUnauthorized)
			return
		}

		if title == "" || Content == "" || len(Categories) == 0 {
			log.Println("addPostPostHandler: failed validation")
			http.Error(w, "Empty fields are not allowed", http.StatusBadRequest)
			return
		}

		categoriesIds := make([]int, len(Categories))

		for i, catName := range Categories {
			intId, _ := strconv.Atoi(string(catName))
			categoriesIds[i] = intId
		}

		dbPostAdd := structs.Post{
			UserId:        sessionUser.Id,
			Title:         title,
			Message:       Content,
			ImageId:       -1,
			Time:          time.Now().UTC(),
			CategoriesIDs: categoriesIds,
			ParentId:      nil,
		}

		// haveImage := pstReq.Image.Size != 0
		// if haveImage && pstReq.Image.Size > maxImageSize {
		// 	http.Error(w, "Image size too large", http.StatusBadRequest)
		// 	return
		// }

		// if haveImage {
		// 	file, err := pstReq.Image.Open()
		// 	if err != nil {
		// 		log.Printf("addPostPostHandler: %s\n", err.Error())
		// 		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		// 		return
		// 	}

		// 	imageBuff, err := io.ReadAll(file)
		// 	if err != nil {
		// 		log.Printf("addPostPostHandler: %s\n", err.Error())
		// 		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		// 		return
		// 	}

		// 	isImage, _ := IsDataImage(imageBuff)
		// 	if !isImage {
		// 		errorServer(w, r, http.StatusUnsupportedMediaType)
		// 		http.Error(w, "file type not allowed", http.StatusUnsupportedMediaType)
		// 		return
		// 	}
		// 	imageId, err := database.UploadImage(imageBuff)
		// 	if err != nil {
		// 		log.Printf("addPostPostHandler: %s\n", err.Error())
		// 		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		// 		return
		// 	}
		// 	log.Printf("addPostPostHandler: image uploaded with id %d\n", imageId)
		// 	dbPostAdd.ImageId = imageId
		// }

		_, err := database.AddPost(dbPostAdd)
		if err != nil {
			log.Printf("addPostPostHandler: %s\n", err.Error())
			http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
			return
		}
	}
	writeToJson(map[string]string{"message": "Operation performed successfully"}, w)

}

func reportUserHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	var reportRequest structs.UserReportRequest

	if !ParseBody(&reportRequest, r) {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// Add error handling for user not logged in
	session, ok := LoggedOrNot(w, r)
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

func reportPostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	var reportRequest structs.PostReportRequest

	if !ParseBody(&reportRequest, r) {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// Add error handling for user not logged in
	session, ok := LoggedOrNot(w, r)
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

func postReactionHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	// get the post id and reaction type
	postId := r.PathValue("post_id")
	reactionType := r.PathValue("reaction_type")

	// validate the reaction type
	if reactionType != string(structs.PostReactionTypeLike) && reactionType != string(structs.PostReactionTypeDislike) && reactionType != string(structs.PostReactionTypeLove) && reactionType != string(structs.PostReactionTypeHaha) && reactionType != string(structs.PostReactionTypeSkull) {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// validate the post id and covert it to int
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// get the post struct from the database
	PostStructForRec, err := database.GetPost(postIdInt)
	if err != nil || PostStructForRec == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	// get the user from the session
	UsersPost := GetUser(r)
	if UsersPost == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	mappedReaction := mapReactionForPost(PostStructForRec, UsersPost.Id, structs.PostReactionType(reactionType), reactionType)
	if mappedReaction == nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	reactionId, err2 := database.GetReactionId(mappedReaction.Type)
	if err2 != nil || reactionId == 0 {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	reaction := structs.PostReaction{
		PostId:     PostStructForRec.Id,
		UserId:     UsersPost.Id,
		ReactionId: reactionId,
	}

	err3 := database.AddReactionToPost(reaction)
	if err3 != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
	// add notification
	// if the user is not the owner of the post
	if UsersPost.Id != PostStructForRec.UserId {
		notification := structs.UserNotification{
			PostReactionID: reactionId,
		}
		_, err4 := database.AddNotification(notification)
		if err4 != nil {
			log.Printf("postReactionHandler: %s\n", err4.Error())
		}
	}
}

func deletePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	postId := r.PathValue("post_id")
	reactionType := r.PathValue("reaction_type")

	// validate the reaction type
	if reactionType != string(structs.PostReactionTypeLike) && reactionType != string(structs.PostReactionTypeDislike) && reactionType != string(structs.PostReactionTypeLove) && reactionType != string(structs.PostReactionTypeHaha) && reactionType != string(structs.PostReactionTypeSkull) {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// validate the post id and covert it to int
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// get the post struct from the database
	PostStructForRec, err := database.GetPost(postIdInt)
	if err != nil || PostStructForRec == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	UsersPost := GetUser(r)
	if UsersPost == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	mappedReaction := mapReactionForPost(PostStructForRec, UsersPost.Id, structs.PostReactionType(reactionType), reactionType)
	if mappedReaction == nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	reactionId, err2 := database.GetReactionId(mappedReaction.Type)
	if err2 != nil {
		log.Printf("deletePostReactionHandler: %s\n", err2.Error())

		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	err3 := database.RemoveReactionFromPost(PostStructForRec.Id, reactionId, UsersPost.Id)
	if err3 != nil {
		log.Printf("deletePostReactionHandler: %s\n", err3.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

// same logic as addPostPostHandler
func addCommentPostHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
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

// same logic as addPostGetHandler
func addCommentGetHandler(w http.ResponseWriter, r *http.Request) {

	sessionUser := GetUser(r)
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

	writeToJson(addPostView, w)
}

func categoryGetHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	categories, err := database.GetCategories()
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		http.Error(w, "cannot get the categories", http.StatusNotFound)
		return
	}
	writeToJson(categories, w)
}

// homepageHandler handles the homepage route and serves the homepage template
func searchPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	view := homeView{
		Posts:      nil,
		User:       nil,
		Categories: nil,
	}

	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		}
	}

	content := r.FormValue("search")

	dbPosts, err := database.SearchContent(content)
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
	}

	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
	}

	if sessionUser == nil {
		view.Posts = mapPosts(dbPosts, -1)
	} else {
		view.Posts = mapPosts(dbPosts, sessionUser.Id)
	}

	view.Categories = mapCategories(dbCategories)
	view.SortOptions = []string{"latest", "most liked", "least liked", "oldest"}

	writeToJson(view, w)

}

// homepageHandler handles the homepage route and serves the homepage template
func LoggedUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
		writeToJson(limiterUsername, w)

	}
	writeToJson(nil, w)

}

func UserTypeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}
	writeToJson(userTypeToResponse(sessionUser.Type), w)
}

func ModeratorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	Moderators, err := database.GetUsersByType(2)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	writeToJson(Moderators, w)
}

func RemoveModeratorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	user := r.FormValue("id")
	intUser, _ := strconv.Atoi(user)
	err := database.UpdateUserType(intUser, 1)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
}

func PromotionRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	reason := r.FormValue("answer")

	promotionRequest := structs.PromoteRequest{
		Reason:    reason,
		Time:      time.Now(),
		IsPending: true,
	}
	// Add error handling for user not found
	user, err := database.GetUserById(sessionUser.Id) // Dereference the pointer value
	if err != nil || user == nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusNotFound)
		return
	}

	// Add error handling for post not found
	promotionRequest.UserId = sessionUser.Id

	err2 := database.AddPromoteRequest(promotionRequest)
	if err2 != nil {
		writeToJson(map[string]string{"message": "Could not perform operation, please try again later"}, w)
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	err2 = writeToJson(promotionRequest, w)
	if err2 != nil {
		log.Printf("reportPostHandler: %s\n", err2.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	password := r.FormValue("password")
	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	hashedPassword, hashErr := GetHash(password)
	if hashErr != HasherErrorNone {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	sessionUser.HashedPassword = hashedPassword
	err := database.UpdateUserInfo(sessionUser)
	if err != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func userMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func PromotionRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	promotionRequests, err := database.GetPromoteRequests()
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	PromoteRequestResponse, err := ConvertToPromoteRequestResponse(promotionRequests)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	writeToJson(PromoteRequestResponse, w)
}

func ShowUserPromotionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	id := r.FormValue("id")
	IntId, _ := strconv.Atoi(id)
	promotionRequests, err := database.GetPromoteRequestByid(IntId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	user, err := database.GetUserById(promotionRequests.UserId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	PromoteRequestResponse := structs.PromoteRequestResponse{
		Id:        promotionRequests.Id,
		UserId:    user.Id,
		Username:  user.Username,
		Reason:    promotionRequests.Reason,
		IsPending: promotionRequests.IsPending,
	}
	writeToJson(PromoteRequestResponse, w)

}

func RejectPromotionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	id := r.FormValue("userId")
	IntId, _ := strconv.Atoi(id)
	err := database.ReomvePromoteRequest(IntId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
}

func ApprovePromotionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	id := r.FormValue("userId")
	IntId, _ := strconv.Atoi(id)
	err := database.UpdateUserType(IntId, 2)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	err = database.ReomvePromoteRequest(IntId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
}

func removeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	categoryId := r.FormValue("id")
	IntId, _ := strconv.Atoi(categoryId)
	err := database.RemoveCategory(IntId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	err = database.RemovePostCategory(IntId)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
}
func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)

	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}
	categoryName := r.FormValue("name")
	categoryDescription := r.FormValue("description")
	category := structs.Category{
		Name:        categoryName,
		Description: categoryDescription,
		Color:       "#000000",
	}
	err := database.AddCategory(category)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	user, err := database.GetUserById(sessionUser.Id)
	if err != nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	writeToJson(user, w)
}
func updateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	sessionUser := GetUser(r)
	if sessionUser == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	dateOfBirth := r.FormValue("dateOfBirth")
	fmt.Print("\n\n\n" + dateOfBirth)
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	country := r.FormValue("country")
	gender := r.FormValue("gender")
	newUserName := sessionUser.Username

	if username != sessionUser.Username {
		// check if the username exists
		exist, err := database.CheckExistance("User", "username", username)
		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}

		if exist {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		newUserName = username
	}

	// check if the email is valid and exists
	if !IsValidEmail(email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	if email != sessionUser.Email {
		exist, err := database.CheckExistance("User", "email", email)
		if err != nil {
			http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
			return
		}

		if exist {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	// Convert dateOfBirth from string to time.Time
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		http.Error(w, "Invalid date of birth format", http.StatusBadRequest)
		return
	}

	newUser := structs.User{
		Id:             sessionUser.Id,
		Username:       newUserName,
		Email:          email,
		DateOfBirth:    dob,
		FirstName:      firstName,
		LastName:       lastName,
		Country:        country,
		HashedPassword: sessionUser.HashedPassword,
		ImageId:        sessionUser.ImageId,
		Type:           sessionUser.Type,
		BannedUntil:    sessionUser.BannedUntil,
		GithubName:     "",
		LinkedinName:   "",
		TwitterName:    "",
		Bio:            sessionUser.Bio,
		Gender:         gender,
	}

	err = database.UpdateUserInfo(&newUser)
	if err != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ChatViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
	if sessionUser == nil {
		var chats []structs.Chats
		users, err := database.GetUsers()
		if err != nil {

			errorServer(w, r, http.StatusNotFound)
			return
		}
		for _, user := range users {
			chat := structs.Chats{
				UserId:   user.Id,
				Username: user.Username,
				Image:    imageIdToUrl(user.ImageId),
				Online:   IsUserOnline(user.Id),
			}
			chats = append(chats, chat)

		}
		sortedChats := SortChatsByOnline(chats)
		writeToJson(sortedChats, w)
	} else {
		users, err := database.GetUsers()
		if err != nil {
			errorServer(w, r, http.StatusNotFound)
			return
		}
		var chats []structs.Chats
		for _, user := range users {
			if user.Id != sessionUser.Id {
				chat := structs.Chats{
					UserId:   user.Id,
					Username: user.Username,
					Image:    imageIdToUrl(user.ImageId),
					Online:   IsUserOnline(user.Id),
				}
				chats = append(chats, chat)
			}
		}
		sortedChats := SortChatsByOnline(chats)
		writeToJson(sortedChats, w)
	}
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	sessionUser := GetUser(r)
	if sessionUser == nil {
		fmt.Print("helli")
		errorServer(w, r, http.StatusUnauthorized)
		return
	}
	userId := r.FormValue("id")
	IntId, _ := strconv.Atoi(userId)
	messages, err := database.GetMessages(sessionUser.Id, IntId)
	if err != nil {
		fmt.Println("Ggggggggggggg")
		errorServer(w, r, http.StatusNotFound)
		return
	}
	writeToJson(messages, w)
}
