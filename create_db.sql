-- Drop existing tables (if necessary)
DROP TABLE IF EXISTS PostCategory;
DROP TABLE IF EXISTS ReactionType;
DROP TABLE IF EXISTS PostReaction;
DROP TABLE IF EXISTS User;
DROP TABLE IF EXISTS UserRole;
DROP TABLE IF EXISTS Report;
DROP TABLE IF EXISTS PromoteRequest;
DROP TABLE IF EXISTS Post;
DROP TABLE IF EXISTS Category;
DROP TABLE IF EXISTS UserNotification;
DROP TABLE IF EXISTS UserSession;
DROP TABLE IF EXISTS UploadedImage;
DROP TABLE IF EXISTS UserMessage;


-- Create the Image table
CREATE TABLE UploadedImage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    data BLOB
);

-- Create the User table
CREATE TABLE User (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id INTEGER NOT NULL,
    username VARCHAR(10) NOT NULL,  
    first_name VARCHAR(16) NOT NULL,
    last_name VARCHAR(16) NOT NULL,
    country VARCHAR(16) NOT NULL,
    date_of_birth DATE,
    email VARCHAR(30) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    image_id INTEGER NOT NULL,
    banned_until DATE,
    github_name VARCHAR(20),
    linkedin_name VARCHAR(20),
    twitter_name VARCHAR(20),
    bio VARCHAR(255),
    gender VARCHAR(10),
    FOREIGN KEY(image_id) REFERENCES UploadedImage(id),
    FOREIGN KEY (type_id) REFERENCES UserRole(id)
);

-- Create the User Role table
CREATE TABLE UserRole (
    id INTEGER PRIMARY KEY,
    role_name VARCHAR(10),
    description VARCHAR(250),
    can_post BOOLEAN,
    can_react BOOLEAN,
    can_manage_category BOOLEAN,
    can_delete BOOLEAN,
    can_report BOOLEAN,
    can_promote BOOLEAN
);

-- Create the Report table
CREATE TABLE Report (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reporter_user_id INTEGER,
    reported_user_id INTEGER,
    report_message VARCHAR(250),
    reported_post_id INTEGER,
    time DATE,
    is_post_report boolean,
    is_pending boolean,
    report_response VARCHAR(250),
    FOREIGN KEY (reported_user_id) REFERENCES User(id),
    FOREIGN KEY (reported_post_id) REFERENCES Post(id),
    FOREIGN KEY (reporter_user_id) REFERENCES User(id)
    
);

-- Create the PromoteRequest table
CREATE TABLE PromoteRequest (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    description TEXT,
    time DATE,
    is_pending boolean,
    FOREIGN KEY (user_id) REFERENCES User(id)
);

-- Create the Post table
CREATE TABLE Post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    parent_id INTEGER,
    title VARCHAR(64),
    message TEXT,
    image_id INTEGER,
    time DATE,
    like_count    INTEGER,
	dislike_count INTEGER,
    FOREIGN KEY (parent_id) REFERENCES Post(id),
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (image_id) REFERENCES UploadedImage(id)
);

-- Create the Notification table
CREATE TABLE UserNotification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    comment_id VARCHAR(10),
    post_reaction_id VARCHAR(150),
    report_id VARCHAR(150),
    promote_request_id VARCHAR(150), -- "accept promotion, depromtion, report accepted, not accepted (replay)"
    read BOOLEAN DEFAULT FALSE, -- Add this line
    FOREIGN KEY (comment_id) REFERENCES Post(id),
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (post_reaction_id) REFERENCES PostReaction(id),
    FOREIGN KEY (report_id) REFERENCES Report(id),
    FOREIGN KEY (promote_request_id) REFERENCES PromoteRequest(id)
);

-- Create the Category table
CREATE TABLE Category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(10),
    description VARCHAR(150),
    color VARCHAR(7)
);

-- Create the PostCategory table
CREATE TABLE PostCategory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES Category(id),
    FOREIGN KEY (post_id) REFERENCES Post(id)
);

-- Create the ReactionType table
CREATE TABLE ReactionType (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type VARCHAR(20)
);

-- Create the PostReaction table
CREATE TABLE PostReaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    reaction_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES User(id),
    FOREIGN KEY (post_id) REFERENCES Post(id),
    FOREIGN KEY (reaction_id) REFERENCES ReactionType(id)
);

-- Create the Session table
CREATE TABLE UserSession (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token CHAR(64),
    user_id INTEGER,
    creation_time INTEGER,
    FOREIGN KEY (user_id) REFERENCES User(id)
);

-- Create the message table
CREATE TABLE UserMessage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    receiver_id INTEGER,
    messag VARCHAR(300),
    time TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES User(id),
    FOREIGN KEY (receiver_id) REFERENCES User(id)
);


-- Insert the default user roles
INSERT INTO UserRole 
    (id, role_name, description, can_post, can_react, can_manage_category, can_delete, can_report, can_promote)
    VALUES 
    (0, 'guest', '', 0, 0, 0, 0, 0, 0),
    (1, 'user', '', 1, 1, 0, 0, 1, 0),
    (2, 'moderator', '', 1, 1, 1, 1, 1, 0),
    (3, 'admin', '', 1, 1, 1, 1, 1, 1);

-- Insert default categories
INSERT INTO Category 
    (id, name, description, color)
    VALUES 
    (1, 'General', 'General discussion', '#000000'),
    (2, 'Announcement', 'Announcements', '#000000'),
    (3, 'Question', 'Questions', '#000000');

-- Insert default reactions
INSERT INTO ReactionType
    (id, type)
    VALUES 
    (1, 'like'),
    (2, 'dislike');


-- insert default users image
INSERT INTO UploadedImage (id, data)
VALUES (
    0,
    (SELECT readfile('/home/mariam/mm/real-time-forum/static/images/user.png'))
);


INSERT INTO User
    (type_id, username, first_name, last_name, country, date_of_birth, email, hashed_password, image_id, banned_until, github_name, linkedin_name, twitter_name, bio, gender) 
    VALUES
    (1, 'johnny_doe', 'John', 'Smith', 'USA', '1992-03-21', 'john.smith@example.com', '$2a$10$yMwcWmnuJKVvHwDR4zRmQO0DmyWzt0wwU2BqdwGAOEcM0MIKjSZ/O', 0, NULL, 'john_smith', 'johnsmith', 'john_smith', 'Software engineer passionate about coding and technology.', 'Male'),
    (1, 'jane_smith', 'Jane', 'Johnson', 'Canada', '1987-08-15', 'jane.johnson@example.com', '$2a$10$yMwcWmnuJKVvHwDR4zRmQO0DmyWzt0wwU2BqdwGAOEcM0MIKjSZ/O', 0, NULL, 'jane_johnson', 'janejohnson', 'jane_johnson', 'Marketing professional with a love for creativity and innovation.', 'Female'),
    (1, 'alex_brown', 'Alex', 'Brown', 'UK', '1990-05-10', 'alex.brown@example.com', '$2a$10$yMwcWmnuJKVvHwDR4zRmQO0DmyWzt0wwU2BqdwGAOEcM0MIKjSZ/O', 0, NULL, 'alexbrown', 'alexbrown', 'alex_brown', 'Art enthusiast and aspiring designer exploring the world.', 'Male'),
    (3, 'rhelal', 'Ruqaya', 'Helal', 'Bahrain', '1998-03-21', 'ruqayahhelal@example.com', '$2a$10$yMwcWmnuJKVvHwDR4zRmQO0DmyWzt0wwU2BqdwGAOEcM0MIKjSZ/O', 0, NULL, 'rhelal', 'rhelal', 'rhelal', 'Travel blogger with a passion for photography and storytelling.', 'Female'),
    (3, 'maabbas', 'mariam', 'Abbas', 'Bahrain', '2004-12-15', 'maiam.abbas@example.com', '$2a$10$yMwcWmnuJKVvHwDR4zRmQO0DmyWzt0wwU2BqdwGAOEcM0MIKjSZ/O', 0, NULL, 'maiam.abbas', 'maiam.abbas', 'maiam.abbas', 'Fitness enthusiast and wellness advocate.', 'Female');


INSERT INTO Post 
    (user_id, parent_id, title, message, image_id, time, like_count, dislike_count) 
    VALUES 
    (1, NULL, 'Exciting News!', 'Just landed a new job today. Feeling grateful and excited for the new opportunity.', 0, '2024-08-15', 0, 0),
    (2, NULL, 'Weekend Adventure', 'Spent the weekend hiking in the mountains. The views were breathtaking!', 0, '2024-08-16', 0, 0),
    (3, NULL, 'Movie Night', 'Watched the latest blockbuster movie last night. It was so good!', 0, '2024-08-17', 0, 0),
    (4, NULL, 'Cooking Experiment', 'Tried a new recipe today. Surprisingly, it turned out delicious!', 0, '2024-08-18', 0, 0),
    (5, NULL, 'Fitness Journey', 'Completed my first 5k run today. Feeling accomplished!', 0, '2024-08-19', 0, 0),
    (1, NULL, 'Travel Plans', 'Excited to announce my upcoming trip to Europe. Cant wait to explore new places!', 0, '2024-08-20', 0, 0),
    (2, NULL, 'Artistic Inspiration', 'Visited a local art gallery today. Feeling inspired to create something beautiful.', 0, '2024-08-21', 0, 0),
    (3, NULL, 'Tech Update', 'Just got my hands on the latest smartphone. The technology is simply amazing!', 0, '2024-08-22', 0, 0),
    (4, NULL, 'Book Recommendation', 'Finished reading an incredible book today. Highly recommend it to all book lovers!', 0, '2024-08-23', 0, 0),
    (5, NULL, 'Gardening Delight', 'Spent the afternoon in my garden. Theres something so therapeutic about gardening.', 0, '2024-08-24', 0, 0),
    (1, NULL, 'New Art Project', 'Started working on a new art project today. Cant wait to see how it turns out!', 0, '2024-08-25', 0, 0),
    (2, NULL, 'Productivity Tips', 'Sharing my top productivity tips to help you stay focused and motivated.', 0, '2024-08-26', 0, 0),
    (3, NULL, 'DIY Project', 'Completed a fun DIY project over the weekend. It turned out better than expected!', 0, '2024-08-27', 0, 0),
    (4, NULL, 'Music Discovery', 'Discovered a new favorite band today. Their music is so catchy!', 0, '2024-08-28', 0, 0),
    (5, NULL, 'Family Time', 'Had a wonderful family get-together today. Cherishing these moments.', 0, '2024-08-29', 0, 0),
    (1, NULL, 'Coding Milestone', 'Reached a coding milestone today. Hard work pays off!', 0, '2024-08-30', 0, 0),
    (2, NULL, 'Outdoor Adventure', 'Spent the day kayaking on the river. Nature never fails to amaze me.', 0, '2024-08-31', 0, 0),
    (3, NULL, 'Healthy Eating', 'Trying out a new healthy eating plan. Heres to a healthier lifestyle!', 0, '2024-09-01', 0, 0),
    (4, NULL, 'Home Decor', 'Redecorated my living room today. Loving the new cozy vibe!', 0, '2024-09-02', 0, 0),
    (5, NULL, 'Volunteering Experience', 'Had a rewarding volunteering experience today. Giving back feels amazing!', 0, '2024-09-03', 0, 0);

--insertion on postCategory table:
INSERT INTO PostCategory (post_id, category_id) VALUES
    (1, 1),
    (1, 2),

    (2, 2),

    (3, 3),

    (4, 1),

    (5, 1),
    (5, 3),

    (6, 1),

    (7, 1),
    (7, 3),

    (8, 1),
    (8, 2),

    (9, 1),
    (9, 2),
    (9, 3),

    (10, 2),
    (10, 3),

    (11, 1),
    (11, 2),

    (12, 1),

    (13, 1),
    (13, 3),

    (14, 1),
    (14, 3),

    (15, 2),
    (15, 3),

    (16, 1),
    (16, 2),

    (17, 1),
    (17, 2),

    (18, 1),

    (19, 2),

    (20, 3);