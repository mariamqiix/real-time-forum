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
                try {
                    // Clean the response text by removing the trailing 'null'
                    const cleanedText = text.replace(/null$/, "").trim();
                    const data1 = JSON.parse(cleanedText); // Parse the cleaned text as JSON
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
                throw new Error(`HTTP error! status: ${response.status}`);
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

function displayPostPage(data) {
    const postPageContent = document.querySelector(".postPageContent");
    postPageContent.innerHTML = ""; // Clear any existing content

    const postPageTime = document.createElement("span");
    postPageTime.classList.add("postPageTime");
    postPageTime.textContent = data.Post.time;
    postPageContent.appendChild(postPageTime);

    const postPageTitle = document.createElement("h1");
    postPageTitle.classList.add("postPageTitle");
    postPageTitle.textContent = data.Post.title;
    postPageContent.appendChild(postPageTitle);

    const postMessage = document.createElement("p");
    postMessage.textContent = data.Post.message;
    postPageContent.appendChild(postMessage);

    const postPageCategories = document.createElement("div");
    postPageCategories.classList.add("postPageCategories");
    console.log(data.Post.categories);
    postPageCategories.textContent = Array.isArray(data.Post.categories) ?
        data.Post.categories.map((category) => category.name).join(", ") // Extract and join category names
        :
        "";
    postPageContent.appendChild(postPageCategories);

    const postPageReaction = document.querySelector(".postPageReaction");
    postPageReaction.innerHTML = ""; // Clear any existing content
    var numOfLike = 0;
    var numOfDislike = 0;
    console.log(data);
    const reactions = data.Post.reactions; // This should be an array of `structs.PostReactionResponse`

    reactions.forEach((reaction) => {
        if (reaction.type === "like") {
            // liskIsClicked = reaction.did_react;

            numOfLike = reaction.count;
        } else if (reaction.type === "dislike") {
            numOfDislike = reaction.count;
            // disliskIsClicked = reaction.did_react;
        }
    });
    const postLike = document.createElement("div");
    postLike.classList.add("postLike");
    postPageReaction.appendChild(postLike);

    const likeCount = document.createElement("div");
    likeCount.classList.add("reactionCount");

    likeCount.textContent = numOfLike;
    postPageReaction.appendChild(likeCount);

    const postDislike = document.createElement("div");
    postDislike.classList.add("postDislike");
    postPageReaction.appendChild(postDislike);

    const dislikeCount = document.createElement("div");
    dislikeCount.classList.add("reactionCount");
    dislikeCount.textContent = numOfDislike;
    postPageReaction.appendChild(dislikeCount);

    // data.comments.forEach((comment) => {
    //     const postPageComment = document.createElement("div");
    //     postPageComment.classList.add("postPageComment");
    //     postPageComment.textContent = comment.text;

    //     const postPageCommentUser = document.createElement("span");
    //     postPageCommentUser.classList.add("postPageCommentUser");
    //     postPageCommentUser.textContent = comment.user;
    //     postPageComment.appendChild(postPageCommentUser);

    //     const postPageCommentTime = document.createElement("span");
    //     postPageCommentTime.classList.add("postPageCommentTime");
    //     postPageCommentTime.textContent = `- ${comment.time}`;
    //     postPageComment.appendChild(postPageCommentTime);

    //     const postPageCommentReaction = document.createElement("div");
    //     postPageCommentReaction.classList.add("postPageCommentReaction");

    //     const commentLike = document.createElement("div");
    //     commentLike.classList.add("postLike");
    //     postPageCommentReaction.appendChild(commentLike);

    //     const commentLikeCount = document.createElement("div");
    //     commentLikeCount.classList.add("reactionCount");
    //     commentLikeCount.textContent = comment.reactions.like;
    //     postPageCommentReaction.appendChild(commentLikeCount);

    //     const commentDislike = document.createElement("div");
    //     commentDislike.classList.add("postDislike");
    //     postPageCommentReaction.appendChild(commentDislike);

    //     const commentDislikeCount = document.createElement("div");
    //     commentDislikeCount.classList.add("reactionCount");
    //     commentDislikeCount.textContent = comment.reactions.dislike;
    //     postPageCommentReaction.appendChild(commentDislikeCount);

    //     postPageComment.appendChild(postPageCommentReaction);

    //     postPageContent.appendChild(postPageComment);
    // });
}