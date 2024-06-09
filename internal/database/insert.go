package database

import (
	"sandbox/internal/structs"

	_ "github.com/mattn/go-sqlite3"
)

// Append New User Info the database, if any error occurs while appending (preparing and executing the SQL statements) it will return an error
func CreateUser(u structs.User) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement
	stmt, err := db.Prepare(`INSERT INTO User 
							(type_id, username, first_name, last_name, 
							date_of_birth, email, hashed_password, image_id, banned_until,
							github_name, linkedin_name, twitter_name) 
							VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the user's data
	_, err = stmt.Exec(
		u.Type,
		u.Username,
		u.FirstName,
		u.LastName,
		u.DateOfBirth,
		u.Email,
		u.HashedPassword,
		u.ImageId,
		u.BannedUntil,
		u.GithubName,
		u.LinkedinName,
		u.TwitterName)

	if err != nil {
		return err
	}

	return nil
}

// saves a session to the UserSession table, including the session token, user ID and creation time.
func AddSession(session structs.Session) error {
	err := RemoveSession(*session.UserId)
	if err != nil {
		return err
	}
	// Lock after remvoing the session to stop a deadlock
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to insert a new session
	stmt, err := db.Prepare("INSERT INTO UserSession (token, user_id, creation_time) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the new session
	_, err = stmt.Exec(session.Token, session.UserId, session.CreationTime)
	if err != nil {
		return err
	}

	return nil
}

// uploads the given image buffer to the database and returns the image ID
func UploadImage(imageBuffer []byte) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to insert the image data
	stmt, err := db.Prepare("INSERT INTO UploadedImage (data) VALUES (?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the image data
	result, err := stmt.Exec(imageBuffer)
	if err != nil {
		return 0, err
	}

	// Retrieve the generated ID of the inserted image
	imageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(imageID), nil
}

// associates a // The `post` function is adding a new post to the database. It prepares an SQL
// statement to insert post data into the `Post` table, including user ID, parent ID,
// title, message, image ID, time, like count, dislike count, love count, haha count,
// and skull count. After executing the SQL statement to insert the post data, it also
// calls the `addCategoriesToPost` function to associate the post with multiple
// category IDs by inserting corresponding rows into the `PostCategory` table.
// post with multiple category IDs in the PostCategory table.
// It takes a post ID and a slice of category IDs as parameters and inserts the corresponding rows.
func addCategoriesToPost(postID int, categoriesIDs []int) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement for inserting into PostCategory table
	stmt, err := tx.Prepare("INSERT INTO PostCategory (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Loop over the category IDs and execute the insert statement for each ID
	for _, categoryID := range categoriesIDs {
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// takes post struct and add it to the DB
func AddPost(post structs.Post) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to insert the image data
	stmt, err := db.Prepare(`INSERT INTO Post 
							(user_id, parent_id, title, message, image_id, time)
	 						VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement to insert the image data
	result, err := stmt.Exec(
		post.UserId,
		post.ParentId,
		post.Title,
		post.Message,
		post.ImageId,
		post.Time)
	if err != nil {
		return 0, err
	}

	// Retrieve the generated ID of the inserted post
	postId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = addCategoriesToPost(int(postId), post.CategoriesIDs)

	if err != nil {
		return 0, err
	}

	return int(postId), nil
}

func AddReactionToPost(reactionPost structs.PostReaction) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Remove any existing reaction by the user on the post
	_, err = tx.Exec("DELETE FROM PostReaction WHERE post_id = ? AND user_id = ?", reactionPost.PostId, reactionPost.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert the new reaction post into the PostReaction table
	stmt, err := tx.Prepare("INSERT INTO PostReaction (post_id, user_id, reaction_id) VALUES (?, ?, ?);")
	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		reactionPost.PostId,
		reactionPost.UserId,
		reactionPost.ReactionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// stores a Report struct in the Report table of the database and takes a report as Report-struct
func AddReport(report structs.Report) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the INSERT statement
	stmt, err := db.Prepare(`INSERT INTO Report (reporter_user_id, reported_user_id, report_message, 
		reported_post_id, time, is_post_report, is_pending, report_response) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	// Execute the INSERT statement with the values from the Report struct
	_, err = stmt.Exec(
		report.ReporterId,
		report.ReportedId,
		report.Reason,
		report.PostId,
		report.Time,
		report.IsPostReport,
		report.IsPending,
		report.ReportResponse)

	if err != nil {
		return err
	}

	return nil
}

// adds a new category to the Category table in the database based on the provided structs.Category object
func AddCategory(category structs.Category) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to insert a new category into the Category table
	stmt, err := db.Prepare(`INSERT INTO Category (name, description, color) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(category.Name, category.Description, category.Color)
	if err != nil {
		return err
	}

	return nil
}

// adds a Request to the PromoteRequest table in the database.
func AddPromoteRequest(request structs.PromoteRequest) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO PromoteRequest (user_id, description, time, is_pending) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with the request's values
	_, err = stmt.Exec(request.UserId, request.Reason, request.Time, request.IsPending)
	if err != nil {
		return err
	}

	return nil
}

func AddNotification(notification structs.UserNotification) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO UserNotification (comment_id, PostReaction_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	result, err := stmt.Exec(notification.CommentID, notification.PostReactionID)
	if err != nil {
		return 0, err
	}

	// Retrieve the generated ID of the inserted notification
	NotificationId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(NotificationId), nil
}
