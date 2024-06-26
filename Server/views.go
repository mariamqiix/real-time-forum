package Server

import (
	"RealTimeForum/structs"
	"html/template"
)

var (
	templateRoot = "www/template/"
	templates    *template.Template
)


// home view is the category page just the path is changed
type homeView struct {
	Posts         []structs.PostResponse
	User          *structs.UserResponse // nil if not logged in
	Categories    []structs.CategoryResponse
	SortOptions   []string                       // index 0 is active sort
	Notifications []structs.NotificationResponse // Add this line
}

type profileView struct {
	User          *structs.UserResponse // nil if not logged in
	UserProfile   structs.UserResponse
	Posts         []structs.PostResponse
	Comments      []structs.PostResponse
	Reactions     []structs.PostResponse
	Notifications []structs.NotificationResponse // Add this line
}

type addPostView struct {
	User          *structs.UserResponse          // nil if not logged in
	Categories    []structs.CategoryResponse     // category names to select from
	ParentId      int                            // init with -1 if not a comment
	OriginalId    int                            // used to edit a post
	Notifications []structs.NotificationResponse // Add this line
}

type discussionView struct {
	User          *structs.UserResponse // nil if not logged in
	Post          structs.PostResponse
	Comments      []structs.PostResponse
	Notifications []structs.NotificationResponse // Add this line
}
type errorView struct {
	User          *structs.UserResponse // nil if not logged in
	Message       string
	Notifications []structs.NotificationResponse // Add this line
}
