package database

// retrieves the count of posts
func GetPostsCount() (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to get the count of posts
	stmt, err := db.Prepare(`SELECT COUNT(*) 
                             FROM Post
							 WHERE parent_id IS NULL`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// retrieves the count of User posts
func GetPostsCountByUser(userId int) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to get the count of posts
	stmt, err := db.Prepare(`SELECT COUNT(*) 
                             FROM Post
							 WHERE user_id = ?
							 AND parent_id IS NULL`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// retrieves the count of posts associated with a given category name by joining the post_category with the category tables, based on the provided category name.
func GetPostsCountByCategory(categoryName string) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to get the count of posts for the given category name
	stmt, err := db.Prepare(`SELECT COUNT(*) 
                             FROM PostCategory 
                             INNER JOIN Category ON PostCategory.category_id = Category.id
                             WHERE Category.name = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(categoryName).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// retrieves the count of post comments
func GetCommentsCountForPost(parentID string) (int, error) {
	// Lock the mutex before accessing the database
	mutex.Lock()
	defer mutex.Unlock()

	// Prepare the SQL statement to get the count of posts comments
	stmt, err := db.Prepare(`SELECT COUNT(*) 
                             FROM Post 
                             WHERE parent_id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(parentID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// retrieves the count of reactions of type 1 and type 2 for a given post
func GetReactionCountsForPost(postId int) (int, int, error) {
    // Lock the mutex before accessing the database
    mutex.Lock()
    defer mutex.Unlock()

    // Prepare the SQL statement to get the count of reactions of type 1 for the given post
    stmt1, err := db.Prepare(`SELECT COUNT(*) 
                              FROM PostReaction 
                              WHERE post_id = ? 
                              AND reaction_id = 1`)
    if err != nil {
        return 0, 0, err
    }
    defer stmt1.Close()

    var count1 int
    err = stmt1.QueryRow(postId).Scan(&count1)
    if err != nil {
        return 0, 0, err
    }

    // Prepare the SQL statement to get the count of reactions of type 2 for the given post
    stmt2, err := db.Prepare(`SELECT COUNT(*) 
                              FROM PostReaction 
                              WHERE post_id = ? 
                              AND reaction_id = 2`)
    if err != nil {
        return 0, 0, err
    }
    defer stmt2.Close()

    var count2 int
    err = stmt2.QueryRow(postId).Scan(&count2)
    if err != nil {
        return 0, 0, err
    }

    return count1, count2, nil
}