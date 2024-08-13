async function submitLoginForm() {
    const usernameInput = document.getElementById("usernameLogin");
    const passwordInput = document.getElementById("passwordLogin");

    const username = usernameInput.value;
    const password = passwordInput.value;

    // Create the request body
    const formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);

    try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Login successful");
            const loginSpan = document.getElementById("loginSpan");
            loginSpan.innerHTML = "LOGOUT";
            GetUserLoggedIn();
            HomePageRequest();
            initializeWebSocket();
        } else {
            const error = await response.text();
            console.log(error);
            alert("Login failed: " + error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

async function submitForm() {
    // Get the form values
    const firstName = document.getElementById("firstName").value;
    const lastName = document.getElementById("lastName").value;
    const email = document.getElementById("email").value;
    const username = document.getElementById("username").value;
    const dob = document.getElementById("dob").value;
    const country = document.getElementById("country").value;
    const password = document.getElementById("password").value;
    const gender = document.getElementById("gender").value;
    // Create the request body
    const formData = new FormData();
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
    formData.append("country", country);
    formData.append("gender", gender);
    formData.append("username", username);
    formData.append("dob", dob);
    formData.append("password", password);

    try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/signup", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("Signup successful");
            GetUserLoggedIn();
            const loginSpan = document.getElementById("loginSpan");
            loginSpan.innerHTML = "Logout";
            HomePageRequest(); // Get the form values
            const firstName = document.getElementById("firstName");
            firstName.innerHTML = "";
            const lastName = document.getElementById("lastName");
            lastName.innerHTML = "";
            const email = document.getElementById("email");
            email.innerHTML = "";
            const username = document.getElementById("username");
            username.innerHTML = "";
            const dob = document.getElementById("dob");
            dob.innerHTML = "";
            const password = document.getElementById("password");
            password.innerHTML = "";
            const country = document.getElementById("country");
            country.innerHTML = "";
            const gender = document.getElementById("gender");
            gender.innerHTML = "";
            initializeWebSocket();
        } else {
            const error = await response.text();
            alert("Signup failed: " + error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}
var mediaQuery2 = window.matchMedia("(min-width: 1800px)");

async function GetUserLoggedIn() {
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
            const setting = document.querySelector(".settingIcon.navigationBarIcons");
            const bellIcon = document.querySelector(".bellIcon.navigationBarIcons");
            const settingBtn = document.querySelector(".settingBtn.navigationBarBtns");
            const notificationsBtn = document.querySelector(".notificationsBtn.navigationBarBtns");
            const penIcon = document.querySelector(".penIcon.mainPageIcons");
            const profileIcon = document.querySelector(".profileIcon.navigationBarIcons");
            const profileBtn = document.querySelector(".profileBtn.navigationBarBtns");
            if (text === "null") {
                setting.style.display = "none";
                bellIcon.style.display = "none";
                settingBtn.style.display = "none";
                notificationsBtn.style.display = "none";
                penIcon.style.display = "none";
                profileIcon.style.display = "none";
                profileBtn.style.display = "none";
                const loginSpan = document.getElementById("loginSpan");
                loginSpan.innerHTML = "Sign up";
                const replayPostButton = document.getElementById("replayPost-button");
                replayPostButton.style.display = "none";
                return;
            } else {
                setting.style.display = "block";
                bellIcon.style.display = "block";
                penIcon.style.display = "block";
                profileIcon.style.display = "block";
                if (mediaQuery2.matches) {
                    notificationsBtn.style.display = "block";
                    settingBtn.style.display = "block";
                    profileBtn.style.display = "block";
                } else {
                    notificationsBtn.style.display = "none";
                    settingBtn.style.display = "none";
                    profileBtn.style.display = "none";
                }
                const loginSpan = document.getElementById("loginSpan");
                loginSpan.innerHTML = "Logout";
                const replayPostButton = document.getElementById("replayPost-button");
                replayPostButton.style.display = "block";
                return text;
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}
GetUserLoggedIn();

async function logout() {
    try {
        const response = await fetch("http://localhost:8080/logout", {
            method: "POST",
        });

        if (response.ok) {
            console.log("Logout successful");
            GetUserLoggedIn();
            initializeWebSocket();
            // Perform any additional actions needed after successful logout
            HomePageRequest();
        } else {
            const error = await response.text();
            console.error("Logout failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

setInterval(GetUserLoggedIn, 5000);