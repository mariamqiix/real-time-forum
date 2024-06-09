package server

import (
	// "fmt"

	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
	"strconv"
)

// reaction type must be of type structs.PostReactionType
// WE WILL NOT ACCEPT ANY OTHER TYPE
// validate post and reaction type
// The user can't add the SAME REACTION twice
func postReactionHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	// get the post id and reaction type
	postId := r.PathValue("post_id")
	reactionType := r.PathValue("reaction_type")

	// validate the reaction type
	if reactionType != string(structs.PostReactionTypeLike) && reactionType != string(structs.PostReactionTypeDislike) && reactionType != string(structs.PostReactionTypeLove) && reactionType != string(structs.PostReactionTypeHaha) && reactionType != string(structs.PostReactionTypeSkull) {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// validate the post id and covert it to int
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// get the post struct from the database
	PostStructForRec, err := database.GetPost(postIdInt)
	if err != nil || PostStructForRec == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	// get the user from the session
	UsersPost := sessionmanager.GetUser(r)
	if UsersPost == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	mappedReaction := mapReactionForPost(PostStructForRec, UsersPost.Id, structs.PostReactionType(reactionType), reactionType)
	if mappedReaction == nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	reactionId, err2 := database.GetReactionId(mappedReaction.Type)
	if err2 != nil || reactionId == 0 {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	reaction := structs.PostReaction{
		PostId:     PostStructForRec.Id,
		UserId:     UsersPost.Id,
		ReactionId: reactionId,
	}

	err3 := database.AddReactionToPost(reaction)
	if err3 != nil {
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
	// add notification
	// if the user is not the owner of the post
	if UsersPost.Id != PostStructForRec.UserId {
		notification := structs.UserNotification{
			PostReactionID: reactionId,
		}
		_, err4 := database.AddNotification(notification)
		if err4 != nil {
			log.Printf("postReactionHandler: %s\n", err4.Error())
		}
	}
}

// same as postReactionHandler
// validate all
// on sucess return count for all reactions
// user can't delete another user reaction
// get the post id and reaction type
func deletePostReactionHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	postId := r.PathValue("post_id")
	reactionType := r.PathValue("reaction_type")

	// validate the reaction type
	if reactionType != string(structs.PostReactionTypeLike) && reactionType != string(structs.PostReactionTypeDislike) && reactionType != string(structs.PostReactionTypeLove) && reactionType != string(structs.PostReactionTypeHaha) && reactionType != string(structs.PostReactionTypeSkull) {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// validate the post id and covert it to int
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}

	// get the post struct from the database
	PostStructForRec, err := database.GetPost(postIdInt)
	if err != nil || PostStructForRec == nil {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	UsersPost := sessionmanager.GetUser(r)
	if UsersPost == nil {
		errorServer(w, r, http.StatusUnauthorized)
		return
	}

	mappedReaction := mapReactionForPost(PostStructForRec, UsersPost.Id, structs.PostReactionType(reactionType), reactionType)
	if mappedReaction == nil {
		errorServer(w, r, http.StatusBadRequest)
		return
	}
	reactionId, err2 := database.GetReactionId(mappedReaction.Type)
	if err2 != nil {
		log.Printf("deletePostReactionHandler: %s\n", err2.Error())

		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	// fmt.Printf("PostId: %d, UserId: %d, ReactionId: %d\n", PostStructForRec.Id, sessionmanager.GetUser(r).Id, reactionId)
	err3 := database.RemoveReactionFromPost(PostStructForRec.Id, reactionId, UsersPost.Id)
	if err3 != nil {
		log.Printf("deletePostReactionHandler: %s\n", err3.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}
}
