package server

import (
	"log"
	"net/http"
	"os"
	"sandbox/internal/authentication"
	"sandbox/internal/cert"
	"sandbox/internal/database"
	"sandbox/internal/ratelimiter"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"time"
)

var userLimiter *ratelimiter.UserRateLimiter

// GoLive starts the server and listens on the specified port and serves the routes
func GoLive(port string) {
	err := database.Connect("forum-db.sqlite")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	userLimiter = ratelimiter.NewUserRateLimiter()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", homepageHandler)                                // ✅
	mux.HandleFunc("GET /static/", staticHandler)                           // ✅
	mux.HandleFunc("GET /uploads/{image_id}", uploadedContentServerHandler) // ✅
	mux.HandleFunc("POST /uploads/add", uploadHandler)                      // ✅
	// profile route
	mux.HandleFunc("GET /user/{user_id}", profileHandler) // ✅
	// user routes
	mux.HandleFunc("GET /signup", signupGetHandler)   // ✅
	mux.HandleFunc("POST /signup", signupPostHandler) // ✅
	mux.HandleFunc("GET /login", loginGetHandler)     // ✅
	mux.HandleFunc("POST /login", loginPostHandler)   // ✅
	mux.HandleFunc("GET /logout", logoutHandler)      // ✅

	// category routes
	mux.HandleFunc("GET /category/{category_name}/", categoryPostsHandler) // ✅
	//  posts
	mux.HandleFunc("GET /post/{post_id}", postsHandler)

	mux.HandleFunc("GET /post/{post_id}/delete", deletePostHandler) // ✅
	mux.HandleFunc("GET /post/{post_id}/edit", editPostGetHandler)  // ✅
	mux.HandleFunc("POST /post/{post_id}/edit", editPostPostHandler)

	// posting
	mux.HandleFunc("GET /post/add", addPostGetHandler)   // ✅
	mux.HandleFunc("POST /post/add", addPostPostHandler) // ✅

	// notifications
	mux.HandleFunc("GET /notifications/", notificationsHandler)
	mux.HandleFunc("POST /notifications/{notification_id}/read", markNotificationReadHandler)

	// reports
	mux.HandleFunc("POST /user/{user_id}/report", reportUserHandler) // ✅
	mux.HandleFunc("POST /post/{post_id}/report", reportPostHandler) // ✅

	// reactions
	mux.HandleFunc("POST /post/{post_id}/{reaction_type}", postReactionHandler)         // ✅
	mux.HandleFunc("DELETE /post/{post_id}/{reaction_type}", deletePostReactionHandler) // ✅

	// comments
	mux.HandleFunc("GET /post/{post_id}/comment", addCommentGetHandler)
	mux.HandleFunc("POST /post/{post_id}/comment", addCommentPostHandler)

	// OAuth
	mux.HandleFunc("GET /login/google", authentication.HandleGoogleLogin)
	mux.HandleFunc("GET /login/github", authentication.HandleGitHubLogin)
	mux.HandleFunc("GET /login/facebook", authentication.HandleFacebookLogin)
	mux.HandleFunc("GET /login/google/callback", authentication.HandleGoogleCallback)
	mux.HandleFunc("GET /github-callback", authentication.HandleGitHubCallback)
	mux.HandleFunc("GET /facebook-callback", authentication.HandleFacebookCallback)

	// Generate TLS certificate and key
	certFilePath := "cert.pem"
	keyFilePath := "key.pem"
	err = cert.GenerateSelfSignedCert(certFilePath, keyFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Listening on https://localhost:%s\n", port)
	server := http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}
	if err := server.ListenAndServeTLS(certFilePath, keyFilePath); err != nil {
		log.Fatal(err.Error())
	}
}

// homepageHandler handles the homepage route and serves the homepage template
func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if r.URL.Path == "/favicon.ico" {
			http.ServeFile(w, r, "www/static/imgs/logo.png")
			return
		}
		errorServer(w, r, http.StatusNotFound)
		return
	}

	sessionUser := sessionmanager.GetUser(r)
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

	err = templates.ExecuteTemplate(w, "index.html", view)
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}

// staticHandler serves the static files aka the frontend files
func staticHandler(w http.ResponseWriter, r *http.Request) {
	if !isAcceptedMethod(w, r, http.MethodGet) {
		return
	}
	filePath := "www" + r.URL.Path
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("staticHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusNotFound)
		return
	}
	fi, err := file.Stat()
	if err != nil {
		log.Printf("staticHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusNotFound)
		return
	}
	if fi.IsDir() {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), file)
}
