package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"fmt"
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

	// Fetch the user from the database
	dbUser, err := database.GetUserByUsername(r.PathValue("user_id"))
	if err != nil {
		log.Printf("profileHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	if dbUser == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

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

	// Fetch the posts based on the query parameter
	switch r.URL.Query().Get("q") {
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
	dob ,_:= strconv.Atoi(r.FormValue("dob"))
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
	cleanedUserData := structs.User{
		Username:       username,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		DateOfBirth:    time.Unix(int64(dob), 0),
		HashedPassword: hashedPassword,
		ImageId:        imageID,
		GithubName:     "",
		LinkedinName:   "",
		TwitterName:    "",
	}

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
	fmt.Print("done")
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

	// Parse the request body
	var loginData structs.UserRequest
	if !ParseBody(&loginData, r) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user *structs.User
	var err error

	// check if the email is valid and exists
	if IsValidEmail(loginData.Username) {
		exist, err := database.CheckExistance("User", "email", loginData.Username)
		if err != nil {
			log.Printf("loginPostHandler: %s\n", err.Error())
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		} else if !exist {
			http.Error(w, "Invalid username or email or password", http.StatusConflict)
			return
		}

		user, err = database.GetUserByEmail(loginData.Username)
		if err != nil {
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		}

	} else {
		user, err = database.GetUserByUsername(loginData.Username)
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

	if err := CompareHashAndPassword(user.HashedPassword, loginData.Password); err != HasherErrorNone {
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get the session token from the cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Get the session from the database
	sessionToken := cookie.Value
	session, err := database.GetSession(sessionToken)
	if err != nil {
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		return
	}

	// Remove the session from the database
	if err := database.RemoveSession(session.Id); err != nil {
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Expires: time.Unix(0, 0),
	})

	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)

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
if r.Method == "GET" {sessionUser := GetUser(r)
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



func addPostHandler(w http.ResponseWriter, r *http.Request) {
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

	if sessionUser == 	nil {
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
} else if r.Method == "POST"  {

	sessionUser := GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}

	title := r.FormValue("title")
	Content := r.FormValue("topic")
	Categories := r.FormValue("selectedCategories")

	if !userLimiter.Allow(limiterUsername) {
		fmt.Print("hello")
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
	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("addPostPostHandler: %s\n", err.Error())
		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		return
	}

	for i, catName := range Categories {
		id := getCategoryIdFromArr(dbCategories, catName)
		if id == -1 {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}
		categoriesIds[i] = id
	}

	dbPostAdd := structs.Post{
		UserId:        sessionUser.Id,
		Title:         Title,
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

	newPostId, err := database.AddPost(dbPostAdd)
	if err != nil {
		log.Printf("addPostPostHandler: %s\n", err.Error())
		http.Error(w, "Internal server error, try again later", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.Itoa(newPostId), http.StatusSeeOther)
}
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

	// fmt.Printf("PostId: %d, UserId: %d, ReactionId: %d\n", PostStructForRec.Id, GetUser(r).Id, reactionId)
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

	categories , err :=  database.GetCategories()
	if err != nil {
			errorServer(w, r, http.StatusNotFound)
			http.Error(w, "cannot get the categories", http.StatusNotFound)
			return
	}
	writeToJson(categories,w)
}