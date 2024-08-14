function calculateAge(birthDateString) {
    const birthDate = new Date(birthDateString);
    const today = new Date();
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDifference = today.getMonth() - birthDate.getMonth();

    // If the birth date hasn't occurred yet this year, subtract one year from the age.
    if (monthDifference < 0 || (monthDifference === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }

    return age;
}

function profile(userId, caseString) {
    const ReporBtutton = document.getElementById("ReportProfile-button");

    // Send the form data to the Go server
    fetch("http://localhost:8080/userType", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                handleErrorResponse(response);
            }
            return response.json();
        })
        .then((type) => {
            if (type.name === "Guest" || type.name === "User") {
                ReporBtutton.style.display = "none";
            } else if (userId != -1) {
                ReporBtutton.style.display = "block";
            }
        })
        .catch((error) => {
            ReporBtutton.style.display = "none";
        });

    const url = new URL("http://localhost:8080/userProfile");
    url.searchParams.append("user_id", userId);
    url.searchParams.append("case", caseString);

    fetch(url, {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                handleErrorResponse(response);
            }
            return response.json();
        })
        .then((data) => {
            if (data.User && data.User.id != data.UserProfile.id) {
                ReporBtutton.setAttribute(
                    "onclick",
                    `ReportPost(${data.UserProfile.id},true,false)`
                );
            } else {
                ReporBtutton.style.display = "none";
            }
            const profileUsername = document.getElementById("profileUsername");
            profileUsername.innerHTML = data.UserProfile.username;

            // Calculate and display the age
            const birthDate = data.UserProfile.DateOfBirth;
            const age = calculateAge(birthDate);
            const profileAge = document.getElementById("userInfo");
            profileAge.innerHTML = `Age: ${age},  location: ${data.UserProfile.location}`;

            // console.log('userPic', data.UserProfile.image_url)
            const userPic = document.getElementById("userPic");
            userPic.style.backgroundImage = `url(data:image/jpeg;base64,${data.UserProfile.image_url})`;
            // Update the onclick attributes with the userId

            // Update the onclick attributes with the userId
            const postsThElement = document.getElementById("PostsTh");
            postsThElement.setAttribute(
                "onclick",
                `posts(${userId},'Posts','Posts', event.target)`
            );

            const commentsThElement = document.getElementById("commentsTh");
            commentsThElement.setAttribute(
                "onclick",
                `posts(${userId},'comments','comments', event.target)`
            );

            const likesThElement = document.getElementById("likesTh");
            likesThElement.setAttribute(
                "onclick",
                `posts(${userId},'likes','likes', event.target)`
            );

            const dislikesThElement = document.getElementById("dislikesTh");
            dislikesThElement.setAttribute(
                "onclick",
                `posts(${userId},'dislikes','dislikes', event.target)`
            );

            posts(userId, "Posts", "Posts", postsThElement);
        })
        .catch((error) => {
            console.error("Error fetching profile:", error);
        });
    toggleVisibility("profile");
}

function posts(userId, caseString, column, element) {
    const url = new URL("http://localhost:8080/userProfile");
    url.searchParams.append("user_id", userId);
    url.searchParams.append("case", caseString);
    changeContent(column, element);

    fetch(url, {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                handleErrorResponse(response);
            }
            return response.json();
        })
        .then((data) => {
            if (data && data.Posts) {
                displayPostOnProfile(data.Posts, caseString);
            } else {
                console.error("Invalid data format. Expected profileView with Posts.");
            }
        })
        .catch((error) => {
            console.error("Error fetching profile:", error);
        });
}

function displayPostOnProfile(Posts, caseString) {
    if (Array.isArray(Posts)) {
        createPost(Posts, "profileContent", caseString);
    } else {
        console.error("Invalid data format. Expected an array of posts.");
    }
}

function ChatView() {
    const messagesBoxDiv = document.getElementById("messagesBoxDiv");
    if (messagesBoxDiv.style.display != "none") {
        fetch("http://localhost:8080/ChatView", {
                method: "GET",
            })
            .then((response) => response.json())
            .then((data) => {
                const usersList = document.getElementById("messagesBoxDiv");
                usersList.innerHTML = "";

                data.forEach((chat) => {
                    const messageBox = document.createElement("div");
                    messageBox.className = "messageBox";
                    messageBox.id = `messageBox-${chat.UserId}`;
                    const chatUserPic = document.createElement("div");
                    chatUserPic.className = "chatUserPic";
                    chatUserPic.style.backgroundImage = `url(data:image/jpeg;base64,${chat.Image})`;
                    chatUserPic.style.backgroundSize = "cover"; // Ensure the image covers the element

                    // Set border color based on online status
                    if (chat.Online) {
                        chatUserPic.style.border = " 3px solid rgb(74, 250, 58);";
                    } else {
                        chatUserPic.style.border = " 3px solid rgb(251, 217, 152, 0.7)";
                    }

                    const chatUserName = document.createElement("div");
                    chatUserName.className = "chatUserName";
                    const userNameP = document.createElement("p");
                    userNameP.textContent = chat.Username;
                    chatUserName.appendChild(userNameP);

                    messageBox.appendChild(chatUserPic);
                    messageBox.appendChild(chatUserName);
                    const typingIcon = document.createElement("div");
                    typingIcon.className = "typingIcon";
                    typingIcon.style.display = "none";
                    messageBox.appendChild(typingIcon);
                    messageBox.setAttribute(
                        "onclick",
                        `OpenMesages('${chat.Username}','${chat.UserId}')`
                    );
                    usersList.appendChild(messageBox);
                });
            })
            .catch((error) => {
                console.error("Error fetching:", error);
            });
    }
}
// Call ChatView initially
ChatView();

// Set an interval to refresh the chat view every 5 seconds
setInterval(ChatView, 5000);