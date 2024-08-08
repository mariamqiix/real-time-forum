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

function OpenMesages(username, id) {
    const messagesBoxDiv = document.getElementById("messagesBoxDiv");
    const msgDiv = document.getElementById("msgDiv");
    const sendMessageButton = document.getElementById("sendMessage");
    const mailIcon = document.getElementById("mailIcon");
    const msgType = document.getElementById("msgType");

    if (messagesBoxDiv.style.display === "none") {
        document.getElementById("messagesTitle").innerHTML = "Messages";
        messagesBoxDiv.style.display = "block";
        msgDiv.style.display = "none";
        msgDiv.setAttribute("data-id", "");
        mailIcon.style.display = "none";
    } else {
        fetch("http://localhost:8080/user", {
            method: "GET",
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.text(); // Get response as text
            })
            .then((text) => {
                if (text != "null") {
                    sendMessageButton.setAttribute("onclick", `SendMessage('${username}')`);
                    msgType.setAttribute("oninput", `notifyTyping('${username}')`);
                    mailIcon.style.display = "block";
                    mailIcon.onclick = function () {
                        OpenMesages();
                    };
                    messagesBoxDiv.style.display = "none";
                    msgDiv.style.display = "block";
                    document.getElementById("messagesTitle").innerHTML = username;
                    msgDiv.setAttribute("data-id", id);
                    const messageBox = document.getElementById(id);

                    if (messageBox) {
                        const newMessageIcon = messageBox.querySelector(".newMessageIcon");
                        if (newMessageIcon) {
                            newMessageIcon.style.backgroundColor = "#fbd998";
                        }
                    }
                    // Create the request body
                    const formData = new FormData();
                    formData.append("id", id);
                    const chats = document.getElementById("UserChat");
                    chats.scrollTop = chats.scrollHeight;

                    chats.innerHTML = "";
                    fetch(`http://localhost:8080/messages`, {
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
                            const OldmsgDiv = document.getElementById("msgDiv");

                            console.log(data);
                            data.forEach((messageData) => {
                                const chatDiv = document.createElement("div");
                                chatDiv.className = "fullMessage";

                                // Create the message div
                                const msgDiv = document.createElement("div");
                                msgDiv.className = "msg";
                                msgDiv.textContent = messageData.Messag;
                                // Create the new chat div
                                const id = OldmsgDiv.getAttribute("data-id");

                                // Create the message time div
                                const msgTimeDiv = document.createElement("div");
                                msgTimeDiv.className = "msgTime";
                                msgTimeDiv.textContent = new Date(
                                    messageData.Time
                                ).toLocaleTimeString();

                                if (id == messageData.ReceiverId) {
                                    chatDiv.classList.add("receiver");
                                    msgDiv.style.backgroundColor = "#fbd998";
                                    chatDiv.appendChild(msgTimeDiv);
                                    chatDiv.appendChild(msgDiv);
                                } else {
                                    chatDiv.classList.add("sender");
                                    msgDiv.classList.add("msg-right");
                                    chatDiv.appendChild(msgDiv);
                                    chatDiv.appendChild(msgTimeDiv);
                                }
                                chats.appendChild(chatDiv);

                                sendMessageButton.setAttribute(
                                    "onclick",
                                    `SendMessage('${username}')`
                                );
                            });
                        })
                        .catch((error) => {
                            console.error("Error:", error);
                        });
                }
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }

    // Set the onclick attribute to call SendMessage with the username
}

function toggleInputBox() {
    var inputBox = document.getElementById("inputBox");
    inputBox.classList.toggle("hidden");
}
