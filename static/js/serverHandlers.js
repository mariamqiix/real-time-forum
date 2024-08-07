async function submitSelection(divName) {
    // Get the selected categories
    const categoryCheckboxes = document.querySelectorAll(
        '.categoriesList input[type="checkbox"]:checked'
    );

    const selectedCategories = Array.from(categoryCheckboxes).map((checkbox) => checkbox.value);
    if (selectedCategories.length === 0) {
        alert("Please select at least one category.");
        return;
    }
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

function fetchAndAppendCategoriesToFilter() {
    // Fetch the categories from the server
    const categoriesContainer = document.querySelector("#categoryFilter table");
    categoriesContainer.innerHTML = "";

    fetch("http://localhost:8080/category")
        .then((response) => response.json())
        .then((data) => {
            // Iterate over the categories and append them to the container
            data.forEach((category) => {
                const row = document.createElement("tr");
                const cell = document.createElement("td");
                const label = document.createElement("label");
                const checkbox = document.createElement("input");

                checkbox.type = "checkbox";
                checkbox.name = category.Name;
                checkbox.value = category.Name;

                label.appendChild(checkbox);
                label.appendChild(document.createTextNode(` ${category.Name}`));
                cell.appendChild(label);
                row.appendChild(cell);
                categoriesContainer.appendChild(row);
            });
            toggleDiv("categoryFilter");
        })
        .catch((error) => {
            console.error("Error fetching categories:", error);
        });
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

            // Select the element with id="sort"
            const sortElement = document.getElementById("sort");

            // Set its attributes to "old"
            sortElement.setAttribute("data-sort", "new");

            // Put the data in its onclick function
            sortElement.onclick = function() {
                sortPosts(data);
            };
            toggleVisibility("home");
        })
        .catch((error) => {
            console.error("Error:", error);
            console.log(response);
        });
}

function PostsByCategories() {
    // Gather selected options
    const selectedOptions = [];
    const checkboxes = document.querySelectorAll('#categoryFilter input[type="checkbox"]:checked');
    checkboxes.forEach((checkbox) => {
        selectedOptions.push(checkbox.value);
    });

    // Ensure there are selected options before making the request
    if (selectedOptions.length === 0) {
        alert("Please select at least one category.");
        return;
    }

    fetch("http://localhost:8080/postsByCategories", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ categories: selectedOptions }), // Ensure the key matches the server's expected key
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            displayPost(data); // Pass the response data to the displayPost function

            // Select the element with id="sort"
            const sortElement = document.getElementById("sort");

            // Set its attributes to "old"
            sortElement.setAttribute("data-sort", "new");

            // Put the data in its onclick function
            sortElement.onclick = function() {
                sortPosts(data);
            };
            toggleDiv("categoryFilter");
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function sortPosts(data) {
    const sortElement = document.getElementById("sort");
    const sortValue = sortElement.getAttribute("data-sort");

    if (Array.isArray(data.Posts)) {
        // Reverse the array based on the current sort value
        if (sortValue === "new") {
            data.Posts.reverse(); // Reverse the array
            sortElement.setAttribute("data-sort", "old");
        } else {
            data.Posts.reverse(); // Reverse the array
            sortElement.setAttribute("data-sort", "new");
        }

        // Display the sorted posts
        displayPost(data);
    } else {
        console.error("Invalid data format. Expected an array of posts.");
    }
}

HomePageRequest();

function displayPost(data) {
    if (Array.isArray(data.Posts)) {
        createPost(data.Posts, "homeContent");
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
    if (search == "" || search.trim() == "") {
        alert("Please enter a search term.");
        return;
    }
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
            postBox.setAttribute("onclick", `PostPageHandler(${JSON.stringify(post)})`);
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
            loginSpan.innerHTML = "LOGOUT";
            GetUserLoggedIn();
            HomePageRequest();
            initializeWebSocket();
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
    const country = document.getElementById("country").value;
    const password = document.getElementById("password").value;
    const gender = document.getElementById("gender").value;
    console.log(gender);
    // Create the request body
    const formData = new FormData();
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
    formData.append("country", country);
    formData.append("gender", gender);
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
            HomePageRequest(); // Get the form values
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
            const country = document.getElementById("country");
            country.innerHTML = "";
            const gender = document.getElementById("gender");
            gender.innerHTML = "";
            initializeWebSocket();
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
                        editPostButton.style.display = "block";
                    } else {
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
                const loginSpan = document.getElementById("loginSpan");
                loginSpan.innerHTML = "sign up";
                const replayPostButton = document.getElementById("replayPost-button");
                replayPostButton.style.display = "none";
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
                const replayPostButton = document.getElementById("replayPost-button");
                replayPostButton.style.display = "block";
                return text;
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}
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
            HomePageRequest();
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
            } else if (data.name === "Admin") {
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
            } else {
                document.querySelectorAll("#settingList li").forEach((li) => {
                    if (
                        li.innerHTML.includes("List of Moderators") ||
                        li.innerHTML.includes("Promotion Requests") ||
                        li.innerHTML.includes("Manage Categories") ||
                        li.innerHTML.includes("Request to be Moderator")
                    ) {
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

    fetch("h  ttp://localhost:8080/RemoveModerator", {
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

        fetch("h    ttp://localhost:8080/changePassword", {
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
    fetch("h ttp://localhost:8080/PromotionRequests", {
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

    fetch("h    ttp://localhost:8080/ShowUserPromotion", {
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

            // Optionally, you can append thenew category to the list
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
    fetch("h ttp://localhost:8080/getUserInfo", {
            method: "GET",
        })
        .then((response) => response.json())
        .then((data) => {
            if (data) {
                document.getElementById("ChangeUserName").value = data.Username || "";
                document.getElementById("ChangeFirstName").value = data.FirstName || "";
                document.getElementById("ChangeLastName").value = data.LastName || "";
                document.getElementById("ChangeEmail").value = data.Email || "";
                console.log(data.DateOfBirth);
                // Parse and format the date to yyyy-MM-dd
                const date = new Date(data.DateOfBirth);
                const formattedDate = date.toISOString().split("T")[0];
                document.getElementById("ChangeDOB").value = formattedDate || "";
                document.getElementById("ChangeCountry").value = data.Country || "";
                // document.getElementById("ChangeGender").value = data.Gender || "";
            } else {
                console.error("Data is undefined");
            }
        })
        .catch((error) => {
            console.error("Error fetching user information:", error);
        });
    toggleDiv("user-info");
}

function UpdateUserInformation() {
    if (!validateUserInfoForm()) {
        return;
    }
    const username = document.getElementById("ChangeUserName").value;
    const firstName = document.getElementById("ChangeFirstName").value;
    const lastName = document.getElementById("ChangeLastName").value;
    const email = document.getElementById("ChangeEmail").value;
    const dateOfBirth = document.getElementById("ChangeDOB").value;
    console.log(dateOfBirth);
    const country = document.getElementById("ChangeCountry").value;
    const gender = document.getElementById("ChangeGender").value;

    if (
        username.trim() === "" ||
        firstName.trim() === "" ||
        lastName.trim() === "" ||
        email.trim() === "" ||
        dateOfBirth.trim() === "" ||
        country.trim() === "" ||
        gender.trim() === ""
    ) {
        alert("Please fill in all fields.");
        return;
    }
    const formData = new FormData();
    formData.append("username", username);
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
    formData.append("dateOfBirth", dateOfBirth);
    formData.append("country", country);
    formData.append("gender", gender);
    if (confirm("Are you sure you want to proceed?")) {
        fetch("h    ttp://localhost:8080/updateUserInfo", {
                method: "POST",
                body: formData,
            })
            .then((response) => {
                if (response.ok) {
                    alert("User information updated successfully");
                    toggleDiv("user-info");
                } else {
                    return response.text().then((error) => {
                        alert(error);
                        console.error("User information update failed:", error);
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
                    messageBox.id = chat.UserId;
                    const chatUserPic = document.createElement("div");
                    chatUserPic.className = "chatUserPic";
                    chatUserPic.style.backgroundImage = `url(${chat.Image})`;
                    // Set border color based on online status
                    if (chat.Online) {
                        chatUserPic.style.border = " 3px solid rgb(74, 250, 58);";
                    } else {
                        chatUserPic.style.border = " 3px solid red";
                    }

                    const chatUserName = document.createElement("div");
                    chatUserName.className = "chatUserName";
                    const userNameP = document.createElement("p");
                    userNameP.textContent = chat.Username;
                    chatUserName.appendChild(userNameP);

                    const newMessageIcon = document.createElement("div");
                    newMessageIcon.className = "newMessageIcon";

                    messageBox.appendChild(chatUserPic);
                    messageBox.appendChild(chatUserName);
                    messageBox.appendChild(newMessageIcon);
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
setInterval(GetUserLoggedIn, 5000);

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
    const url = new URL("http://localhost:8080/userProfile");
    url.searchParams.append("user_id", userId);
    url.searchParams.append("case", caseString);

    fetch(url, {
            method: "GET",
        })
        .then((response) => response.json())
        .then((data) => {
            console.log(data);
            console.log(data.UserProfile.image_url); // Change imageURL to ImageURL
            const profileUsername = document.getElementById("profileUsername");
            profileUsername.innerHTML = data.UserProfile.username;

            // Calculate and display the age
            const birthDate = data.UserProfile.DateOfBirth;
            const age = calculateAge(birthDate);
            const profileAge = document.getElementById("userInfo");
            profileAge.innerHTML = `Age: ${age},  location: ${data.UserProfile.location}`;

            const userPic = document.getElementById("userPic");
            userPic.style.backgroundImage = `url(${data.UserProfile.image_url})`;
            const postsThElement = document.querySelector(
                "th[onclick=\"posts(-1,'Posts','Posts', event.target)\"]"
            );
            posts(-1, "Posts", "Posts", postsThElement);
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
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then((data) => {
            console.log(caseString);
            if (data && data.Posts) {
                displayPostOnProfile(data.Posts);
            } else {
                console.error("Invalid data format. Expected profileView with Posts.");
            }
        })
        .catch((error) => {
            console.error("Error fetching profile:", error);
        });
}

function displayPostOnProfile(Posts) {
    if (Array.isArray(Posts)) {
        createPost(Posts, "profileContent");
    } else {
        console.error("Invalid data format. Expected an array of posts.");
    }
}

function createPost(Posts, divName) {
    const homeNavigationContent = document.getElementById(divName);
    homeNavigationContent.innerHTML = "";

    Posts.forEach((post) => {
        let numOfLike = 0;
        let numOfDislike = 0;
        console.log(post);

        let liskIsClicked = false;
        let disliskIsClicked = false;
        // Assuming `post` is of type `structs.PostResponse`
        const reactions = post.reactions; // This should be an array of `structs.PostReactionResponse`
        // if (reactions.length > 0) {
        reactions.forEach((reaction) => {
            console.log("hi");
            if (reaction.type === "like") {
                liskIsClicked = reaction.did_react;

                numOfLike = reaction.count;
            } else if (reaction.type === "dislike") {
                numOfDislike = reaction.count;
                disliskIsClicked = reaction.did_react;
            }
        });
        // }

        const postBox = document.createElement("div");
        postBox.classList.add("postBox");
        postBox.setAttribute("onclick", `PostPageHandler(${JSON.stringify(post)})`);
        postBox.setAttribute("id", `${post.id}`);
        console.log(post.id);
        const postUserPic = document.createElement("div");
        postUserPic.classList.add("postUserPic");
        postBox.appendChild(postUserPic);

        const postTitle = document.createElement("div");
        postTitle.classList.add("postStyle");

        const titleElement = document.createElement("span");
        titleElement.classList.add("postTitle");
        titleElement.textContent = post.title;

        const postUserName = document.createElement("span");
        postUserName.classList.add("postUserName");

        postUserName.textContent = post.author.username;

        const postContent = document.createElement("span");
        postContent.classList.add("postContent");

        postContent.textContent = post.message;

        // Add rxn stuff
        const postReactions = document.createElement("span");
        postReactions.classList.add("postReaction");

        const postLikeIcone = document.createElement("button");
        postLikeIcone.classList.add("postLike");
        if (liskIsClicked) {
            postLikeIcone.classList.toggle("clicked", liskIsClicked);
        }

        const likeReactionCount = document.createElement("span");
        likeReactionCount.classList.add("reactionCount");

        likeReactionCount.textContent = numOfLike;

        const postDislikeIcone = document.createElement("button");
        postDislikeIcone.classList.add("postDislike");
        if (disliskIsClicked) {
            postDislikeIcone.classList.toggle("clicked", disliskIsClicked);
        }
        const dislikeReactionCount = document.createElement("span");
        dislikeReactionCount.classList.add("reactionCount");

        dislikeReactionCount.textContent = numOfDislike;

        postTitle.appendChild(titleElement);
        postTitle.appendChild(postUserName);
        postTitle.appendChild(postContent);

        postReactions.appendChild(dislikeReactionCount);
        postReactions.appendChild(postDislikeIcone);

        postReactions.appendChild(likeReactionCount);
        postReactions.appendChild(postLikeIcone);

        postBox.appendChild(postTitle);
        postBox.appendChild(postReactions);

        homeNavigationContent.appendChild(postBox);
        postLikeIcone.addEventListener("click", (event) => {
            // Prevent the click event on the button from bubbling up to the div
            event.stopPropagation();
            const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
            if (profileIcon.style.display === "block") {
                liskIsClicked = !liskIsClicked && !disliskIsClicked;
                postLikeIcone.classList.toggle("clicked", liskIsClicked);
                if (liskIsClicked) {
                    numOfLike++;
                    console.log(post.id);
                    AddReaction(1, post.id);
                } else {
                    numOfLike--;
                    deleteReaction(1, post.id);
                }
                likeReactionCount.textContent = numOfLike;
            }
        });

        postDislikeIcone.addEventListener("click", (event) => {
            // Prevent the click event on the button from bubbling up to the div
            event.stopPropagation();
            const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
            if (profileIcon.style.display === "block") {
                disliskIsClicked = !disliskIsClicked && !liskIsClicked;
                postDislikeIcone.classList.toggle("clicked", disliskIsClicked);
                if (disliskIsClicked) {
                    numOfDislike++;
                    AddReaction(2, post.id);
                    console.log("hi");
                } else {
                    numOfDislike--;
                    deleteReaction(2, post.id);
                }
                dislikeReactionCount.textContent = numOfDislike;
            }
        });
    });
}

async function deleteReaction(reaction, postId) {
    try {
        const response = await fetch(`/posts/${postId}/reactions/${reaction}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
        });
        if (!response.ok) {
            throw new Error("Failed to delete reaction");
        }
        const result = await response.json();
        return result;
    } catch (error) {
        console.error("Error deleting reaction:", error);
    }
}

function deleteReaction(reaction, postId) {
    const Form = new FormData();
    Form.append("reaction", reaction);
    Form.append("postId", postId);
    fetch(`/post/reaction/delete`, {
            method: "POST",
            body: Form,
        })
        .then((response) => {
            console.log("hi");

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        })

    .catch((error) => {
        console.error("Error removing reaction:", error);
    });
}

function AddReaction(reaction, postId) {
    console.log("hi");
    const Form = new FormData();
    Form.append("reaction", reaction);
    Form.append("postId", postId);
    fetch(`/posts/AddReactions`, {
            method: "POST",
            body: Form,
        })
        .then((response) => {
            console.log("hi");

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        })

    .catch((error) => {
        console.error("Error adding reaction:", error);
    });
}

function validateUserInfoForm() {
    const username = document.getElementById("ChangeUserName").value.trim();
    const firstName = document.getElementById("ChangeFirstName").value.trim();
    const lastName = document.getElementById("ChangeLastName").value.trim();
    const email = document.getElementById("ChangeEmail").value.trim();
    const dateOfBirth = document.getElementById("ChangeDOB").value.trim();

    if (!username) {
        alert("Username is required.");
        return false;
    } else if (username.length < 4 || username.length > 20) {
        alert("Username must be between 4 and 20 characters.");
        return false;
    }
    if (!firstName) {
        alert("First Name is required.");
        return false;
    } else if (firstName.length < 3 || firstName.length > 20) {
        alert("First Name must be between 3 and 20 characters.");
        return false;
    }
    if (!lastName) {
        alert("Last Name is required.");
        return false;
    } else if (lastName.length < 3 || lastName.length > 20) {
        alert("Last Name must be between 3 and 20 characters.");
        return false;
    }
    if (!email) {
        alert("Email is required.");
        return false;
    } else if (!isValidEmail(email)) {
        alert("Invalid email format.");
        return false;
    }
    if (!dateOfBirth) {
        alert("Date of Birth is required.");
        return false;
    } else if (dateOfBirth > new Date()) {
        alert("Invalid date of birth.");
        return false;
    }
    // Additional validation can be added here (e.g., email format, date format)
    return true;
}

var element = document.getElementById("UserChat");
element.scrollTop = element.scrollHeight;

// To always keep the div scrolled to the bottom
element.addEventListener("DOMSubtreeModified", function() {
    element.scrollTop = element.scrollHeight;
});