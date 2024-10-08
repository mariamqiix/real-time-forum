package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//     "crypto/sha256"
	// "fmt"
	// "io/ioutil"
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

	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static"+r.URL.Path)
	})
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/style.css")
	})
	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static"+r.URL.Path)
	})

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/homePageDataHuncler", homePageDataHuncler)
	http.HandleFunc("/posts/AddReactions", addReactionHandler)
	http.HandleFunc("/uploads/{image_id}", uploadedContentServerHandler)
	http.HandleFunc("/uploads/add", uploadHandler)
	http.HandleFunc("/userProfile", profileHandler)
	http.HandleFunc("/user", LoggedUserHandler)
	http.HandleFunc("/userType", UserTypeHandler)
	http.HandleFunc("/Moderator", ModeratorHandler)
	http.HandleFunc("/RemoveModerator", RemoveModeratorHandler)
	http.HandleFunc("/PromotionRequest", PromotionRequestHandler)
	http.HandleFunc("/getUserInfo", getUserInfoHandler)
	http.HandleFunc("/signup", signupPostHandler)
	http.HandleFunc("/login", loginPostHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/changePassword", changePasswordHandler)
	http.HandleFunc("/updateUserInfo", updateUserInfoHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/postsByCategories", categoryPostsHandler)
	http.HandleFunc("/category", categoryGetHandler)
	http.HandleFunc("/post/{post_id}", postsHandler)
	http.HandleFunc("/post/{post_id}/delete", deletePostHandler)
	http.HandleFunc("/post/{post_id}/edit", editPostHandler)
	http.HandleFunc("/post/add/Get", addPostHandlerGet)
	http.HandleFunc("/post/add/Post", addPostHandlerPost)
	http.HandleFunc("/search", searchPostHandler)
	http.HandleFunc("/checkUserOnline", checkUserOnlineHandler)
	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/PromotionRequests", PromotionRequestsHandler)
	http.HandleFunc("/ShowUserPromotion", ShowUserPromotionHandler)
	http.HandleFunc("/RejectPromotion", RejectPromotionHandler)
	http.HandleFunc("/ApprovePromotion", ApprovePromotionHandler)
	http.HandleFunc("/removeCategory", removeCategoryHandler)
	http.HandleFunc("/addCategory", addCategoryHandler)
	http.HandleFunc("/ChatView", ChatViewHandler)
	http.HandleFunc("/Reports", ReportsHandler)
	http.HandleFunc("/banUser", BanUserHandler)
	http.HandleFunc("/updateReport", updateReportHandler)
	http.HandleFunc("/ReportsByUser", ReportsByUserHandler)
	http.HandleFunc("/notifications", notificationsHandler)
	http.HandleFunc("/user/{user_id}/report", reportUserHandler)
	http.HandleFunc("/post/{post_id}/report", reportPostHandler)
	http.HandleFunc("/post/reaction/delete", deletePostReactionHandler)
	http.HandleFunc("/post/{post_id}/comment", addCommentGetHandler)
	http.HandleFunc("/post/comment", addCommentPostHandler)
	http.HandleFunc("/userMessage", userMessageHandler)
	
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// homepageHandler handles the homepage route and serves the homepage template
func homepageHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "static/home.html"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	hash := sha256.Sum256(content)
	scc := fmt.Sprintf("%x", hash)

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-SCC", scc)
	http.ServeFile(w, r, filePath)

}

func homePageDataHuncler(w http.ResponseWriter, r *http.Request) {
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
			Id:          sessionUser.Id,
			Username:    sessionUser.Username,
			FirstName:   sessionUser.FirstName,
			LastName:    sessionUser.LastName,
			DateOfBirth: sessionUser.DateOfBirth,
			Location:    sessionUser.Country,
			ImageURL:    GetImageData(sessionUser.ImageId),
			Type:        userTypeToResponse(sessionUser.Type),
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

// staticHandler serves the static files aka the frontend files
func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "static" + r.URL.Path
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("staticHandler: %s\n", err.Error())
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Printf("staticHandler: %s\n", err.Error())
		http.NotFound(w, r)
		return
	}
	if fi.IsDir() {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), file)
}
