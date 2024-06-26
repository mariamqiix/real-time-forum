package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"log"
	"net/http"
	"os"
)

var userLimiter *UserRateLimiter

// GoLive starts the server and listens on the specified port and serves the routes
func GoLive(port string) {

	err := database.Connect("forum-db.sqlite")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	userLimiter = NewUserRateLimiter()
	
	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/uploads/{image_id}", uploadedContentServerHandler)
	http.HandleFunc("/uploads/add", uploadHandler)
	http.HandleFunc("/user/{user_id}", profileHandler)

	// http.HandleFunc("/signup", signupGetHandler)   // should be implemented in the js (get)
	http.HandleFunc("/signup", signupPostHandler)

	// http.HandleFunc("/login", loginGetHandler)   // should be implemented in the js (get)
	http.HandleFunc("/login", loginPostHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/category/{category_name}/", categoryPostsHandler)
	http.HandleFunc("/category", categoryGetHandler)

	http.HandleFunc("/post/{post_id}", postsHandler)
	http.HandleFunc("/post/{post_id}/delete", deletePostHandler)

	http.HandleFunc("/post/{post_id}/edit", editPostHandler)
	http.HandleFunc("/post/add", addPostHandler)

	// http.HandleFunc("GET /notifications/", notificationsHandler)
	// http.HandleFunc("POST /notifications/{notification_id}/read", markNotificationReadHandler)

	http.HandleFunc("/user/{user_id}/report", reportUserHandler)
	http.HandleFunc("/post/{post_id}/report", reportPostHandler)
	http.HandleFunc("/post/{post_id}/{reaction_type}", postReactionHandler)
	http.HandleFunc("/post/{post_id}/{reaction_type}/delete", deletePostReactionHandler)
	http.HandleFunc("/post/{post_id}/comment", addCommentGetHandler)
	http.HandleFunc("POST /post/{post_id}/comment", addCommentPostHandler)

	// http.HandleFunc("GET /login/google", authentication.HandleGoogleLogin)
	// http.HandleFunc("GET /login/github", authentication.HandleGitHubLogin)
	// http.HandleFunc("GET /login/facebook", authentication.HandleFacebookLogin)
	// http.HandleFunc("GET /login/google/callback", authentication.HandleGoogleCallback)
	// http.HandleFunc("GET /github-callback", authentication.HandleGitHubCallback)
	// http.HandleFunc("GET /facebook-callback", authentication.HandleFacebookCallback)

		http.ListenAndServe(":8080", nil)

}

// homepageHandler handles the homepage route and serves the homepage template
func homepageHandler(w http.ResponseWriter, r *http.Request) {
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

	dbPosts, err := database.GetPosts(-1, 0, "time DESC")
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

