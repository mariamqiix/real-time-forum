package database

// Remove reaction from Post by taking ReactionPost ID
func RemoveReactionFromPost(postId, reactionTypeId, reacterId int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Retrieve the reaction ID before deleting it
	var reactionId int64
	err = tx.QueryRow("SELECT id FROM PostReaction WHERE post_id = ? AND user_id = ? AND reaction_id = ?", postId, reacterId, reactionTypeId).Scan(&reactionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the reaction from the PostReaction table
	_, err = tx.Exec("DELETE FROM PostReaction WHERE id = ?", reactionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the corresponding notification from the UserNotification table
	_, err = tx.Exec("DELETE FROM UserNotification WHERE post_reaction_id = ?", reactionId)
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

func RemovePost(postID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete the notifications associated with the post reactions from the UserNotification table
	_, err = tx.Exec("DELETE FROM UserNotification WHERE post_reaction_id IN (SELECT id FROM PostReaction WHERE post_id = ?)", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the notifications associated with comments on the post from the UserNotification table
	_, err = tx.Exec("DELETE FROM UserNotification WHERE comment_id IN (SELECT id FROM Post WHERE parent_id = ?)", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the post reactions associated with the post from the PostReaction table
	_, err = tx.Exec("DELETE FROM PostReaction WHERE post_id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the post categories associated with the post from the PostCategory table
	_, err = tx.Exec("DELETE FROM PostCategory WHERE post_id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete post and the child posts of the post from the Post table
	_, err = tx.Exec("DELETE FROM Post WHERE id = ? OR parent_id = ?", postID, postID)
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

// delete session from UserSession table, by using session token
func RemoveSession(userId int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement to delete a session
	stmt, err := tx.Prepare("DELETE FROM UserSession WHERE user_id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the session
	_, err = stmt.Exec(userId)
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

// RemoveImage removes an image from the UploadedImage table by its ID.
func RemoveImage(imageID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement within the transaction
	stmt, err := tx.Prepare("DELETE FROM UploadedImage WHERE id = ?")
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the delete statement within the transaction
	_, err = stmt.Exec(imageID)
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	err = tx.Commit()
	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	return nil
}

func RemoveNotification(notificationID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	//Delete the notification
	_, err := db.Exec("DELETE FROM UserNotification WHERE id = ?", notificationID)
	if err != nil {
		return err
	}

	// Return no error if the deletion was successful
	return nil
}

// ReomvePromoteRequest removes all promotion requests with the given userId from the PromoteRequest table.
func ReomvePromoteRequest(userId int) error {

	mutex.Lock()
	defer mutex.Unlock()

	//Delete the notification
	_, err := db.Exec("DELETE FROM PromoteRequest WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	// Return no error if the deletion was successful
	return nil
}

// ReomvePromoteRequest removes all promotion requests with the given userId from the PromoteRequest table.
func RemoveCategory(Id int) error {

	mutex.Lock()
	defer mutex.Unlock()

	//Delete the notification
	_, err := db.Exec("DELETE FROM Category WHERE id = ?", Id)
	if err != nil {
		return err
	}

	// Return no error if the deletion was successful
	return nil
}

// RemovePostCategory removes all entries from the PostCategory table with the given categoryId.
func RemovePostCategory(categoryId int) error {

	mutex.Lock()
	defer mutex.Unlock()

	// Delete the post categories associated with the category ID from the PostCategory table
	_, err := db.Exec("DELETE FROM PostCategory WHERE category_id = ?", categoryId)
	if err != nil {
		return err
	}

	// Return no error if the deletion was successful
	return nil
}

// RemoveUser removes a user by their ID
func RemoveUser(userID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the SQL statement to delete a user
	stmt, err := tx.Prepare("DELETE FROM User WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the user
	_, err = stmt.Exec(userID)
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
