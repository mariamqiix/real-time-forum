package database

import (
	"RealTimeForum/structs"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func UpdateRepor(reportID int, atype, value string) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the UPDATE statement
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE Report SET %s = ? WHERE id = ?", atype))
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the UPDATE statement to mark the report as resolved
	_, err = stmt.Exec(value, reportID)
	if err != nil {
		return err
	}

	return nil
}

// mark the report as resolved -not pendding-
func UpdateReportStatus(reportID int) error {
	return UpdateRepor(reportID, "is_pending", "false")
}

// add response to the report
func UpdateReportResponse(reportID int, response string) error {
	return UpdateRepor(reportID, "report_response", response)
}

// updates the information of a user in the User table.
func UpdateUserInfo(newUserInfo *structs.User) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to update the user
	stmt, err := db.Prepare(`UPDATE User SET type_id = ?, username = ?, first_name = ?, last_name = ?,country = ?, 
	date_of_birth = ?, email = ?, hashed_password = ?, image_id = ?, banned_until = ?, github_name = ?, 
	linkedin_name = ?, twitter_name = ?,bio = ?,gender = ?  WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the user
	_, err = stmt.Exec(
		newUserInfo.Type,
		newUserInfo.Username,
		newUserInfo.FirstName,
		newUserInfo.LastName,
		newUserInfo.Country,
		newUserInfo.DateOfBirth,
		newUserInfo.Email,
		newUserInfo.HashedPassword,
		newUserInfo.ImageId,
		newUserInfo.BannedUntil,
		newUserInfo.GithubName,
		newUserInfo.LinkedinName,
		newUserInfo.TwitterName,
		newUserInfo.Bio,
		newUserInfo.Gender,
		newUserInfo.Id)
	if err != nil {
		return err
	}

	return nil
}

// updates the IsPending status of a request to false in the PromoteRequest table.
func UpdateRequestStatus(requestID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE PromoteRequest SET is_pending = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the IsPending status
	_, err = stmt.Exec(false, requestID)
	if err != nil {
		return err
	}

	return nil
}

// updates the IsPending status of all requests for a specific user to false in the PromoteRequest table
// and removes every UserNotification that has PromoteRequestID and UserId == userID.
// returns the ID of the last PromoteRequestID for the user.
func UpdateRequestStatusByUserID(userID int) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// Prepare the SQL statement to update the IsPending status
	stmt, err := tx.Prepare("UPDATE PromoteRequest SET is_pending = ? WHERE user_id = ?")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the IsPending status
	_, err = stmt.Exec(false, userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Prepare the SQL statement to delete the UserNotifications
	stmt, err = tx.Prepare("DELETE FROM UserNotification WHERE promote_request_id IN (SELECT id FROM PromoteRequest WHERE user_id = ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the UserNotifications
	_, err = stmt.Exec(userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Fetch the last PromoteRequestID for the user
	var lastPromoteRequestID int
	err = tx.QueryRow("SELECT id FROM PromoteRequest WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID).Scan(&lastPromoteRequestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return lastPromoteRequestID, nil
}

// updates the information of a post in the Post table.
func UpdatePostInfo(newpost *structs.Post) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to update the post
	stmt, err := db.Prepare(`UPDATE Post SET title = ?, message = ?, image_id = ? WHERE id = ?`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the post
	_, err = stmt.Exec(
		newpost.Title,
		newpost.Message,
		newpost.ImageId,
		newpost.Id)
	if err != nil {
		return err
	}

	return nil
}

func MarkNotificationAsRead(notificationID int) error {
	// Implement the logic to mark the notification as read in the database
	_, err := db.Exec("UPDATE UserNotification SET read = 1 WHERE id = ?", notificationID)
	return err
}

// UpdateUserType updates the type_id of a user in the User table.
func UpdateUserType(userID int, newTypeID int) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to update the user's type_id
	stmt, err := db.Prepare(`UPDATE User SET type_id = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the user's type_id
	_, err = stmt.Exec(newTypeID, userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserPassword updates the password for a given user.
func UpdateUserPassword(userID int, newPasswordHash string) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	query := `UPDATE User SET hashed_password = ? WHERE id = ?`
	result, err := db.Exec(query, newPasswordHash, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected, user not found")
	}

	return nil
}

// BanUser bans a user for a week and removes their session if they are logged in
func BanUser(userId int) error {
	// Calculate the ban end date (one week from now)
	banEndDate := time.Now().AddDate(0, 0, 7)

	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Update the user's BannedUntil field
	stmt, err := tx.Prepare("UPDATE User SET banned_until = ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(banEndDate, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove the user's session
	err = removeSession(tx, userId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// removeSession removes the user's session
func removeSession(tx *sql.Tx, userId int) error {
	stmt, err := tx.Prepare("DELETE FROM UserSession WHERE user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

// UpdateImage updates the image data in the UploadedImage table.
func UpdateImage(imageID int, imageData []byte) error {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to update the image data
	stmt, err := db.Prepare(`UPDATE UploadedImage SET data = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to update the image data
	_, err = stmt.Exec(imageData, imageID)
	if err != nil {
		return err
	}

	return nil
}
