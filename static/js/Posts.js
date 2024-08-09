function createPost(Posts, divName) {
    const homeNavigationContent = document.getElementById(divName);
    homeNavigationContent.innerHTML = "";

    Posts.forEach((post) => {
        let numOfLike = 0;
        let numOfDislike = 0;
        let liskIsClicked = false;
        let disliskIsClicked = false;
        // Assuming `post` is of type `structs.PostResponse`
        const reactions = post.reactions; // This should be an array of `structs.PostReactionResponse`
        // if (reactions.length > 0) {
        reactions.forEach((reaction) => {
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
                    AddReaction(1, post.id);
                } else if (!disliskIsClicked) {
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
                } else if (!liskIsClicked) {
                    numOfDislike--;
                    deleteReaction(2, post.id);
                }
                dislikeReactionCount.textContent = numOfDislike;
            }
        });
    });
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
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        })

    .catch((error) => {
        console.error("Error removing reaction:", error);
    });
}

function AddReaction(reaction, postId) {
    const Form = new FormData();
    Form.append("reaction", reaction);
    Form.append("postId", postId);
    fetch(`/posts/AddReactions`, {
            method: "POST",
            body: Form,
        })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        })

    .catch((error) => {
        console.error("Error adding reaction:", error);
    });
}