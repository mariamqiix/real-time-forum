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

            // Loop through each checked checkbox and uncheck it
            categoryCheckboxes.forEach(function(checkbox) {
                checkbox.checked = false;
            });
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
                handleErrorResponse(response);
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
                handleErrorResponse(response);
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


    // Add the change event listener
    checkbox.addEventListener("change", function() {
        validateCreatePost();
    });


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

function PostPageHandler(data) {
    const ReporBtutton = document.getElementById("Report-button");
    // Send the form data to the Go server
    fetch("http://localhost:8080/userType", {
            method: "GET",
        })
        .then((response) => {
            if (!response.ok) {
                ReporBtutton.style.display = "none";
            }
            return response.json();
        })
        .then((type) => {
            if (type.name === "Guest" || type.name === "User") {
                ReporBtutton.style.display = "none";
            } else {
                ReporBtutton.style.display = "block";
            }
        })
        .catch((error) => {
            ReporBtutton.style.display = "none";
        });

    const editPostButton = document.getElementById("editPost-button");
    toggleVisibility("postPage");
    if (editPostButton) {
        fetch("http://localhost:8080/user", {
                method: "GET",
            })
            .then((response) => {
                if (!response.ok) {
                    handleErrorResponse(response);
                }
                return response.text(); // Get response as text
            })
            .then((text) => {
                try {
                    // Clean the response text by removing the trailing 'null'
                    const cleanedText = text.replace(/null$/, "").trim();
                    const data1 = JSON.parse(cleanedText); // Parse the cleaned text as JSON
                    if (data1 === data.author.username) {
                        console.log("User is the author of the post.");
                        editPostButton.style.display = "block";
                        editPostButton.setAttribute("onclick", `editPost(${data.id})`);
                        ReporBtutton.style.display = "none";
                    } else {
                        editPostButton.style.display = "none";
                        ReporBtutton.setAttribute("onclick", `ReportPost(${data.id},false,true)`);
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

    PostInfo(data.id);
}

function EditPostHandler() {}

function PostInfo(id) {
    const formData = new FormData();
    formData.append("post_id", id);

    fetch(`http://localhost:8080/post/${id}`, {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (!response.ok) {
                handleErrorResponse(response);
            }
            return response.json();
        })
        .then((data) => {
            displayPostPage(data);
        })
        .catch((error) => {
            console.error("Error fetching post:", error);
        });
}
// // To always keep the div scrolled to the bottom
// element.addEventListener("DOMSubtreeModified", function() {
//     element.scrollTop = element.scrollHeight;
// });
function formatDateTime(dateTimeString) {
    const date = new Date(dateTimeString);
    const options = { hour: "2-digit", minute: "2-digit" };
    const time = date.toLocaleTimeString([], options);
    const day = date.getDate();
    const month = date.toLocaleString("default", { month: "long" });
    const year = date.getFullYear();

    return `${time} - ${day}/${month}/${year}`;
}

function displayPostPage(data) {
    const postPageContent = document.querySelector(".postPageContent");
    postPageContent.innerHTML = ""; // Clear any existing content
    let liskIsClicked = false;
    let disliskIsClicked = false;
    const replayButton = document.getElementById("SendReplay");
    replayButton.setAttribute(
        "onclick",
        `SendReplay(${data.Post.id}, '${data.Post.author.username}')`
    );

    const postPageTime = document.createElement("span");
    postPageTime.classList.add("postPageTime");
    postPageTime.textContent = formatDateTime(data.Post.created_at);
    postPageContent.appendChild(postPageTime);

    const postPageUsername = document.createElement("span");
    postPageUsername.classList.add("postPageUsername");
    postPageUsername.textContent = data.Post.author.username;
    postPageContent.appendChild(postPageUsername);

    const postPageTitle = document.createElement("h1");
    postPageTitle.classList.add("postPageTitle");
    postPageTitle.textContent = data.Post.title;
    postPageContent.appendChild(postPageTitle);

    const postMessage = document.createElement("p");
    postMessage.textContent = data.Post.message;
    postPageContent.appendChild(postMessage);

    const postPageCategories = document.createElement("div");
    postPageCategories.classList.add("postPageCategories");
    postPageCategories.textContent = Array.isArray(data.Post.categories) ?
        data.Post.categories.map((category) => category.name).join(", ") // Extract and join category names
        :
        "";
    postPageContent.appendChild(postPageCategories);

    const postPageReaction = document.querySelector(".postPageReaction");
    postPageReaction.innerHTML = ""; // Clear any existing content
    var numOfLike = 0;
    var numOfDislike = 0;
    const reactions = data.Post.reactions; // This should be an array of `structs.PostReactionResponse`

    reactions.forEach((reaction) => {
        if (reaction.type === "like") {
            liskIsClicked = reaction.did_react;

            numOfLike = reaction.count;
        } else if (reaction.type === "dislike") {
            numOfDislike = reaction.count;
            disliskIsClicked = reaction.did_react;
        }
    });

    const likeCount = document.createElement("div");
    likeCount.classList.add("reactionCount");
    likeCount.textContent = numOfLike;
    postPageReaction.appendChild(likeCount);

    const postLike = document.createElement("div");
    postLike.classList.add("postLike");
    postPageReaction.appendChild(postLike);
    if (liskIsClicked) {
        postLike.classList.toggle("clicked", liskIsClicked);
    }

    const dislikeCount = document.createElement("div");
    dislikeCount.classList.add("reactionCount");
    dislikeCount.textContent = numOfDislike;
    postPageReaction.appendChild(dislikeCount);

    const postDislike = document.createElement("div");
    postDislike.classList.add("postDislike");
    postPageReaction.appendChild(postDislike);
    if (disliskIsClicked) {
        postDislike.classList.toggle("clicked", disliskIsClicked);
    }
    postLike.addEventListener("click", (event) => {
        // Prevent the click event on the button from bubbling up to the div
        event.stopPropagation();
        const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
        if (profileIcon.style.display === "block") {
            liskIsClicked = !liskIsClicked && !disliskIsClicked;
            postLike.classList.toggle("clicked", liskIsClicked);
            if (liskIsClicked) {
                numOfLike++;
                AddReaction(1, data.Post.id);
            } else if (!disliskIsClicked) {
                numOfLike--;
                deleteReaction(1, data.Post.id);
            }
            likeCount.textContent = numOfLike;
        }
    });

    postDislike.addEventListener("click", (event) => {
        // Prevent the click event on the button from bubbling up to the div
        event.stopPropagation();
        const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
        if (profileIcon.style.display === "block") {
            disliskIsClicked = !disliskIsClicked && !liskIsClicked;
            postDislike.classList.toggle("clicked", disliskIsClicked);
            if (disliskIsClicked) {
                numOfDislike++;
                AddReaction(2, data.Post.id);
            } else if (!liskIsClicked) {
                numOfDislike--;
                deleteReaction(2, data.Post.id);
            }
            dislikeCount.textContent = numOfDislike;
        }
    });
    const postPageCommentsContent = document.getElementById("comments");
    postPageCommentsContent.innerHTML = ""; // Clear any existing content
    data.Comments.forEach((comment) => {
        createComments(comment);
    });
}

function createComments(comment) {
    let liskIsClicked = false;
    let disliskIsClicked = false;
    const postPageContent = document.getElementById("comments");

    const postPageComment = document.createElement("div");
    postPageComment.classList.add("postPageComment");
    postPageComment.textContent = comment.message;

    const postPageCommentUser = document.createElement("span");
    postPageCommentUser.classList.add("postPageCommentUser");
    postPageCommentUser.textContent = `    ${comment.author.username}`;
    postPageComment.appendChild(postPageCommentUser);

    const postPageCommentTime = document.createElement("span");
    postPageCommentTime.classList.add("postPageCommentTime");
    postPageCommentTime.textContent = `- ${comment.created_at}`;
    postPageComment.appendChild(postPageCommentTime);
    const postPageCommentReaction = document.createElement("div");
    postPageCommentReaction.classList.add("postPageCommentReaction");
    // const commentLike = document.createElement("div");
    // commentLike.classList.add("postLike");
    // postPageCommentReaction.appendChild(commentLike);
    // const commentLikeCount = document.createElement("div");
    // commentLikeCount.classList.add("reactionCount");
    // commentLikeCount.textContent = comment.reactions.like;
    // postPageCommentReaction.appendChild(commentLikeCount);
    // const commentDislike = document.createElement("div");
    // commentDislike.classList.add("postDislike");
    // postPageCommentReaction.appendChild(commentDislike);
    // const commentDislikeCount = document.createElement("div");
    // commentDislikeCount.classList.add("reactionCount");
    // commentDislikeCount.textContent = comment.reactions.dislike;
    // postPageCommentReaction.appendChild(commentDislikeCount);
    // postPageComment.appendChild(postPageCommentReaction);

    var numOfLike = 0;
    var numOfDislike = 0;
    const reactions = comment.reactions; // This should be an array of `structs.PostReactionResponse`

    reactions.forEach((reaction) => {
        if (reaction.type === "like") {
            liskIsClicked = reaction.did_react;

            numOfLike = reaction.count;
        } else if (reaction.type === "dislike") {
            numOfDislike = reaction.count;
            disliskIsClicked = reaction.did_react;
        }
    });

    const likeCount = document.createElement("div");
    likeCount.classList.add("reactionCount");
    likeCount.textContent = numOfLike;
    postPageCommentReaction.appendChild(likeCount);

    const postLike = document.createElement("div");
    postLike.classList.add("postLike");
    postPageCommentReaction.appendChild(postLike);
    if (liskIsClicked) {
        postLike.classList.toggle("clicked", liskIsClicked);
    }

    const dislikeCount = document.createElement("div");
    dislikeCount.classList.add("reactionCount");
    dislikeCount.textContent = numOfDislike;
    postPageCommentReaction.appendChild(dislikeCount);

    const postDislike = document.createElement("div");
    postDislike.classList.add("postDislike");
    postPageCommentReaction.appendChild(postDislike);
    if (disliskIsClicked) {
        postDislike.classList.toggle("clicked", disliskIsClicked);
    }
    postPageComment.appendChild(postPageCommentReaction);
    postLike.addEventListener("click", (event) => {
        // Prevent the click event on the button from bubbling up to the div
        event.stopPropagation();
        const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
        if (profileIcon.style.display === "block") {
            liskIsClicked = !liskIsClicked && !disliskIsClicked;
            postLike.classList.toggle("clicked", liskIsClicked);
            if (liskIsClicked) {
                numOfLike++;
                AddReaction(1, comment.id);
            } else if (!disliskIsClicked) {
                numOfLike--;
                deleteReaction(1, comment.id);
            }
            likeCount.textContent = numOfLike;
        }
    });

    postDislike.addEventListener("click", (event) => {
        // Prevent the click event on the button from bubbling up to the div
        event.stopPropagation();
        const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
        if (profileIcon.style.display === "block") {
            disliskIsClicked = !disliskIsClicked && !liskIsClicked;
            postDislike.classList.toggle("clicked", disliskIsClicked);
            if (disliskIsClicked) {
                numOfDislike++;
                AddReaction(2, comment.id);
            } else if (!liskIsClicked) {
                numOfDislike--;
                deleteReaction(2, comment.id);
            }
            dislikeCount.textContent = numOfDislike;
        }
    });
    postPageContent.appendChild(postPageComment);
}

function SendReplay(postId, Username) {
    const message = document.getElementById("PostReplay").value;
    if (message.trim() != "") {
        const formData = new FormData();
        formData.append("title", `Replay to (${Username})`);
        formData.append("content", message);
        formData.append("post_id", postId);

        fetch(`/post/comment`, {
                method: "POST",
                body: formData,
            })
            .then((response) => {
                if (response.ok) {
                    // Handle success (e.g., clear the input field, update the comments section)
                    document.getElementById("PostReplay").value = "";
                    return response.json();
                    // Optionally, you can refresh the comments section to show the new comment
                } else {
                    // Handle error
                    console.error("Error:", data.error);
                }
            })
            .then((data) => {
                createComments(data.Post);
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }
}

function toggleNewPostForm() {
    const EditBtn = document.getElementById("EditBtn");
    const deletePostBtn = document.getElementById("deletePostBtn");
    EditBtn.style.display = "none";
    deletePostBtn.style.display = "none";
    const submitBtn = document.getElementById("submitBtn");
    submitBtn.style.display = "block";
    // Populate the form fields with the data
    document.getElementById("newPostTitle").value = "";
    document.getElementById("topic").value = "";
    submitBtn.style.display = "block";
    // Select the categories
    const categoriesList = document.querySelectorAll('.categoriesList input[type="checkbox"]');
    categoriesList.forEach((checkbox) => {
        checkbox.checked = false;
    });
    toggleDiv("newPostForm");
}

function editPost(Id) {
    toggleDiv("newPostForm");
    fetch(`http://localhost:8080/post/${Id}/edit`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
        .then((response) => {
            if (!response.ok) {
                handleErrorResponse(response);
            }
            return response.json();
        })
        .then((data) => {
            // Populate the form fields with the data
            document.getElementById("newPostTitle").value = data.Post.title;
            document.getElementById("topic").value = data.Post.message;
            const EditBtn = document.getElementById("EditBtn");
            const deletePostBtn = document.getElementById("deletePostBtn");
            EditBtn.style.display = "block";
            deletePostBtn.style.display = "block";
            const submitBtn = document.getElementById("submitBtn");
            submitBtn.style.display = "none";
            // Select the categories
            const categoriesList = document.querySelectorAll(
                '.categoriesList input[type="checkbox"]'
            );
            categoriesList.forEach((checkbox) => {
                checkbox.checked = data.Post.categories.some(
                    (category) => category.name === checkbox.name
                );
            });
            deletePostBtn.setAttribute("onclick", `deletePost(${data.Post.id})`);

            // Set the submit button's onclick function
            document.getElementById("EditBtn").onclick = function() {
                updatePost(data.Post.id);
            };
            // Toggle the newPostForm div
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function updatePost(id) {
    const title = document.getElementById("newPostTitle").value;
    const topic = document.getElementById("topic").value;
    const categoryCheckboxes = document.querySelectorAll(
        '.categoriesList input[type="checkbox"]:checked'
    );
    const selectedCategories = Array.from(categoryCheckboxes).map((checkbox) => checkbox.value);
    if (selectedCategories.length === 0) {
        alert("Please select at least one category.");
        return;
    }

    const formData = new FormData();
    formData.append("title", title);
    formData.append("content", topic);
    formData.append("categories", JSON.stringify(selectedCategories));

    fetch(`http://localhost:8080/post/${id}/edit`, {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                console.log("Post updated successfully");
                HomePageRequest();
                toggleDiv("newPostForm");
            } else {
                const error = response.text();
                console.error("Post update failed:", error);
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function deletePost(id) {
    fetch(`http://localhost:8080/post/${id}/delete`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
        })
        .then((response) => {
            if (response.ok) {
                console.log("Post deleted successfully");
                toggleNewPostForm();
                HomePageRequest();
            } else {
                const error = response.text();
                console.error("Post deletion failed:", error);
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function ReportPost(id, user, post) {
    toggleDiv("ReportDiv");
    const ReportBtn = document.getElementById("ReportDiv-button");
    ReportBtn.onclick = function() {
        ReportPostHandler(id, user, post);
    };
}

function ReportPostHandler(id, user, post) {
    const reportDescription = document.getElementById("ReportDescription").value;
    console.log(reportDescription);
    if (post) {
        const postReportRequest = {
            post_id: id,
            reason: reportDescription,
        };

        fetch(`http://localhost:8080/post/${id}/report`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postReportRequest),
            })
            .then((response) => {
                if (response.ok) {
                    const reportDescription2 = document.getElementById("ReportDescription");
                    reportDescription2.innerHTML = "";
                    toggleDiv("ReportDiv");
                    HomePageRequest();
                    console.log("Post reported successfully");
                } else {
                    response.text().then((error) => {
                        console.error("Post reporting failed:", error);
                    });
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    } else if (user) {
        const postReportRequest = {
            username: id,
            reason: reportDescription,
        };
        console.log(user);
        fetch(`http://localhost:8080/user/${id}/report`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postReportRequest),
            })
            .then((response) => {
                if (response.ok) {
                    const reportDescription2 = document.getElementById("ReportDescription");
                    reportDescription2.innerHTML = "";
                    toggleDiv("ReportDiv");
                    HomePageRequest();
                    console.log("Post reported successfully");
                } else {
                    response.text().then((error) => {
                        console.error("user reporting failed:", error);
                    });
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }
}

var chatContainer = document.getElementById("UserChat");

// Function to scroll the chat container to the bottom
function scrollToBottom() {
    chatContainer.scrollTop = chatContainer.scrollHeight;
}

// Call the function to scroll to the bottom when the DOM is fully loaded
document.addEventListener("DOMContentLoaded", function() {
    scrollToBottom();
});

// Auto scroll to bottom when content is added
chatContainer.addEventListener("DOMSubtreeModified", function() {
    scrollToBottom();
});

function handleErrorResponse(response) {
    if (!response.ok) {
        response
            .json()
            .then((error) => {
                console.error("Error:", error.message);
                // if (error.user) {
                //     alert("User Info:", error.user);
                // }
                // Display the error message to the user
                alert(`Error: ${error.message}`);
            })
            .catch((err) => {
                alert("Failed to parse error response:", err);
            });
    }
}