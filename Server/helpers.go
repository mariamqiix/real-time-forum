package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[\w]+@[\w]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsDataImage(buff []byte) (bool, string) {
	// the function that actually does the trick
	t := http.DetectContentType(buff)
	return strings.HasPrefix(t, "image"), t
}

// Returns false if an error happens
func ParseBody(data any, r *http.Request) bool {
	return json.NewDecoder(r.Body).Decode(&data) == nil
}

func imageIdToUrl(imageId int) string {
	return "/images/" + strconv.Itoa(imageId) + ".png"
}

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

func mapPosts(oldArr []structs.Post, loggedUserId int) []structs.PostResponse {
	newArr := make([]structs.PostResponse, len(oldArr))
	for i, old := range oldArr {
		user, err := database.GetUserById(old.UserId)

		author := structs.UserResponse{}
		// if the user is not found, set the author to unknown
		if err != nil {
			log.Printf("mapPosts: %s\n", err.Error())
			author.Username = "[Unknown]"
			author.Type = userTypeToResponse(structs.UserTypeIdUser)
		} else {
			author.Id = user.Id
			author.Username = user.Username
			author.FirstName = user.FirstName
			author.LastName = user.LastName
			author.DateOfBirth = user.DateOfBirth
			author.Location = user.Country
			author.Type = userTypeToResponse(user.Type)
			author.ImageURL = GetImageData(user.ImageId)
		}
		newArr[i] = structs.PostResponse{
			Id:         old.Id,
			Author:     author,
			Title:      old.Title,
			Message:    old.Message,
			ImageURL:   imageIdToUrl(old.ImageId),
			Categories: mapCategoriesForPost(old.CategoriesIDs),
			Reactions:  mapReactionsForPost(&old, loggedUserId),
			CreatedAt:  old.Time.UTC().Format("2006-01-02 15:04:05"),
		}
		if old.ParentId == nil {
			newArr[i].ParentId = -1
		} else {
			newArr[i].ParentId = *old.ParentId

		}
	}
	return newArr

}

func mapPost(old structs.Post, loggedUserId int) structs.PostResponse {
	user, err := database.GetUserById(old.UserId)

	author := structs.UserResponse{}
	// if the user is not found, set the author to unknown
	if err != nil {
		log.Printf("mapPost: %s\n", err.Error())
		author.Username = "[Unknown]"
		author.Type = userTypeToResponse(structs.UserTypeIdUser)
	} else {
		author.Id = user.Id
		author.Username = user.Username
		author.FirstName = user.FirstName
		author.LastName = user.LastName
		author.DateOfBirth = user.DateOfBirth
		author.Location = user.Country
		author.Type = userTypeToResponse(user.Type)
		author.ImageURL = GetImageData(user.ImageId)
	}

	postResponse := structs.PostResponse{
		Id:         old.Id,
		Author:     author,
		Title:      old.Title,
		Message:    old.Message,
		ImageURL:   imageIdToUrl(old.ImageId),
		Categories: mapCategoriesForPost(old.CategoriesIDs),
		Reactions:  mapReactionsForPost(&old, loggedUserId),
		CreatedAt:  old.Time.UTC().Format("2006-01-02 15:04:05"),
	}

	if old.ParentId == nil {
		postResponse.ParentId = -1
	} else {
		postResponse.ParentId = *old.ParentId
	}

	return postResponse
}

func mapReactionsForPost(post *structs.Post, loggedUserId int) []structs.PostReactionResponse {
	reactionsResp := []structs.PostReactionResponse{}

	// like
	likeResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeLike, "👍")
	if likeResp != nil {
		reactionsResp = append(reactionsResp, *likeResp)
	}

	// dislike
	dislikeResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeDislike, "👎")
	if dislikeResp != nil {
		reactionsResp = append(reactionsResp, *dislikeResp)
	}

	return reactionsResp
}

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

func mapCategories(oldArr []structs.Category) []structs.CategoryResponse {
	newArr := make([]structs.CategoryResponse, len(oldArr))
	for i, old := range oldArr {
		newArr[i] = structs.CategoryResponse{
			Name:        old.Name,
			Description: old.Description,
			Color:       old.Color,
			IconURL:     "",
		}
	}
	return newArr
}

func parsePostForm(data *structs.AddPostRequest, r *http.Request) bool {
	err := r.ParseMultipartForm(22 << 10)
	if err != nil {
		log.Printf("parseForm: %s\n", err.Error())
		return false
	}
	data.Title = r.FormValue("title")
	data.Content = r.FormValue("content")
	data.Categories = r.Form["categories"]

	// _, fh, err := r.FormFile("image")
	// if err == nil {
	// 	data.Image = *fh
	// }
	return true
}

func getCategoryIdFromArr(arr []structs.Category, name string) int {
	for _, category := range arr {
		if category.Name == name {
			return category.Id
		}
	}
	return -1
}

func writeToJson(v any, w http.ResponseWriter) error {
	buff, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(buff)
	return err
}

func mapReactionForPost(post *structs.Post, loggedUserId int, reactionType structs.PostReactionType, emoji string) *structs.PostReactionResponse {
	reactionResp := &structs.PostReactionResponse{
		Type:     reactionType,
		Reaction: emoji,
		DidReact: false,
	}

	reactionId, err := database.GetReactionId(reactionType)
	if err != nil {
		return nil
	}

	reactersIds, err := database.GetReactionUsers(post.Id, reactionId)
	if err != nil {
		return nil
	}

	reactionResp.Count = len(reactersIds)
	if loggedUserId != -1 {
		reactionResp.DidReact = slices.Contains(reactersIds, loggedUserId)
	}

	return reactionResp
}

func mapCategoriesForPost(categoriesIDs []int) []structs.CategoryResponse {
	respCategories := []structs.CategoryResponse{}
	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("mapCategoryForPost: %s\n", err.Error())
		return nil
	}
	for _, catId := range categoriesIDs {
		for _, dbCategory := range dbCategories {
			if dbCategory.Id == catId {
				respCategories = append(respCategories, structs.CategoryResponse{
					Name:        dbCategory.Name,
					Description: dbCategory.Description,
				})
			}
		}
	}
	return respCategories
}

func ConvertToPromoteRequestResponse(promotionRequests []structs.PromoteRequest) ([]structs.PromoteRequestResponse, error) {
	var responses []structs.PromoteRequestResponse

	for _, request := range promotionRequests {
		user, err := database.GetUserById(request.UserId)
		if err != nil {
			return nil, err
		}

		response := structs.PromoteRequestResponse{
			Id:        request.Id,
			UserId:    request.UserId,
			Username:  user.Username,
			Reason:    request.Reason,
			IsPending: request.IsPending,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// SortChatsByOnline sorts the chats slice so that online chats are at the top
func SortChatsByOnline(chats []structs.Chats) []structs.Chats {
	sort.SliceStable(chats, func(i, j int) bool {
		return chats[i].Online && !chats[j].Online
	})
	return chats
}

func ConvertToReportRequestResponse(reports []structs.Report) ([]structs.ReportRequestResponse, error) {
	var responses []structs.ReportRequestResponse

	for _, report := range reports {
		reporter, err := database.GetUserById(report.ReporterId)
		if err != nil {
			return nil, err
		}
		reported, err := database.GetUserById(report.ReportedId)
		if err != nil {
			return nil, err
		}

		var post *structs.Post
		if report.PostId != -1 && report.IsPostReport {
			post, err = database.GetPost(report.PostId)
			if err != nil {
				return nil, err
			}
		} else {
			report.PostId = 0
		}

		// Check if post is nil before accessing its fields
		var postTitle string
		if post != nil {
			postTitle = post.Title
		} else {
			postTitle = ""
		}
		response := structs.ReportRequestResponse{
			Id:                report.Id,
			ReporterId:        report.ReporterId,
			ReporterUsername:  reporter.Username,
			ReportedId:        report.ReportedId,
			ReportedUsername:  reported.Username,
			ReportedPostId:    report.PostId,
			ReportedPostTitle: postTitle, // Assuming you need to fetch the post title separately if required
			Time:              report.Time,
			Reason:            report.Reason,
			IsPostReported:    report.IsPostReport,
			IsPending:         report.IsPending,
			ReportResponse:    report.ReportResponse,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// / GetImageData retrieves the image data from the database and returns it as a base64-encoded string.
func GetImageData(imageID int) string {
	// Retrieve the image data from the database
	imageData, err := database.GetImage(imageID)
	if err != nil {
		return ""
	}

	if imageData == nil {
		return ""
	}

	// Encode the image data to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image
}

// ImageURLToBytes reads an image from the given file path and returns it as a byte slice.
func ImageURLToBytes(imagePath string) ([]byte, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the image data from the file
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func ConvertToNotificationResponse(notifications []structs.UserNotification) ([]structs.NotificationResponse, error) {
	var responses []structs.NotificationResponse
	for _, notification := range notifications {
		if notification.PostReactionID == 0 && notification.CommentID == 0 && notification.ReportID == 0 && notification.PromoteRequestID == 0 {
			continue
		}
		react := false
		comment := false
		report := false
		PromoteRequestID := false
		if notification.PostReactionID != 0 {
			react = true
		} else if notification.CommentID != 0 {
			comment = true
		} else if notification.ReportID != 0 {
			report = true
		} else if notification.PromoteRequestID != 0 {
			PromoteRequestID = true
		}
		if react {
			response := structs.NotificationResponse{
				IsReact:          react,
				IsComment:        comment,
				IsReport:         report,
				IsPromoteRequest: PromoteRequestID,
			}
			React, err := database.GetReactionById(notification.PostReactionID)
			if err != nil {
				return nil, err
			}
			if React == nil {
				return nil, fmt.Errorf("reaction not found for PostReactionID: %d", notification.PostReactionID)
			}
			reactionn := "like"
			if React.ReactionId == 2 {
				reactionn = "dislike"
			}
			username, err := database.GetUsernameByUserId(React.UserId)
			if err != nil {
				return nil, err
			}
			Post, err := database.GetPost(React.PostId)
			if err != nil {
				return nil, err
			}
			if Post == nil {
				return nil, fmt.Errorf("post not found for PostId: %d", React.PostId)
			}
			reactionNotification := structs.ReactionNotification{
				PostId:   React.PostId,
				Username: username,
				Post:     mapPost(*Post, notification.UserId),
				Reaction: reactionn,
			}
			response.ReactionNotifi = reactionNotification
			responses = append(responses, response)
		} else if comment {
			response := structs.NotificationResponse{
				IsReact:          react,
				IsComment:        comment,
				IsReport:         report,
				IsPromoteRequest: PromoteRequestID,
			}
			Post, err := database.GetPost(notification.CommentID)
			if err != nil {
				return nil, err
			}
			if Post == nil {
				return nil, fmt.Errorf("post not found for CommentID: %d", notification.CommentID)
			}

			var parentId int
			if Post.ParentId != nil {
				parentId = *Post.ParentId
			}
			ParentPost, err := database.GetPost(parentId)
			if err != nil {
				return nil, err
			}
			if ParentPost == nil {
				return nil, fmt.Errorf("post not found for CommentID: %d", notification.CommentID)
			}
			username, err := database.GetUsernameByUserId(Post.UserId)
			if err != nil {
				return nil, err
			}
			commentNotification := structs.CommentNotification{
				ParentId: *Post.ParentId,
				Username: username,
				Post:     mapPost(*ParentPost, notification.UserId),
			}
			response.CommentNotifi = commentNotification
			responses = append(responses, response)
		} else if PromoteRequestID {
			response := structs.NotificationResponse{
				IsReact:          react,
				IsComment:        comment,
				IsReport:         report,
				IsPromoteRequest: PromoteRequestID,
			}
			PromoteRequest, err := database.GetRequest(notification.PromoteRequestID)
			if err != nil {
				return nil, err
			}
			User, err := database.GetUserById(PromoteRequest.UserId)
			if err != nil {
				return nil, err
			}
			accepted := false
			if User.Type == 2 {
				accepted = true
			}
			promoteNotification := structs.PromoteRequestNotification{
				Reason:   PromoteRequest.Reason,
				Accepted: accepted,
			}
			response.PromoteRequestNotifi = promoteNotification
			responses = append(responses, response)
		} else if report {
			response := structs.NotificationResponse{
				IsReact:          react,
				IsComment:        comment,
				IsReport:         report,
				IsPromoteRequest: PromoteRequestID,
			}
			Report, err := database.GetReport(notification.ReportID)
			if err != nil {
				return nil, err
			}
			reported, err := database.GetUsernameByUserId(Report.ReportedId)
			if err != nil {
				return nil, err
			}
			accepted := false
			if Report.ReportResponse != "" && Report.ReportResponse != "Report Rejected" {
				fmt.Print((Report.ReportResponse))
				accepted = true
			}
			PostTitle := ""
			if Report.IsPostReport && Report.ReportResponse == "Report Rejected" {
				fmt.Print("hi")
				Post, err := database.GetPost(Report.PostId)
				if err != nil {
					return nil, err
				}
				if Post == nil {
					return nil, fmt.Errorf("post not found for PostId: %d", Report.PostId)
				}
				PostTitle = Post.Title
			}
			reportNotification := structs.ReportRequestNotification{
				Reason:            Report.Reason,
				Accepted:          accepted,
				ReportedUsername:  reported,
				ReportedPostTitle: PostTitle,
			}
			response.ReportRequestNotifi = reportNotification
			responses = append(responses, response)
		}
	}
	return responses, nil
}
