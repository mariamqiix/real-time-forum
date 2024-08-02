const socket = new WebSocket(`ws://localhost:8080/ws`);

socket.onopen = function(event) {
    console.log("WebSocket is open now.");
};

socket.onmessage = function(event) {
    console.log("WebSocket message received:", event.data);
};

socket.onclose = function(event) {
    console.log("WebSocket is closed now.");
};

socket.onerror = function(error) {
    console.error("WebSocket error observed:", error);
};

document.getElementById("submitButton").addEventListener("click", function() {
    const recipientID = document.getElementById("recipientID").value;
    const message = document.getElementById("message").value;
    sendMessage(recipientID, message);
});

function sendMessage(recipientID, message) {
    const messageObject = {
        recipientID: recipientID,
        message: message,
    };
    socket.send(JSON.stringify(messageObject));
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