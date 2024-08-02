async function submitSelection(divName) {
    // Get the selected categories
    const categoryCheckboxes = document.querySelectorAll(
        '.categoriesList input[type="checkbox"]:checked'
    );

    const selectedCategories = Array.from(categoryCheckboxes).map((checkbox) => checkbox.value);

    // Get the title and topic
    const title = document.getElementById("newPostTitle").value;
    const topic = document.getElementById("topic").value;

    console.log("Title:", title);
    console.log("Topic:", topic);
    console.log("Selected Categories:", selectedCategories);

    // Create the request body
    const formData = new FormData();
    formData.append("title", title);
    formData.append("topic", topic);
    formData.append("selectedCategories", JSON.stringify(selectedCategories));

    try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/post/add/Post", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Post added successfully");
            HomePageRequest();
            toggleDiv(divName); // Hide the content div after submission
            document.getElementById("newPostTitle").value = "";
            document.getElementById("topic").value = "";
        } else {
            const error = await response.text();
            console.error("Post addition failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// display posts in home
function HomePageRequest() {
    GetUserLoggedIn();
    // Send the form data to the Go server
    fetch("http://localhost:8080/homePageDataHuncler", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            displayPost(data); // Pass the response data to the displayScores function
        })
        .catch((error) => {
            console.error("Error:", error);
            console.log(response);
        });
}

HomePageRequest();

function displayPost(data) {
    console.log("Received data:", data); // Log the received data to check its structure

    if (Array.isArray(data.Posts)) {
        const homeNavigationContent = document.getElementById("home navigationContent");

        data.Posts.forEach((post) => {
            const postBox = document.createElement("div");
            postBox.classList.add("postBox");
            postBox.setAttribute("onclick", `PostPageHandler(${JSON.stringify(post)})`);
            postBox.setAttribute("id", `${post.Id}`);

            // // Create a EditButton element
            // const editButton = document.createElement("button");
            // editButton.classList.add("editPost-button");
            // editButton.setAttribute("onclick", "EditPost()");
            // editButton.textContent = "Edit";
            // postBox.appendChild(editButton);

            const postUserPic = document.createElement("div");
            postUserPic.classList.add("postUserPic");
            postBox.appendChild(postUserPic);

            const postTitle = document.createElement("div");
            postTitle.classList.add("postTitle");

            const titleElement = document.createElement("span");
            titleElement.classList.add("title");
            titleElement.textContent = post.title;

            const postUserName = document.createElement("span");
            postUserName.classList.add("postUserName");

            postUserName.textContent = post.author.username;

            const postContent = document.createElement("span");
            postContent.classList.add("postContent");

            postContent.textContent = post.message;

            postTitle.appendChild(titleElement);
            postTitle.appendChild(postUserName);
            postTitle.appendChild(postContent);

            postBox.appendChild(postTitle);

            homeNavigationContent.appendChild(postBox);
        });
    } else {
        console.error("Invalid data format. Expected an array of posts.");
    }
}

function fetchAndAppendCategories() {
    // Fetch the categories from the server
    const categoriesContainer = document.querySelector(".categoriesList");
    categoriesContainer.innerHTML = "";

    fetch("http://localhost:8080/category")
        .then((response) => response.json())
        .then((data) => {
            // Iterate over the categories and append them to the container
            data.forEach((category) => {
                appendCategory(category.Name, category.Id);
            });
        })
        .catch((error) => {
            console.error("Error fetching categories:", error);
        });
}

function appendCategory(name, Id) {
    // Get the container element where the categories will be added
    const categoriesContainer = document.querySelector(".categoriesList");

    // Create a new label element
    const label = document.createElement("label");

    // Create a new checkbox input element
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.name = name;
    checkbox.value = Id;
    checkbox.classList.add("checkBoxCateg");

    // Create the category name text node
    const categoryText = document.createTextNode(name);

    // Append the checkbox and category name to the label
    label.appendChild(checkbox);
    label.appendChild(categoryText);

    // Append the label to the container
    categoriesContainer.appendChild(label);
    categoriesContainer.appendChild(document.createElement("br"));
}

fetchAndAppendCategories();
HomePageRequest();

function Search() {
    const search = document.getElementById("search").value;

    // Create the request body
    const formData = new FormData();
    formData.append("search", search);

    fetch("http://localhost:8080/search", {
            method: "POST",
            body: formData,
        })
        .then((response) => response.json())
        .then((data) => {
            displaySearchPosts(data);
        })
        .catch((error) => {
            console.error("Error fetching posts:", error);
        });
}

const searchoutput = document.getElementById("search-output");

function displaySearchPosts(data) {
    // Clear any existing posts
    searchoutput.innerHTML = "";

    console.log("Received data:", data); // Log the received data to check its structure

    if (Array.isArray(data.Posts) && data.Posts.length != 0) {
        data.Posts.forEach((post) => {
            const postBox = document.createElement("div");
            postBox.classList.add("postResult");
            postBox.setAttribute("onclick", "toggleVisibility('postPage')");

            const postTitleElement = document.createElement("h2");
            postTitleElement.classList.add("title");
            postTitleElement.setAttribute(
                "style",
                '--font-selector:R0Y7S3VsaW0gUGFyay1yZWd1bGFy;--framer-font-family:"Kulim Park";--framer-font-size:40px;color:black;'
            );
            postTitleElement.textContent = post.title;
            postBox.appendChild(postTitleElement);

            searchoutput.appendChild(postBox);
        });
    } else {
        console.error("Invalid data format. Expected an array of posts.");
        const noResultElement = document.createElement("h2");
        noResultElement.setAttribute(
            "style",
            '--font-selector:R0Y7S3VsaW0gUGFyay1yZWd1bGFy;--framer-font-family:"Kulim Park";--framer-font-size:40px;color:black;'
        );
        noResultElement.textContent = "No result availabe";
        searchoutput.appendChild(noResultElement);
    }
}

async function submitLoginForm() {
    const usernameInput = document.getElementById("usernameLogin");
    const passwordInput = document.getElementById("passwordLogin");

    const username = usernameInput.value;
    const password = passwordInput.value;
    console.log(username);

    // Create the request body
    const formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);

    try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Login successful");
            const loginSpan = document.getElementById("loginSpan");
            loginSpan.innerHTML = "logout";
            GetUserLoggedIn();
            toggleVisibility("home");
        } else {
            const error = await response.text();
            console.log(error);
            alert("Login failed: " + error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

async function submitForm() {
    // Get the form values
    const firstName = document.getElementById("firstName").value;
    const lastName = document.getElementById("lastName").value;
    const email = document.getElementById("email").value;
    const username = document.getElementById("username").value;
    const dob = document.getElementById("dob").value;
    const password = document.getElementById("password").value;

    // Create the request body
    const formData = new FormData();
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
    formData.append("username", username);
    formData.append("dob", dob);
    formData.append("password", password);

    try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/signup", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Signup successful");
            GetUserLoggedIn();
            const loginSpan = document.getElementById("loginSpan");
            loginSpan.innerHTML = "logout";
            toggleVisibility("home");
            // Get the form values
            const firstName = document.getElementById("firstName");
            firstName.innerHTML = "";
            const lastName = document.getElementById("lastName");
            lastName.innerHTML = "";
            const email = document.getElementById("email");
            email.innerHTML = "";
            const username = document.getElementById("username");
            username.innerHTML = "";
            const dob = document.getElementById("dob");
            dob.innerHTML = "";
            const password = document.getElementById("password");
            password.innerHTML = "";
        } else {
            const error = await response.text();
            alert("Signup failed: " + error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    GetUserLoggedIn();
    const editPostButton = document.getElementById("editPost-button");
    if (editPostButton) {
        editPostButton.style.display = "none";
    }
});

function PostPageHandler(data) {
    const editPostButton = document.getElementById("editPost-button");
    toggleVisibility("postPage");
    if (editPostButton) {
        fetch("http://localhost:8080/user", {
                method: "GET",
            })
            .then((response) => {
                if (!response.ok) {
                    throw new Error("Error: " + response.status);
                }
                return response.text(); // Get response as text
            })
            .then((text) => {
                console.log("Raw response text:", text); // Log the raw response text
                try {
                    // Clean the response text by removing the trailing 'null'
                    const cleanedText = text.replace(/null$/, "").trim();
                    const data1 = JSON.parse(cleanedText); // Parse the cleaned text as JSON
                    console.log(data1);
                    if (data1 === data.author.username) {
                        console.log("hello");
                        editPostButton.style.display = "block";
                    } else {
                        console.log("hi");
                        editPostButton.style.display = "none";
                    }
                } catch (error) {
                    console.error("JSON parsing error:", error);
                    console.error("Raw response text causing error:", text); // Log the problematic text
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }
}

function EditPostHandler() {}

async function GetUserLoggedIn() {
    fetch("http://localhost:8080/user", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.text(); // Get response as text
        })
        .then((text) => {
            const setting = document.querySelector(".settingIcon.navigationBarIcons");
            const bellIcon = document.querySelector(".bellIcon.navigationBarIcons");
            const settingBtn = document.querySelector(".settingBtn.navigationBarBtns");
            const notificationsBtn = document.querySelector(".notificationsBtn.navigationBarBtns");
            const penIcon = document.querySelector(".penIcon.mainPageIcons");
            const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
            const profileBtn = document.querySelector(".profileBtn.navigationBarBtns");
            if (text === "null") {
                setting.style.display = "none";
                bellIcon.style.display = "none";
                settingBtn.style.display = "none";
                notificationsBtn.style.display = "none";
                penIcon.style.display = "none";
                profileIcon.style.display = "none";
                profileBtn.style.display = "none";
                console.log("User not logged in");
                const loginSpan = document.getElementById("loginSpan");
                loginSpan.innerHTML = "sign up";
                return;
            } else {
                setting.style.display = "block";
                bellIcon.style.display = "block";
                settingBtn.style.display = "block";
                notificationsBtn.style.display = "block";
                penIcon.style.display = "block";
                profileIcon.style.display = "block";
                profileBtn.style.display = "block";
                const loginSpan = document.getElementById("loginSpan");
                loginSpan.innerHTML = "logout";
                console.log("User is logged in");
                return data;
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

window.onload(GetUserLoggedIn());
GetUserLoggedIn();

async function logout() {
    try {
        const response = await fetch("http://localhost:8080/logout", {
            method: "POST",
        });

        if (response.ok) {
            console.log("Logout successful");
            GetUserLoggedIn();
            // Perform any additional actions needed after successful logout
            toggleVisibility("home");
        } else {
            const error = await response.text();
            console.error("Logout failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

function settingsHandler() {
    // Send the form data to the Go server
    fetch("http://localhost:8080/userType", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            console.log("User type:", data.name);
            if (data.name === "Guest" || data.name === "User") {
                console.log("User type:boooo");
                document.querySelectorAll("#settingList li").forEach((li) => {
                    if (
                        li.innerHTML.includes("List of Moderators") ||
                        li.innerHTML.includes("Promotion Requests") ||
                        li.innerHTML.includes("Manage Categories")
                    ) {
                        li.style.display = "none";
                    }
                });
            } else if (data.name === "Moderator" || data.name === "Admin") {
                document.querySelectorAll("#settingList li").forEach((li) => {
                    if (
                        li.innerHTML.includes("List of Moderators") ||
                        li.innerHTML.includes("Promotion Requests") ||
                        li.innerHTML.includes("Manage Categories")
                    ) {
                        li.style.display = "block";
                    } else if (li.innerHTML.includes("Request to be Moderator")) {
                        li.style.display = "none";
                    }
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });

    toggleVisibility("setting");
}

function fetchAndAppendModerators() {
    // Fetch the list of moderators from the server
    fetch("http://localhost:8080/Moderator")
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            const moderatorsList = document.querySelector("#list-moderators ul");
            moderatorsList.innerHTML = ""; // Clear existing list items

            data.forEach((moderator) => {
                const listItem = document.createElement("li");
                listItem.textContent = moderator.Username;
                listItem.dataset.id = moderator.Id; // Set the moderator ID as a data attribute

                const removeButton = document.createElement("button");
                removeButton.classList.add("remove-button");
                removeButton.textContent = "Remove";
                removeButton.setAttribute("onclick", `RemoveModerator(${moderator.Id})`);

                listItem.appendChild(removeButton);
                moderatorsList.appendChild(listItem);
            });
        })
        .catch((error) => {
            console.error("Error fetching moderators:", error);
        });
    toggleDiv("list-moderators");
}

function RemoveModerator(id) {
    confirmAction();
    console.log("hello");
    // Create the request body
    const formData = new FormData();
    formData.append("id", id);

    fetch("http://localhost:8080/RemoveModerator", {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                alert("Moderator removed successfully");
                fetchAndAppendModerators(); // Refresh the list of moderators
            } else {
                return response.text().then((error) => {
                    console.error("Moderator removal failed:", error);
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function saveNewPassword() {
    var newPassword = document.getElementById("password2").value.trim();
    var confirmPassword = document.getElementById("confirmPassword2").value.trim();
    var passwordError = document.getElementById("passwordError2");
    var confirmPasswordError = document.getElementById("confirmPasswordError2");

    passwordError.textContent = "";
    confirmPasswordError.textContent = "";

    // Log the values to the console for debugging
    console.log("New Password:", newPassword);
    console.log("Confirm Password:", confirmPassword);

    if (newPassword.length < 8) {
        passwordError.textContent = "Password must be at least 8 characters long!";
    } else if (newPassword !== confirmPassword) {
        confirmPasswordError.textContent = "Password confirmation does not match!";
    } else {
        const formData = new FormData();
        formData.append("password", newPassword);

        fetch("http://localhost:8080/changePassword", {
                method: "POST",
                body: formData,
            })
            .then((response) => {
                if (response.ok) {
                    alert("Password changed successfully!");
                } else {
                    alert("error changig the password");
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }
}

function PromotionRequests() {
    fetch("http://localhost:8080/PromotionRequests", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            const promotionRequestsList = document.querySelector("#promotion-requests ul");
            promotionRequestsList.innerHTML = ""; // Clear existing list items

            data.forEach((request) => {
                const listItem = document.createElement("li");
                listItem.textContent = request.Username;
                listItem.dataset.id = request.Id; // Set the moderator ID as a data attribute

                const approveButton = document.createElement("button");
                approveButton.classList.add("Show-button");
                approveButton.textContent = "Show";
                approveButton.setAttribute("onclick", `ShowPromotion(${request.Id})`);

                const rejectButton = document.createElement("button");
                rejectButton.classList.add("reject-button");
                rejectButton.textContent = "Reject";
                rejectButton.setAttribute("onclick", `RejectPromotion(${request.UserId})`);

                listItem.appendChild(approveButton);
                listItem.appendChild(rejectButton);
                promotionRequestsList.appendChild(listItem);
            });
        })
        .catch((error) => {
            console.error("Error fetching promotion requests:", error);
        });
    toggleDiv("promotion-requests");
}

function ShowPromotion(Id) {
    console.log("hello");

    const formData = new FormData();
    formData.append("id", Id);

    fetch("http://localhost:8080/ShowUserPromotion", {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            const promotionRequestsList = document.querySelector("#ShowPromotionRequest");

            // Clear previous buttons
            document.getElementById("ShowPromotionRequestUser").textContent = data.Username;
            const contentParagraph = document.querySelector("#ShowPromotionRequest p");
            contentParagraph.textContent = data.Reason;

            const approveButton = document.querySelector(".ShowPromotionRequestApprove-button");
            approveButton.setAttribute("onclick", `ApprovePromotion(${data.UserId})`);

            const rejectButton = document.querySelector(".ShowPromotionRequestReject-button");
            rejectButton.setAttribute("onclick", `RejectPromotion(${data.UserId})`);

            promotionRequestsList.appendChild(approveButton);
            promotionRequestsList.appendChild(rejectButton);
            console.log("Promotion data:", data);
            // Handle the response data here
        })
        .catch((error) => {
            console.error("Error:", error);
        });
    toggleDiv("ShowPromotionRequest");
}
// toggleDiv("request-moderator");

function RejectPromotion(userId) {
    const formData = new FormData();
    formData.append("userId", userId);

    fetch("http://localhost:8080/RejectPromotion", {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                alert("Promotion request rejected successfully");
                PromotionRequests();
            } else {
                return response.text().then((error) => {
                    console.error("Promotion request rejection failed:", error);
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
    const contentDiv = document.getElementById("ShowPromotionRequest");
    contentDiv.style.display = "none";
}

function ApprovePromotion(userId) {
    const formData = new FormData();
    formData.append("userId", userId);

    fetch("http://localhost:8080/ApprovePromotion", {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                alert("Promotion request approved successfully");
                PromotionRequests();
            } else {
                return response.text().then((error) => {
                    console.error("Promotion request approval failed:", error);
                });
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
    toggleDiv("ShowPromotionRequest");
}

function ManageCategories() {
    // Fetch the categories from the server
    const categoriesContainer = document.querySelector("#manage-categories ul");
    categoriesContainer.innerHTML = "";

    fetch("http://localhost:8080/category")
        .then((response) => response.json())
        .then((data) => {
            // Iterate over the categories and append them to the container
            data.forEach((category) => {
                appendCategoryToList(category.Name, category.Id);
            });
        })
        .catch((error) => {
            console.error("Error fetching categories:", error);
        });
    toggleDiv("manage-categories");
}

function appendCategoryToList(name, id) {
    // Get the container element where the categories will be added
    const categoriesContainer = document.querySelector("#manage-categories ul");

    // Create a new list item element
    const listItem = document.createElement("li");

    // Create the category name text node
    const categoryText = document.createTextNode(name);

    // Create the remove button
    const removeButton = document.createElement("button");
    removeButton.classList.add("remove-button");
    removeButton.textContent = "Remove";
    removeButton.onclick = () => removeCategory(id);

    // Append the category name and remove button to the list item
    listItem.appendChild(categoryText);
    listItem.appendChild(removeButton);

    // Append the list item to the container
    categoriesContainer.appendChild(listItem);
}

function removeCategory(id) {
    // Create the request body
    const formData = new FormData();
    formData.append("id", id);

    if (confirm("Are you sure you want to proceed?")) {
        fetch("http://localhost:8080/removeCategory", {
                method: "POST",
                body: formData,
            })
            .then((response) => {
                if (response.ok) {
                    alert("Category removed successfully");
                    ManageCategories(); // Refresh the list of categories
                } else {
                    return response.text().then((error) => {
                        console.error("Category removal failed:", error);
                    });
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    } else {
        // User clicked "Cancel" or closed the dialog
        console.log("Action canceled.");
    }
}

async function addCategory() {
    // Get the category name and description from the form
    const categoryName = document.getElementById("categoryName").value;
    const categoryDescription = document.getElementById("CategoryDescription").value;

    // Create the request body
    const formData = new FormData();
    formData.append("name", categoryName);
    formData.append("description", categoryDescription);

    try {
        // Send the POST request to add the category
        const response = await fetch("http://localhost:8080/addCategory", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Category added successfully");

            // Optionally, you can append the new category to the list
            appendCategoryToList(categoryName, await response.text());

            // Clear the form fields
            document.getElementById("categoryName").value = "";
            document.getElementById("CategoryDescription").value = "";

            // Hide the add category popup
            toggleDiv("addCategory");
        } else {
            const error = await response.text();
            console.error("Category addition failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

function ChangeUserInformation() {
    fetch("http://localhost:8080/getUserInfo", {
            method: "GET",
        })
        .then((response) => response.json())
        .then((data) => {
            document.getElementById("ChangeUserName").value = data.Username;
            document.getElementById("ChangeFirstName").value = data.FirstName;
            document.getElementById("ChangeLastName").value = data.LastName;
            document.getElementById("ChangeEmail").value = data.Email;
            console.log(data.DateOfBirth);
            // Parse and format the date to yyyy-MM-dd
            const date = new Date(data.DateOfBirth);
            const formattedDate = date.toISOString().split("T")[0];
            document.getElementById("ChangeDOB").value = formattedDate;
        })
        .catch((error) => {
            console.error("Error fetching user information:", error);
        });
    toggleDiv("user-info");
}