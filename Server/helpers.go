package Server

import (
	"RealTimeForum/database"
	"RealTimeForum/structs"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"slices"
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
	return "/uploads/" + strconv.Itoa(imageId)
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
			author.Username = user.Username
			author.FirstName = user.FirstName
			author.LastName = user.LastName
			author.Type = userTypeToResponse(user.Type)
			author.ImageURL = imageIdToUrl(user.ImageId)
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

func mapReactionsForPost(post *structs.Post, loggedUserId int) []structs.PostReactionResponse {
	reactionsResp := []structs.PostReactionResponse{}

	// like
	likeResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeLike, "👍")
	if likeResp != nil {
		reactionsResp = append(reactionsResp, *likeResp)
	}

	// skull
	skullResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeSkull, "💀")
	if skullResp != nil {
		reactionsResp = append(reactionsResp, *skullResp)
	}

	// dislike
	dislikeResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeDislike, "👎")
	if dislikeResp != nil {
		reactionsResp = append(reactionsResp, *dislikeResp)
	}

	// haha
	hahaResp := mapReactionForPost(post, loggedUserId, structs.PostReactionTypeHaha, "😂")
	if hahaResp != nil {
		reactionsResp = append(reactionsResp, *hahaResp)
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

	_, fh, err := r.FormFile("image")
	if err == nil {
		data.Image = *fh
	}
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
		log.Printf("mapReactionForPost: reaction id: %s\n", err.Error())
		return nil
	}

	reactersIds, err := database.GetReactionUsers(post.Id, reactionId)
	if err != nil {
		log.Printf("mapReactionForPost: reacters Ids: %s\n", err.Error())
		return nil
	}

	reactionResp.Count = len(reactersIds)
	reactionResp.DidReact = slices.Contains(reactersIds, loggedUserId)

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
