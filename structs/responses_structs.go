package structs

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
	Username  string           `json:"username"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	ImageURL  string           `json:"image_url"`
	Type      UserTypeResponse `json:"type"`
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
	Id    int    `json:"id"` // passed to `/notifications/{notification_id}/read` to mark as read
	Title string `json:"title"`
	Body  string `json:"body"`
	Link  string `json:"link"` // link to be opened when the notification is clicked
}
