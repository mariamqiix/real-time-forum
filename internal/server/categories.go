package server

import (
	"log"
	"net/http"
	"sandbox/internal/database"
	"sandbox/internal/sessionmanager"
	"sandbox/internal/structs"
)

func getCategoryIdFromArr(arr []structs.Category, name string) int {
	for _, category := range arr {
		if category.Name == name {
			return category.Id
		}
	}
	return -1
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

func categoryPostsHandler(w http.ResponseWriter, r *http.Request) {
	sessionUser := sessionmanager.GetUser(r)
	limiterUsername := "[GUESTS]"
	if sessionUser != nil {
		limiterUsername = sessionUser.Username
	}
	if !userLimiter.Allow(limiterUsername) {
		errorServer(w, r, http.StatusTooManyRequests)
		return
	}

	categoryName := r.PathValue("category_name")

	// check if category exists
	categoryExists, _ := database.CheckExistance("Category", "name", categoryName)
	if !categoryExists {
		errorServer(w, r, http.StatusNotFound)
		return
	}

	posts_count, err := database.GetPostsCountByCategory(categoryName)
	if err != nil {
		log.Printf("error getting posts count by category: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

	dbPosts, err := database.GetPostsByCategory(categoryName, posts_count, 0, "latest")
	if err != nil {
		log.Printf("error getting posts by category: %s\n", err.Error())
	}

	dbCategories, err := database.GetCategories()
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
	}

	mappedPosts := mapPosts(dbPosts, -1)
	mappedCategories := mapCategories(dbCategories)

	view := homeView{
		Posts:       mappedPosts,
		User:        nil,
		Categories:  mappedCategories,
		SortOptions: []string{"latest", "most liked", "least liked", "oldestq"},
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

	err = templates.ExecuteTemplate(w, "index.html", view)
	if err != nil {
		log.Printf("homepageHandler: %s\n", err.Error())
		errorServer(w, r, http.StatusInternalServerError)
		return
	}

}
