function toggleDiv(divName) {
    const contentDiv = document.getElementById(divName);
    const overlayDiv = document.getElementById('overlay');

    if (contentDiv.style.display === 'none') {
        // Show the content div
        overlayDiv.style.display = 'block';
        contentDiv.style.display = 'block';
    } else {
        // Hide the content div
        contentDiv.style.display = 'none';
        overlayDiv.style.display = 'none';
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