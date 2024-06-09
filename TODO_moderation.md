# Forum Moderation System To-Do List

## 1. Implement User Roles (files involved: `create_db.sql`, `internal/structs/structs.go`, `internal/database/get.go`, `internal/database/insert.go`, `internal/server/user.go`)
### 1.1. Update Database Schema
- **Objective**: Define the roles and permissions for different types of users.
- **Roles**:
  - **Guest**: Can only view posts and comments.
  - **User**: Can create, comment, like, and dislike posts.
  - **Moderator**: Can delete posts, report posts to the admin, and approve posts.
  - **Administrator**: Can promote/demote users, manage categories, delete posts and comments, and respond to reports.
- **Steps**:
  - Modify `create_db.sql` to include the following roles and permissions:
    ```sql
    -- Create the UserRole table
    CREATE TABLE UserRole (
        id INTEGER PRIMARY KEY,
        role_name VARCHAR(10),
        description VARCHAR(250),
        can_post BOOLEAN,
        can_comment BOOLEAN,
        can_like BOOLEAN,
        can_dislike BOOLEAN,
        can_delete BOOLEAN,
        can_manage_category BOOLEAN,
        can_approve_post BOOLEAN,
        can_promote BOOLEAN
    );

    -- Insert the default user roles
    INSERT INTO UserRole 
        (id, role_name, description, can_post, can_comment, can_like, can_dislike, can_delete, can_manage_category, can_approve_post, can_promote)
        VALUES 
        (0, 'guest', 'Can only view posts and comments', 0, 0, 0, 0, 0, 0, 0, 0),
        (1, 'user', 'Can create, comment, like, and dislike posts', 1, 1, 1, 1, 0, 0, 0, 0),
        (2, 'moderator', 'Can delete posts, report posts to admin, and approve posts', 1, 1, 1, 1, 1, 0, 1, 0),
        (3, 'admin', 'Can promote/demote users, manage categories, delete posts and comments, and respond to reports', 1, 1, 1, 1, 1, 1, 1, 1);
    ```

### 1.2. Update Structs
- **Objective**: Reflect the new roles and permissions in the Go structs.
- **Steps**:
  - Update `internal/structs/structs.go` to include new user roles and permissions:
    ```go
    type UserRole struct {
        Id               int
        RoleName         string
        Description      string
        CanPost          bool
        CanComment       bool
        CanLike          bool
        CanDislike       bool
        CanDelete        bool
        CanManageCategory bool
        CanApprovePost   bool
        CanPromote       bool
    }
    ```

### 1.3. Update Database Queries
- **Objective**: Ensure the database queries handle the new roles and permissions.
- **Steps**:
  - Modify `internal/database/get.go` to retrieve user roles and permissions:
    ```go
    func GetUserRoles() ([]structs.UserRole, error) {
        rows, err := db.Query(`SELECT id, role_name, description, can_post, can_comment, can_like, can_dislike, can_delete, can_manage_category, can_approve_post, can_promote FROM UserRole`)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        var roles []structs.UserRole
        for rows.Next() {
            var role structs.UserRole
            err := rows.Scan(&role.Id, &role.RoleName, &role.Description, &role.CanPost, &role.CanComment, &role.CanLike, &role.CanDislike, &role.CanDelete, &role.CanManageCategory, &role.CanApprovePost, &role.CanPromote)
            if err != nil {
                return nil, err
            }
            roles = append(roles, role)
        }
        return roles, nil
    }
    ```

### 1.4. Update User Handlers
- **Objective**: Implement role-based access control in user-related handlers.
- **Steps**:
  - Modify `internal/server/user.go` to check user roles and permissions before performing actions like posting, commenting, liking, etc.
    ```go
    func canPerformAction(user *structs.User, action string) bool {
        switch action {
        case "post":
            return user.Role.CanPost
        case "comment":
            return user.Role.CanComment
        case "like":
            return user.Role.CanLike
        case "dislike":
            return user.Role.CanDislike
        case "delete":
            return user.Role.CanDelete
        case "manage_category":
            return user.Role.CanManageCategory
        case "approve_post":
            return user.Role.CanApprovePost
        case "promote":
            return user.Role.CanPromote
        default:
            return false
        }
    }
    ```

## 2. Implement Post Approval System (files involved: `create_db.sql`, `internal/structs/structs.go`, `internal/database/get.go`, `internal/database/insert.go`, `internal/server/post.go`)
### 2.1. Update Database Schema
- **Objective**: Add a mechanism to mark posts as approved or pending approval.
- **Steps**:
  - Modify `create_db.sql` to include a column `is_approved` in the `Post` table:
    ```sql
    -- Create the Post table
    CREATE TABLE Post (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        parent_id INTEGER,
        title VARCHAR(64),
        message TEXT,
        image_id INTEGER,
        time DATE,
        is_approved BOOLEAN DEFAULT 0,
        like_count INTEGER,
        dislike_count INTEGER,
        love_count INTEGER,
        haha_count INTEGER,
        skull_count INTEGER,
        FOREIGN KEY (parent_id) REFERENCES Post(id),
        FOREIGN KEY (user_id) REFERENCES User(id),
        FOREIGN KEY (image_id) REFERENCES UploadedImage(id)
    );
    ```

### 2.2. Update Structs
- **Objective**: Reflect the post approval status in the Go structs.
- **Steps**:
  - Update `internal/structs/structs.go` to include a field `IsApproved` in the `Post` struct:
    ```go
    type Post struct {
        Id            int
        UserId        int
        ParentId      *int
        Title         string
        Message       string
        ImageId       int
        Time          time.Time
        IsApproved    bool
        CategoriesIDs []int
    }
    ```

### 2.3. Update Database Queries
- **Objective**: Ensure the database queries handle post approval status.
- **Steps**:
  - Modify `internal/database/get.go` to retrieve posts based on their approval status:
    ```go
    func GetApprovedPosts(count, offset int) ([]structs.Post, error) {
        rows, err := db.Query(`SELECT id, user_id, parent_id, title, message, image_id, time FROM Post WHERE is_approved = 1 ORDER BY time DESC LIMIT ? OFFSET ?`, count, offset)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        return getPostsHelper(rows)
    }
    ```

### 2.4. Update Post Handlers
- **Objective**: Implement logic to handle post approval and filtering based on categories.
- **Steps**:
  - Modify `internal/server/post.go` to check if the user is a moderator or admin before approving posts:
    ```go
    func approvePostHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanApprovePost {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postId := r.URL.Query().Get("post_id")
        err := database.ApprovePost(postId)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
    ```

## 3. Implement Moderator Functions (files involved: `internal/server/moderator.go`, `internal/database/get.go`, `internal/database/insert.go`)
### 3.1. Create Moderator Handlers
- **Objective**: Allow moderators to delete or report posts.
- **Steps**:
  - Create `internal/server/moderator.go` to handle moderator-specific functions:
    ```go
    func deletePostHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanDelete {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postId := r.URL.Query().Get("post_id")
        err := database.DeletePost(postId)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }

    func reportPostHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanReport {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var report structs.Report
        if !helpers.ParseBody(&report, r) {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        err := database.AddReport(report)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
    ```

### 3.2. Update Database Queries
- **Objective**: Ensure the database queries handle moderator actions.
- **Steps**:
  - Modify `internal/database/get.go` to retrieve posts reported by moderators:
    ```go
    func GetReportedPosts() ([]structs.Post, error) {
        rows, err := db.Query(`SELECT id, user_id, parent_id, title, message, image_id, time FROM Post WHERE is_reported = 1 ORDER BY time DESC`)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        return getPostsHelper(rows)
    }
    ```

  - Modify `internal/database/insert.go` to insert and update reports made by moderators:
    ```go
    func AddReport(report structs.Report) error {
        _, err := db.Exec(`INSERT INTO Report (reporter_user_id, reported_user_id, report_message, reported_post_id, time, is_post_report, is_pending) VALUES (?, ?, ?, ?, ?, ?, ?)`,
            report.ReporterId, report.ReportedId, report.Reason, report.PostId, report.Time, report.IsPostReport, true)
        return err
    }
    ```

## 4. Implement Administrator Functions (files involved: `internal/server/admin.go`, `internal/database/get.go`, `internal/database/insert.go`)
### 4.1. Create Admin Handlers
- **Objective**: Allow administrators to manage users, categories, and respond to reports.
- **Steps**:
  - Create `internal/server/admin.go` to handle admin-specific functions:
    ```go
    func promoteUserHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanPromote {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var promoteRequest structs.PromoteUserRequest
        if !helpers.ParseBody(&promoteRequest, r) {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        err := database.PromoteUser(promoteRequest.Username)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }

    func manageCategoryHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanManageCategory {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var categoryRequest structs.CategoryRequest
        if !helpers.ParseBody(&categoryRequest, r) {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        err := database.CreateCategory(categoryRequest)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }

    func respondToReportHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanRespondToReports {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var reportResponse structs.ReportResponse
        if !helpers.ParseBody(&reportResponse, r) {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        err := database.RespondToReport(reportResponse)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
    ```

### 4.2. Update Database Queries
- **Objective**: Ensure the database queries handle admin actions.
- **Steps**:
  - Modify `internal/database/get.go` to retrieve user roles and reports:
    ```go
    func GetUserRoles() ([]structs.UserRole, error) {
        rows, err := db.Query(`SELECT id, role_name, description, can_post, can_comment, can_like, can_dislike, can_delete, can_manage_category, can_approve_post, can_promote FROM UserRole`)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        var roles []structs.UserRole
        for rows.Next() {
            var role structs.UserRole
            err := rows.Scan(&role.Id, &role.RoleName, &role.Description, &role.CanPost, &role.CanComment, &role.CanLike, &role.CanDislike, &role.CanDelete, &role.CanManageCategory, &role.CanApprovePost, &role.CanPromote)
            if err != nil {
                return nil, err
            }
            roles = append(roles, role)
        }
        return roles, nil
    }

    func GetReports() ([]structs.Report, error) {
        rows, err := db.Query(`SELECT id, reporter_user_id, reported_user_id, report_message, reported_post_id, time, is_post_report, is_pending FROM Report WHERE is_pending = 1`)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        var reports []structs.Report
        for rows.Next() {
            var report structs.Report
            err := rows.Scan(&report.Id, &report.ReporterId, &report.ReportedId, &report.Reason, &report.PostId, &report.Time, &report.IsPostReport, &report.IsPending)
            if err != nil {
                return nil, err
            }
            reports = append(reports, report)
        }
        return reports, nil
    }
    ```

  - Modify `internal/database/insert.go` to insert and update user roles and categories:
    ```go
    func PromoteUser(username string) error {
        _, err := db.Exec(`UPDATE User SET type_id = 2 WHERE username = ?`, username)
        return err
    }

    func CreateCategory(category structs.CategoryRequest) error {
        _, err := db.Exec(`INSERT INTO Category (name, description, color) VALUES (?, ?, ?)`, category.Name, category.Description, category.Color)
        return err
    }

    func RespondToReport(response structs.ReportResponse) error {
        _, err := db.Exec(`UPDATE Report SET is_pending = 0, report_response = ? WHERE id = ?`, response.Response, response.ReportId)
        return err
    }
    ```

## 5. Implement Category Management (files involved: `internal/server/category.go`, `internal/database/get.go`, `internal/database/insert.go`)
### 5.1. Create Category Handlers
- **Objective**: Allow administrators to create and delete categories.
- **Steps**:
  - Create `internal/server/category.go` to handle category creation and deletion:
    ```go
    func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanManageCategory {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var categoryRequest structs.CategoryRequest
        if !helpers.ParseBody(&categoryRequest, r) {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        err := database.CreateCategory(categoryRequest)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }

    func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
        user := sessionmanager.GetUser(r)
        if user == nil || !user.Role.CanManageCategory {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        categoryId := r.URL.Query().Get("category_id")
        err := database.DeleteCategory(categoryId)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
    ```

### 5.2. Update Database Queries
- **Objective**: Ensure the database queries handle category management.
- **Steps**:
  - Modify `internal/database/get.go` to retrieve categories:
    ```go
    func GetCategories() ([]structs.Category, error) {
        rows, err := db.Query(`SELECT id, name, description, color FROM Category`)
        if err != nil {
            return nil, err
        }
        defer rows.Close()

        var categories []structs.Category
        for rows.Next() {
            var category structs.Category
            err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.Color)
            if err != nil {
                return nil, err
            }
            categories = append(categories, category)
        }
        return categories, nil
    }
    ```

  - Modify `internal/database/insert.go` to insert and delete categories:
    ```go
    func CreateCategory(category structs.CategoryRequest) error {
        _, err := db.Exec(`INSERT INTO Category (name, description, color) VALUES (?, ?, ?)`, category.Name, category.Description, category.Color)
        return err
    }

    func DeleteCategory(categoryId string) error {
        _, err := db.Exec(`DELETE FROM Category WHERE id = ?`, categoryId)
        return err
    }
    ```

## 6. Update Frontend (files involved: `www/template/*.html`, `www/static/css/*.css`, `www/static/js/*.js`)
### 6.1. Update Templates
- **Objective**: Reflect the new roles and post approval status in the UI.
- **Steps**:
  - Modify `www/template/*.html` to include UI changes for different user roles and post approval status:
    ```html
    <!-- Example: Add a section for moderators to approve posts -->
    {{ if .User.Role.CanApprovePost }}
    <div class="approve-posts-section">
        <h2>Pending Posts</h2>
        {{ range .PendingPosts }}
        <div class="post">
            <h3>{{ .Title }}</h3>
            <p>{{ .Message }}</p>
            <button onclick="approvePost({{ .Id }})">Approve</button>
            <button onclick="deletePost({{ .Id }})">Delete</button>
        </div>
        {{ end }}
    </div>
    {{ end }}
    ```

### 6.2. Update Styles
- **Objective**: Style the new UI elements.
- **Steps**:
  - Modify `www/static/css/*.css` to style the new UI elements for different user roles and post approval status:
    ```css
    .approve-posts-section {
        margin: 20px 0;
    }

    .approve-posts-section .post {
        border: 1px solid #ccc;
        padding: 10px;
        margin-bottom: 10px;
    }

    .approve-posts-section .post h3 {
        margin: 0 0 10px;
    }

    .approve-posts-section .post button {
        margin-right: 10px;
    }
    ```

### 6.3. Update Scripts
- **Objective**: Handle frontend logic for different user roles and post approval status.
- **Steps**:
  - Modify `www/static/js/*.js` to handle actions like approving, deleting, and reporting posts:
    ```js
    function approvePost(postId) {
        fetch(`/post/approve?post_id=${postId}`, {
            method: 'POST',
        })
        .then(response => {
            if (response.ok) {
                alert('Post approved successfully');
                location.reload();
            } else {
                alert('Failed to approve post');
            }
        });
    }

    function deletePost(postId) {
        fetch(`/post/delete?post_id=${postId}`, {
            method: 'POST',
        })
        .then(response => {
            if (response.ok) {
                alert('Post deleted successfully');
                location.reload();
            } else {
                alert('Failed to delete post');
            }
        });
    }
    ```

## 7. Testing and Debugging
### 7.1. Unit Tests
- **Objective**: Ensure the new functionalities work as expected.
- **Steps**:
  - Write unit tests for new functionalities in `internal/server/*_test.go`, `internal/database/*_test.go`:
    ```go
    func TestApprovePost(t *testing.T) {
        // Setup test data
        post := structs.Post{Id: 1, IsApproved: false}
        database.CreatePost(post)

        // Call the function
        err := database.ApprovePost(post.Id)
        if err != nil {
            t.Fatalf("Expected no error, got %v", err)
        }

        // Verify the result
        updatedPost, _ := database.GetPostById(post.Id)
        if !updatedPost.IsApproved {
            t.Fatalf("Expected post to be approved")
        }
    }
    ```

### 7.2. Integration Tests
- **Objective**: Ensure the system works as a whole.
- **Steps**:
  - Write integration tests to ensure the system works as expected:
    ```go
    func TestUserRolePermissions(t *testing.T) {
        // Setup test data
        user := structs.User{Id: 1, Role: structs.UserRole{CanPost: true, CanApprovePost: false}}
        database.CreateUser(user)

        // Call the function
        canPost := canPerformAction(&user, "post")
        canApprove := canPerformAction(&user, "approve_post")

        // Verify the result
        if !canPost {
            t.Fatalf("Expected user to be able to post")
        }
        if canApprove {
            t.Fatalf("Expected user not to be able to approve posts")
        }
    }
    ```

### 7.3. Manual Testing
- **Objective**: Identify and fix any issues.
- **Steps**:
  - Perform manual testing to identify and fix any issues:
    - Test creating, approving, and deleting posts.
    - Test promoting and demoting users.
    - Test managing categories.
    - Test reporting posts and responding to reports.

## 8. Documentation
### 8.1. Update README
- **Objective**: Provide information about the new moderation system and user roles.
- **Steps**:
  - Update the README file to include information about the new moderation system and user roles:
    ```markdown
    ## Forum Moderation System

    This forum includes a moderation system with the following user roles:

    - **Guest**: Can only view posts and comments.
    - **User**: Can create, comment, like, and dislike posts.
    - **Moderator**: Can delete posts, report posts to the admin, and approve posts.
    - **Administrator**: Can promote/demote users, manage categories, delete posts and comments, and respond to reports.

    ### Post Approval System

    Posts created by users need to be approved by a moderator or admin before they become publicly visible. Moderators and admins can approve or delete pending posts.

    ### Category Management

    Administrators can create and delete categories to organize posts.
    ```

### 8.2. Add Documentation
- **Objective**: Provide detailed documentation for developers and users.
- **Steps**:
  - Add detailed documentation for developers and users on how to use the new features:
    ```markdown
    ## Developer Documentation

    ### User Roles and Permissions

    The system includes the following user roles with specific permissions:

    - **Guest**: Can only view posts and comments.
    - **User**: Can create, comment, like, and dislike posts.
    - **Moderator**: Can delete posts, report posts to the admin, and approve posts.
    - **Administrator**: Can promote/demote users, manage categories, delete posts and comments, and respond to reports.

    ### Post Approval System

    - **Creating a Post**: Users can create posts, but they need to be approved by a moderator or admin before they become publicly visible.
    - **Approving a Post**: Moderators and admins can approve or delete pending posts.

    ### Category Management

    - **Creating a Category**: Administrators can create new categories to organize posts.
    - **Deleting a Category**: Administrators can delete categories.

    ## User Documentation

    ### Creating a Post

    1. Log in as a user.
    2. Navigate to the "Create Post" page.
    3. Fill in the post details and submit the form.
    4. Wait for a moderator or admin to approve your post.

    ### Approving a Post

    1. Log in as a moderator or admin.
    2. Navigate to the "Pending Posts" section.
    3. Review the pending posts and click "Approve" or "Delete" as needed.

    ### Managing Categories

    1. Log in as an admin.
    2. Navigate to the "Manage Categories" page.
    3. Use the form to create or delete categories.
    ```

