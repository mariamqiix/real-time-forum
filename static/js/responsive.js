var div = document.getElementById("messagesBar");
var mediaQuery = window.matchMedia("(min-width: 1200px)");

function checkWidth(mediaQuery) {
    if (mediaQuery.matches) {
        div.style.display = "block";
    } else {
        div.style.display = "none";
    }
}

mediaQuery.addListener(checkWidth); // Add listener for changes in screen width

// Initial check when the page loads
checkWidth(mediaQuery);

var mediaQuery2 = window.matchMedia("(min-width: 1800px)");

//naviagtionBar resize
function checkSize(mediaQuery2) {
    const navigationBar = document.querySelector(".navigationBar");
    const navigationBarBtns = document.querySelectorAll(".navigationBarBtns");
    const logoName = document.querySelectorAll(".logoName")[0];
    const messagesBar = document.getElementById("messagesBar");
    const navigationContent = document.querySelectorAll(".navigationContent");
    if (!mediaQuery2.matches) {
        navigationBarBtns.forEach((btn) => {
            btn.style.display = "none";
        });
        navigationContent.forEach((content) => {
            messagesBar.style.display = "none";
            content.style.width = "100%";
        });

        logoName.style.display = "none";
        messagesBar.style.display = "none";
        // navigationBar.classList.add("center-content");
    } else {
        navigationBarBtns.forEach((btn) => {
            btn.style.display = "block";
            GetUserLoggedIn();
        });
        navigationContent.forEach((content) => {
            content.style.width = "70%";
            messagesBar.style.display = "block";
        });

        logoName.style.display = "block";
        GetUserLoggedIn();

        // navigationBar.classList.remove("center-content");
    }
}

mediaQuery2.addListener(checkSize); // Add listener for changes in screen width

// Initial check when the page loads
checkSize(mediaQuery2);