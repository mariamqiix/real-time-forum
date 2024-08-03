function toggleDiv(divName) {
    const contentDiv = document.getElementById(divName);
    const overlayDiv = document.getElementById("overlay");

    if (contentDiv.style.display === "none") {
        // Show the content div
        overlayDiv.style.display = "block";
        contentDiv.style.display = "block";
    } else {
        // Hide the content div
        contentDiv.style.display = "none";
        overlayDiv.style.display = "none";
    }
}

function toggleLogout(className) {
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
            if (text != "null") {
                logout();
            } else {
                toggleVisibility(className);
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function toggleVisibility(className) {
    var contentDivs = document.getElementsByClassName("navigationContent");

    for (var i = 0; i < contentDivs.length; i++) {
        var div = contentDivs[i];

        if (div.classList.contains(className)) {
            div.style.display = "block";
        } else {
            div.style.display = "none";
        }
    }
}

function OpenMesages(username) {
    const messagesBoxDiv = document.getElementById("messagesBoxDiv");
    const msgDiv = document.getElementById("msgDiv");
    const sendMessageButton = document.getElementById("sendMessage");

    if (messagesBoxDiv.style.display === "none") {
        document.getElementById("messagesTitle").innerHTML = "Messages";
        messagesBoxDiv.style.display = "block";
        msgDiv.style.display = "none";
    } else {
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
                if (text != "null") {
                    messagesBoxDiv.style.display = "none";
                    msgDiv.style.display = "block";
                    document.getElementById("messagesTitle").innerHTML = username;
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }

    // Set the onclick attribute to call SendMessage with the username
    sendMessageButton.setAttribute("onclick", `SendMessage('${username}')`);
}