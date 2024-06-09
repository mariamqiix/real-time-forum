package server

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"slices"
	"strconv"
)

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

// helper function to prepare the reations, if used not logged in, pass -1 as the loggedUserId
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

func postsHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	postId, err := strconv.Atoi(r.PathValue("post_id"))
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	post, err := database.GetPost(postId)
	if err != nil || post == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}
	if post.ParentId != nil {
		newPath := "/post/" + strconv.Itoa(*post.ParentId) + "#" + r.PathValue("post_id")
		http.Redirect(w, r, newPath, http.StatusFound)
		return
	}

	view := discussionView{
		User:     nil,
		Post:     structs.PostResponse{},
		Comments: nil,
	}

	comments, err := database.GetCommentsForPost(post.Id, -1, 0)
	if err != nil {
		log.Printf("postsHandler: %s\n", err.Error())
	}

	if sessionUser != nil {
		view.User = &structs.UserResponse{
			Username:  sessionUser.Username,
			FirstName: sessionUser.FirstName,
			LastName:  sessionUser.LastName,
			ImageURL:  imageIdToUrl(sessionUser.ImageId),
			Type:      userTypeToResponse(sessionUser.Type),
		}
		view.Post = mapPosts([]structs.Post{*post}, sessionUser.Id)[0]
		view.Comments = mapPosts(comments, sessionUser.Id)
	} else {
		view.Post = mapPosts([]structs.Post{*post}, -1)[0]
		view.Comments = mapPosts(comments, -1)

	}

	err = templates.ExecuteTemplate(w, "discussion.html", view)
	if err != nil {
		log.Printf("postsHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}
