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
            if (data.name === "Guest" || data.name === "User") {
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
    fetch("http://localhost:8080/Moderator", {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            console.log("Fetched data:", data); // Log the fetched data for debugging
            const moderatorsList = document.querySelector("#list-moderators ul");
            moderatorsList.innerHTML = ""; // Clear existing list items
            if (Array.isArray(data)) {
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
            } else {
                const paregraph = document.createElement("p");
                paregraph.textContent = "No moderators found.";
                moderatorsList.appendChild(paregraph);
            }
        })
        .catch((error) => {
            console.error("Error fetching moderators:", error);
        });
    toggleDiv("list-moderators");
}

function RemoveModerator(id) {
    confirmAction();
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
    const promotionRequestsTable = document.querySelector("#promotion-requests table");
    promotionRequestsTable.innerHTML = ""; // Clear existing table rows
    fetch("http://localhost:8080/PromotionRequests", {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.text(); // Read the response as text
        })
        .then((text) => {
            try {
                const data = JSON.parse(text); // Attempt to parse the text as JSON

                data.forEach((request) => {
                    const row = document.createElement("tr");

                    const usernameCell = document.createElement("td");
                    usernameCell.textContent = request.Username;
                    usernameCell.dataset.id = request.Id; // Set the moderator ID as a data attribute
                    row.appendChild(usernameCell);

                    const approveButtonCell = document.createElement("td");
                    const approveButton = document.createElement("button");

                    approveButton.classList.add("show-button");
                    approveButton.textContent = "Show";
                    approveButton.setAttribute("onclick", `ShowPromotion(${request.Id})`);

                    approveButtonCell.appendChild(approveButton);
                    row.appendChild(approveButtonCell);

                    const rejectButtonCell = document.createElement("td");
                    const rejectButton = document.createElement("button");

                    rejectButton.classList.add("reject-button");
                    rejectButton.textContent = "Reject";
                    rejectButton.setAttribute("onclick", `RejectPromotion(${request.UserId})`);

                    rejectButtonCell.appendChild(rejectButton);
                    row.appendChild(rejectButtonCell);

                    promotionRequestsTable.appendChild(row);
                });
            } catch (error) {
                const paregraph = document.createElement("p");
                paregraph.textContent = "No requests found.";
                promotionRequestsTable.appendChild(paregraph);
            }
        })
        .catch((error) => {
            console.error("Error fetching promotion requests:", error);
        });
    toggleDiv("promotion-requests");
}

function ShowPromotion(Id) {
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
            return response.text(); // Read the response as text
        })
        .then((text) => {
            try {
                const data = JSON.parse(text); // Attempt to parse the text as JSON
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
                // Handle the response data here
            } catch (error) {
                console.error("Error parsing JSON:", error);
                console.log("Response text:", text); // Log the response text for debugging
            }
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
    fetch("http://localhost:8080/getUserInfo", {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.text(); // Get response as text
        })
        .then((text) => {
            const data = JSON.parse(text); // Parse the text as JSON
            if (data) {
                document.getElementById("ChangeUserName").value = data.Username || "";
                document.getElementById("ChangeFirstName").value = data.FirstName || "";
                document.getElementById("ChangeLastName").value = data.LastName || "";
                document.getElementById("ChangeEmail").value = data.Email || "";
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
