package server

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/hasher"
	"sandbox/internal/helpers"
	"sandbox/internal/sessionmanager"

	"sandbox/internal/structs"
	"slices"
	"time"
)

func filterPostsByReaction(posts []structs.Post, userId int, reactionType structs.PostReactionType) []structs.Post {
	filteredPosts := []structs.Post{}

	reactionId, err := database.GetReactionId(reactionType)
	if err != nil || reactionId == -1 {
		return filteredPosts
	}

	for _, post := range posts {
		usersIds, _ := database.GetReactionUsers(post.Id, reactionId)
		if slices.Contains(usersIds, userId) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

// this function takes a UserTypeId struct and returns a UserTypeResponse struct
// this function help us to see what is the type of the user (guest, user, moderator, admin) so we can know what's its permissions
func userTypeToResponse(userType structs.UserTypeId) structs.UserTypeResponse {
	switch userType {
	case structs.UserTypeIdGuest:
		return structs.UserTypeResponse{
			Name:        "Guest",
			Description: "A user that is not logged in",
			Color:       "#000000",
		}
	case structs.UserTypeIdUser:
		return structs.UserTypeResponse{
			Name:        "User",
			Description: "A user that is logged in",
			Color:       "#0000FF",
		}
	case structs.UserTypeIdModerator:
		return structs.UserTypeResponse{
			Name:        "Moderator",
			Description: "A user that is logged in and has moderation rights",
			Color:       "#00FF00",
		}
	case structs.UserTypeIdAdmin:
		return structs.UserTypeResponse{
			Name:        "Admin",
			Description: "A user that is logged in and has admin rights",
			Color:       "#FF0000",
		}
	}
	return structs.UserTypeResponse{}
}

// this function is used to get the user's profile and display it and serve the profile.html template
// this function use a struct type of profileView to pass the data to the template
// since this function will display the user's profile, we need to get the user's data from the database aka number of posts, username, first name, last name, image, and type
func profileHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	// get the user
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

	switch r.URL.Query().Get("q") {
	case "comments":
		// get comments only
		posts, err := database.GetPostsByUser(dbUser.Id, -1, 0, true)
		if err == nil {
			view.Posts = mapPosts(posts, sessionUser.Id)
		}

	case "likes":
		// get likes only
		// database.
		posts, err := database.GetUserReactions(dbUser.Id)
		if err != nil {
			log.Printf("profileHandler: %s\n", err.Error())
		} else {
			filterdPosts := filterPostsByReaction(posts, dbUser.Id, structs.PostReactionTypeLike)
			view.Posts = mapPosts(filterdPosts, sessionUser.Id)
		}

	case "dislikes":
		// get dislikes only
		posts, err := database.GetUserReactions(dbUser.Id)
		if err != nil {
			log.Printf("profileHandler: %s\n", err.Error())
		} else {
			filterdPosts := filterPostsByReaction(posts, dbUser.Id, structs.PostReactionTypeDislike)
			view.Posts = mapPosts(filterdPosts, sessionUser.Id)
		}
	default:
		posts, err := database.GetPostsByUser(dbUser.Id, -1, 0, false)
		if err == nil {
			view.Posts = mapPosts(posts, sessionUser.Id)
		}
	}

	err = templates.ExecuteTemplate(w, "profile.html", view)
	if err != nil {
		log.Printf("profileHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

// signupGetHandler handles the GET request for signup
func signupGetHandler(w http.ResponseWriter, r *http.Request) {
	// handle the case where if logged in, redirect to the homepage
	// we need to implement the session handling first
	err := templates.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		log.Printf("signupGetHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

// signupPostHandler handles the POST request for signup
// this function is used to create a new user and add it to the database and check if the username and email are already exists in the database as well checking the password length
// this function will handel the image upload as well so if the user wants to put as img or not in the profile
func signupPostHandler(w http.ResponseWriter, r *http.Request) {
	var userData structs.UserRequest
	if !helpers.ParseBody(&userData, r) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	// check password length
	if len(userData.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// check if the username exists
	exist, err := database.CheckExistance("User", "username", userData.Username)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}

	if exist {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// check if the email is valid and exists
	if !helpers.IsValidEmail(userData.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	exist, err = database.CheckExistance("User", "email", userData.Email)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}

	if exist {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	imageID := 0
	if userData.Image != nil {
		isImage, _ := helpers.IsDataImage(userData.Image)
		if isImage {
			imageID, err = database.UploadImage(userData.Image)
			if err != nil {
				log.Printf("SignupHandler: %s\n", err.Error())
			}
		}
	}
	// structure
	hashedPassword, hashErr := hasher.GetHash(userData.Password)
	if hashErr != hasher.HasherErrorNone {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}
	cleanedUserData := structs.User{
		Username:       userData.Username,
		Email:          userData.Email,
		FirstName:      userData.FirstName,
		LastName:       userData.LastName,
		DateOfBirth:    time.Unix(userData.DateOfBirth, 0),
		HashedPassword: hashedPassword,
		ImageId:        imageID,
		GithubName:     userData.GithubName,
		LinkedinName:   userData.LinkedinName,
		TwitterName:    userData.TwitterName,
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
	err = sessionmanager.CreateSessionAndSetCookie("", w, finalUser)
	if err != nil {
		http.Error(w, "something went wrong, please try again later", http.StatusInternalServerError)
		return
	}
}

// loginGetHandler handles the GET request for login
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		log.Printf("loginGetHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

// loginPostHandler handles the POST request for login
// this function is used to check if the user is already exists in the database and if the password is correct and if the user is banned or not
// if the user is not banned and the password is correct, the function will create a new session and add it to the database and create a new cookie for the user
// if the user is banned, the function will return a 403 forbidden error
func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 1024 {
		http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Parse the request body
	var loginData structs.UserRequest
	if !helpers.ParseBody(&loginData, r) {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user *structs.User
	var err error

	// check if the email is valid and exists
	if helpers.IsValidEmail(loginData.Username) {
		exist, err := database.CheckExistance("User", "email", loginData.Username)
		if err != nil {
			log.Printf("loginPostHandler: %s\n", err.Error())
			http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
			return
		}

		if !exist {
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

	if err := hasher.CompareHashAndPassword(user.HashedPassword, loginData.Password); err != hasher.HasherErrorNone {
		http.Error(w, "Invalid password", http.StatusConflict)
		return
	}

	if !user.BannedUntil.IsZero() && user.BannedUntil.After(time.Now()) {
		http.Error(w, "User is blocked", http.StatusForbidden)
		return
	}

	// Create a new session and set the cookie
	err2 := sessionmanager.CreateSessionAndSetCookie("", w, user)
	if err2 != nil {
		http.Error(w, "something went wromg, plaese try again later", http.StatusInternalServerError)
		return
	}
}

// logoutHandler handles the request for logout
// this function is used to remove the user's session from the database and remove the cookie from the user's browser
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		// No session token, user is not logged in
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Get the session from the Sessions map
	sessionToken := cookie.Value
	// Delete the session?.
	TokenInfo, err1 := database.GetSession(sessionToken)
	if err1 != nil {
		log.Printf("logoutHandler: %s\n", err1.Error())
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		return
	}
	err2 := database.RemoveSession(TokenInfo.Id)
	if err2 != nil {
		log.Printf("logoutHandler: %s\n", err2.Error())
		http.Error(w, "Something went wrong, contact server administrator", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Expires: time.Unix(0, 0),
	})
	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
