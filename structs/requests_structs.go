package structs

import "mime/multipart"

type UserRequest struct {
	Username     string `json:"username"`   // must for signup and login
	Email        string `json:"email"`      // must for signup
	FirstName    string `json:"first_name"` // must for signup
	LastName     string `json:"last_name"`  // must for signup
	Password     string `json:"password"`   // must for signup and login
	Image        []byte `json:"image"`
	DateOfBirth  int64  `json:"date_of_birth"` // must for signup
	GithubName   string `json:"github_name"`
	LinkedinName string `json:"linkedin_name"`
	TwitterName  string `json:"twitter_name"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"Description"`
}

type PostReportRequest struct {
	PostID int    `json:"post_id"`
	Reason string `json:"reason"`
}

type UserReportRequest struct {
	Username int    `json:"username"`
	Reason   string `json:"reason"`
}

type PromoteUserRequest struct {
	Reason string `json:"reason"`
}

type AddPostRequest struct {
	Title      string               `json:"title"`
	Content    string               `json:"content"`
	Categories []string             `json:"categories"` // category names to select from
	Image      multipart.FileHeader `json:"image"`      // probbaly will display as banner
}
