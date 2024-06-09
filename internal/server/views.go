package server

import (
	"html/template"
	"log"
	"sandbox/internal/structs"
)

var (
	templateRoot = "www/template/"
	templates    *template.Template
)

func init() {
	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// go func() {
	// 	for {
	// 		select {
	// 		case event, ok := <-watcher.Events:
	// 			if !ok {
	// 				return
	// 			}
	// 			if event.Has(fsnotify.Write) {
	// 				log.Println("updating templates")
	// 				updateTemplates()
	// 			}
	// 		case err, ok := <-watcher.Errors:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("error watching file:", err)
	// 		}
	// 	}
	// }()

	// if err := filepath.Walk(templateRoot, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !info.IsDir() {
	// 		err = watcher.Add(path)
	// 		if err != nil {
	// 			log.Println("error watching file:", err)
	// 		}
	// 	}
	// 	return nil
	// }); err != nil {
	// 	log.Println("error walking file:", err)
	// }
	updateTemplates()
}

func updateTemplates() {
	tmpl := template.New("")

	// Register custom functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	})

	// Parse template files
	tmpl, err := tmpl.ParseFiles(
		templateRoot+"index.html",
		templateRoot+"login.html",
		templateRoot+"signup.html",
		templateRoot+"profile.html",
		templateRoot+"newPost.html",
		templateRoot+"discussion.html",
		templateRoot+"error.html",
		templateRoot+"editPost.html",
		templateRoot+"views/notifications.html",
		templateRoot+"views/navbar.html",
		templateRoot+"views/footer.html",
		templateRoot+"views/sidebar.html",
		templateRoot+"views/posts.html",
		templateRoot+"views/discussions.html")

	if err != nil {
		log.Println("error updating templates:", err)
		return
	}

	// Set the updated templates
	templates = tmpl
}

// home view is the category page just the path is changed
type homeView struct {
	Posts       []structs.PostResponse
	User        *structs.UserResponse // nil if not logged in
	Categories  []structs.CategoryResponse
	SortOptions []string // index 0 is active sort
	Notifications []structs.NotificationResponse // Add this line
}

type profileView struct {
	User        *structs.UserResponse // nil if not logged in
	UserProfile structs.UserResponse
	Posts       []structs.PostResponse
	Comments    []structs.PostResponse
	Reactions   []structs.PostResponse
	Notifications []structs.NotificationResponse // Add this line
}

type addPostView struct {
	User       *structs.UserResponse      // nil if not logged in
	Categories []structs.CategoryResponse // category names to select from
	ParentId   int                        // init with -1 if not a comment
	OriginalId int                        // used to edit a post
	Notifications []structs.NotificationResponse // Add this line
}

type discussionView struct {
	User     *structs.UserResponse // nil if not logged in
	Post     structs.PostResponse
	Comments []structs.PostResponse
	Notifications []structs.NotificationResponse // Add this line
}
type errorView struct {
	User    *structs.UserResponse // nil if not logged in
	Message string
	Notifications []structs.NotificationResponse // Add this line
}
