const socket = new WebSocket(`ws://localhost:8080/ws`);

socket.onopen = function(event) {
    console.log("WebSocket is open now.");
};

socket.onmessage = function(event) {
    console.log("WebSocket message received:", event.data);
    // Assuming the received data is a JSON string representing a structs.Message
    const messageData = JSON.parse(event.data);
    const messageBox = document.getElementById(messageData.SenderId);

    // Create the new chat div
    const chatDiv = document.createElement("div");
    chatDiv.className = "fullMessage";

    // Create the message div
    const msgDiv = document.createElement("div");
    msgDiv.className = "msg";
    msgDiv.textContent = messageData.Message;
    msgDiv.style.backgroundColor = "white";

    // Create the message time div
    const msgTimeDiv = document.createElement("div");
    msgTimeDiv.className = "msgTime";
    msgTimeDiv.textContent = new Date(messageData.Time).toLocaleTimeString();

    // Append the message and time divs to the chat div
    chatDiv.appendChild(msgDiv);
    chatDiv.appendChild(msgTimeDiv);

    const chats = document.getElementById("UserChat");
    chats.appendChild(chatDiv);
    // Make the newMessageIcon green
    if (messageBox) {
        const newMessageIcon = messageBox.querySelector(".newMessageIcon");
        if (newMessageIcon) {
            newMessageIcon.style.backgroundColor = "lightgreen";
        }
    }
};

socket.onclose = function(event) {
    console.log("WebSocket is closed now.");
};

socket.onerror = function(error) {
    console.error("WebSocket error observed:", error);
};

// submitButton for the messages
// document.getElementById("submitButton").addEventListener("click", function() {
//     const recipientID = document.getElementById("recipientID").value;
//     const message = document.getElementById("message").value;
//     sendMessage(recipientID, message);
// });
function SendMessage(RecipientId) {
    const Message = document.getElementById("msgType").value;
    document.getElementById("msgType").value = "";
    fetch(`http://localhost:8080/user`, {
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
                // Clean the text by removing quotation marks and trimming whitespace
                const cleanedText = text.replace(/['"]/g, "").trim();
                const newcleanedText = cleanedText.replace("null", "");

                const messageObject = {
                    SenderId: newcleanedText,
                    RecipientId: RecipientId,
                    Message: Message,
                    Time: new Date().toISOString(),
                };
                console.log("Sending message:", messageObject);
                socket.send(JSON.stringify(messageObject));

                // Create the new chat div
                const chatDiv = document.createElement("div");
                chatDiv.className = "fullMessage";

                // Create the message div
                const msgDiv = document.createElement("div");
                msgDiv.className = "msg";
                msgDiv.textContent = Message;

                // Create the message time div
                const msgTimeDiv = document.createElement("div");
                msgTimeDiv.className = "msgTime";
                msgTimeDiv.textContent = new Date(messageObject.Time).toLocaleTimeString();

                // Append the message and time divs to the chat div
                chatDiv.appendChild(msgTimeDiv);

                chatDiv.appendChild(msgDiv);

                // Append the chat div to the msgUser chat container

                const chats = document.getElementById("UserChat");
                chats.appendChild(chatDiv);
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}

function checkUserOnline(userID) {
    fetch(`http://localhost:8080/checkUserOnline?userID=${userID}`, {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                console.log(`User ${userID} is online`);
            } else {
                console.log(`User ${userID} is offline`);
            }
        })
        .catch((error) => {
            console.error("Error checking user status:", error);
        });
}