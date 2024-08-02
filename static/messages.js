function sendMessage(username){
    const usersDiv = document.getElementById(messagesBoxDiv);
    const msgDiv = document.getElementById(msgDiv);

    if (usersDiv.style.display === 'none') {
        // Show the content div
        usersDiv.style.display = 'block';
        msgDiv.style.display = 'none'
    } else {
        // Hide the content div
        usersDiv.style.display = 'none';
        msgDiv.style.display = 'block'
    }
}