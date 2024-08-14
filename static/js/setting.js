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
                        li.innerHTML.includes("Manage Categories") ||
                        li.innerHTML.includes("Manage Reports")
                    ) {
                        li.style.display = "none";
                    }
                });
            } else if (data.name === "Admin") {
                document.querySelectorAll("#settingList li").forEach((li) => {
                    if (
                        li.innerHTML.includes("List of Moderators") ||
                        li.innerHTML.includes("Promotion Requests") ||
                        li.innerHTML.includes("Manage Categories") ||
                        li.innerHTML.includes("Manage Reports")
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
                        li.innerHTML.includes("Request to be Moderator") ||
                        li.innerHTML.includes("Manage Reports")
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
    console.log("Removing moderator with ID:", id);
    // Create the request body
    const formData = new FormData();
    formData.append("Id", id);

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
                if (category.Name != "") {
                    appendCategoryToList(category.Name, category.Id);
                }
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
    const image = document.getElementById("ChangeUploadedImage").files[0];

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
    formData.append("image", image);

    if (confirm("Are you sure you want to proceed?")) {
        fetch("http://localhost:8080/updateUserInfo", {
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

function ManageReports() {
    const reportRequestsList = document.querySelector("#Report-requests ul");
    reportRequestsList.innerHTML = ""; // Clear existing list items

    fetch("http://localhost:8080/Reports", {
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
                    const listItem = document.createElement("li");

                    const reporterUsernameSpan = document.createElement("span");
                    reporterUsernameSpan.classList.add("ReporterUsername");
                    reporterUsernameSpan.textContent = request.reporter_username;
                    listItem.appendChild(reporterUsernameSpan);

                    const reportedUsernameSpan = document.createElement("span");
                    reportedUsernameSpan.classList.add("ReportedUsername");
                    reportedUsernameSpan.textContent = request.reported_username;
                    listItem.appendChild(reportedUsernameSpan);

                    const showButton = document.createElement("button");
                    showButton.classList.add("show-button");
                    showButton.textContent = "Show";
                    showButton.setAttribute(
                        "onclick",
                        `showReportRequest(${JSON.stringify(request)})`
                    );
                    listItem.appendChild(showButton);

                    const rejectButton = document.createElement("button");
                    rejectButton.classList.add("reject-button");
                    rejectButton.textContent = "Reject";
                    rejectButton.setAttribute(
                        "onclick",
                        `reportUserOrPost(${JSON.stringify(request)}, "Report Rejected", true)`
                    );
                    listItem.appendChild(rejectButton);

                    reportRequestsList.appendChild(listItem);
                });
            } catch (error) {
                const paragraph = document.createElement("p");
                paragraph.textContent = "No requests found.";
                reportRequestsList.appendChild(paragraph);
            }
        })
        .catch((error) => {
            console.error("Error fetching promotion requests:", error);
        });

    toggleDiv("Report-requests");
}

function showReportRequest(report) {
    toggleDiv("Report-requests");

    console.log("Report:", report); // Debugging line

    const reportDetailsDiv = document.getElementById("ShowReportDetails");
    reportDetailsDiv.innerHTML = ""; // Clear existing content

    const ul = document.createElement("ul");

    const fields = [
        { title: "Reporter Username", info: report.reporter_username },
        { title: "Reported Username", info: report.reported_username },
        { title: "Time", info: new Date(report.time).toLocaleString() },
        { title: "Reason", info: report.reason },
    ];
    const replyButton = document.createElement("button");
    replyButton.textContent = "Report";
    replyButton.classList.add("reply-button");
    replyButton.disabled = true; // Disable the button initially

    console.log(report);
    console.log(report);
    if (report.is_post_reported) {
        console.log;
        fields.push({ title: "Reported Post Title", info: report.reported_post_title });
        replyButton.textContent = "Delete Post";
    }

    fields.forEach((field) => {
        const li = document.createElement("li");

        const titleSpan = document.createElement("span");
        titleSpan.textContent = field.title + ": ";
        titleSpan.classList.add("title-span");

        const infoSpan = document.createElement("span");
        infoSpan.textContent = field.info;
        infoSpan.classList.add("info-span");

        li.appendChild(titleSpan);
        li.appendChild(infoSpan);
        ul.appendChild(li);
    });

    // Add reply label, input, and button in a new li
    const replyLi = document.createElement("li");

    const replyLabel = document.createElement("label");
    replyLabel.textContent = "Response :";
    replyLabel.classList.add("reply-label");

    const replyInput = document.createElement("textarea");
    replyInput.id = "replyInput";
    replyInput.classList.add("reply-input");
    replyInput.onfocus = replyInput.oninput = function() {
        replyButton.disabled = replyInput.value.trim() === "";
    };
    replyButton.onclick = function() {
        const replyText = document.getElementById("replyInput").value;
        reportUserOrPost(report, replyText, false);
        document.getElementById("replyInput").innerHTML = "";
        console.log("Reply:", replyText);
        // Add your reply handling logic here
    };
    replyLi.appendChild(replyLabel);
    replyLi.appendChild(replyInput);
    ul.appendChild(replyLi);

    const rejectButton = document.createElement("button");
    rejectButton.textContent = "Reject";
    rejectButton.classList.add("rejectReport-button");
    rejectButton.onclick = function() {
        reportUserOrPost(report, "Report Rejected", true);
    };
    reportDetailsDiv.appendChild(ul);
    reportDetailsDiv.appendChild(replyButton);
    reportDetailsDiv.appendChild(rejectButton);

    toggleDiv("ReportDetails");
}

function reportUserOrPost(report, replyText, rejected) {
    const reportInfo = {
        report_id: report.id,
        response: replyText,
    };

    fetch("/updateReport", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(reportInfo),
        })
        .then((response) => {
            if (!response.ok) {
                return response.text().then((errorText) => {
                    throw new Error("Failed to update report: " + errorText);
                });
            }
            return response.json();
        })
        .then((result) => {
            console.log(result.message);

            if (report.is_post_reported && !rejected) {
                deletePost(report.id);
            } else if (!rejected) {
                banUser(report.reported_id);
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
    toggleDiv("Report-requests");
}

function banUser(userId) {
    fetch(`/banUser?userId=${userId}`, {
            method: "GET", // or 'POST' if you prefer
            headers: {
                "Content-Type": "application/json",
            },
        })
        .then((response) => {
            if (!response.ok) {
                return response.text().then((errorText) => {
                    throw new Error(`Failed to ban user: ${errorText}`);
                });
            }
            return response.text();
        })
        .then((result) => {
            console.log(result); // "User banned successfully"
            alert(result);
        })
        .catch((error) => {
            console.error("Error:", error);
            alert(`Error: ${error.message}`);
        });
}