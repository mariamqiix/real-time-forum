package structs

import "time"

// MARK: Session
// Might not be needed
// type SessionResponse struct {
// 	Token  string
// 	Expiry int
// }

// MARK: Category
type CategoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IconURL     string `json:"icon_url"`
}

type MessageResponse struct {
	Type       string    `json:"type"`
	SenderId   int       `json:"SenderId"`
	ReceiverId int       `json:"ReceiverId"`
	Messag     string    `json:"Messag"`
	Time       time.Time `json:"Time"`
}

// MARK: Reaction
type PostReactionResponse struct {
	Reaction   string           `json:"reaction"`
	Type       PostReactionType `json:"type"`
	Count      int              `json:"count"`
	DidReact   bool             `json:"did_react"`
	WhoReacted []string         `json:"who_reacted"`
}

// MARK: Post
type PostResponse struct {
	Id         int                    `json:"id"`
	Author     UserResponse           `json:"author"`
	ParentId   int                    `json:"parent_id"`
	Title      string                 `json:"title"`
	Message    string                 `json:"message"`
	ImageURL   string                 `json:"image_url"`
	Categories []CategoryResponse     `json:"categories"`
	Reactions  []PostReactionResponse `json:"reactions"`
	CreatedAt  string                 `json:"created_at"`
}

// MARK: User
type UserTypeResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UserResponse struct {
	Id          int              `json:"id"`
	Username    string           `json:"username"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	DateOfBirth time.Time        `json:"DateOfBirth"`
	Location    string           `json:"location"`
	ImageURL    string           `json:"image_url"`
	Type        UserTypeResponse `json:"type"`
}

// Badges: github, twitter ...
type UserBadgeResponse struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
	Link    string `json:"link"`
}

type UserProfileResponse struct {
	User   UserResponse        `json:"user"`
	Badges []UserBadgeResponse `json:"badges"`
}

type NotificationResponse struct {
	IsReact              bool                       `json:"is_react"`
	IsComment            bool                       `json:"is_comment"`
	IsReport             bool                       `json:"is_report"`
	IsPromoteRequest     bool                       `json:"is_promote_request"`
	ReactionNotifi       ReactionNotification       `json:"ReactionNotifi"`
	PromoteRequestNotifi PromoteRequestNotification `json:"PromoteRequestNotification"`
	CommentNotifi        CommentNotification        `json:"CommentNotifi"`
	ReportRequestNotifi  ReportRequestNotification  `json:"ReportRequestNotifi"`
}

type CommentNotification struct {
	ParentId int          `json:"parent_id"`
	Username string       `json:"username"`
	Post     PostResponse `json:"PostResponse"`
}

type ReactionNotification struct {
	PostId   int          `json:"post_id"`
	Username string       `json:"username"`
	Post     PostResponse `json:"PostResponse"`
	Reaction string       `json:"reaction"`
}
type ReportRequestNotification struct {
	Reason            string `json:"reason"`
	Accepted          bool   `json:"accepted"`
	ReportedUsername  string `json:"reported_username"`
	ReportedPostTitle string `json:"reported_post_id"`
}

type PromoteRequestNotification struct {
	Reason   string `json:"reason"`
	Accepted bool   `json:"accepted"`
}

type PromoteRequestResponse struct {
	Id        int
	UserId    int
	Username  string
	Reason    string
	IsPending bool
}

type ReportRequestResponse struct {
	Id                int       `json:"id"`
	ReporterId        int       `json:"reporter_id"`
	ReporterUsername  string    `json:"reporter_username"`
	ReportedId        int       `json:"reported_id"`
	ReportedUsername  string    `json:"reported_username"`
	ReportedPostId    int       `json:"reported_post_id"`
	ReportedPostTitle string    `json:"reported_post_title"`
	Time              time.Time `json:"time"`
	Reason            string    `json:"reason"`
	IsPostReported    bool      `json:"is_post_reported"`
	IsPending         bool      `json:"is_pending"`
	ReportResponse    string    `json:"report_response"`
}
