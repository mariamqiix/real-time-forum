function validateCreatePost(input) {
    var title = document.getElementById("newPostTitle").value;
    var topic = document.getElementById("topic").value;
    var submitButton = document.getElementById("submitBtn");
    var EditButton = document.getElementById("EditBtn");

    const categoryCheckboxes = document.querySelectorAll(
        '.categoriesList input[type="checkbox"]:checked'
    );

    if (title.trim() !== "" && topic.trim() !== "" && categoryCheckboxes.length > 0) {
        submitButton.style.cursor = "pointer";
        submitButton.style.opacity = "1";
        EditButton.style.cursor = "pointer";
        EditButton.style.opacity = "1";
    } else {
        submitButton.style.cursor = "not-allowed";
        submitButton.style.opacity = "0.5";
        EditButton.style.cursor = "not-allowed";
        EditButton.style.opacity = "0.5";
    }
}

function validateLogInField(input) {
    const fieldName = input.id;
    const errorId = fieldName + "Error";
    const errorElement = document.getElementById(errorId);
    const value = input.value.trim();

    // Validate username or email
    if (fieldName === "usernameLogin") {
        if (value === "") {
            errorElement.textContent = "Username or email is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate password
    if (fieldName === "passwordLogin") {
        if (value === "") {
            errorElement.textContent = "Password is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    checkLogFormValidity();
}

function validateField(input, type = 0) {
    const fieldName = input.id;
    const errorId = fieldName + "Error";
    const errorElement = document.getElementById(errorId);
    const value = input.value.trim();

    // Validate first name
    if (fieldName === "firstName") {
        if (value === "" || value.length < 4 || value.length > 20) {
            errorElement.textContent = "First name is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate last name
    if (fieldName === "lastName") {
        if (value === "" || value.length < 4 || value.length > 20) {
            errorElement.textContent = "Last name is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate email
    if (fieldName === "email") {
        if (value === "") {
            errorElement.textContent = "Email is required.";
        } else if (!isValidEmail(value)) {
            errorElement.textContent = "Invalid email format.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate username
    if (fieldName === "username") {
        if (value === "" || value.length < 4 || value.length > 20) {
            errorElement.textContent = "Username is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate date of birth
    if (fieldName === "dob") {
        if (value === "") {
            errorElement.textContent = "Date of birth is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate country
    if (fieldName === "country") {
        if (value === "" || value.length < 4 || value.length > 20) {
            errorElement.textContent = "Country is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate password
    if (fieldName === "password") {
        if (value === "") {
            errorElement.textContent = "Password is required.";
            console.log("Password is required.");
        } else if (value.length < 8) {
            errorElement.textContent = "Password must be at least 8 characters long.";
            console.log("Password must be at least 8 characters long.");
        } else {
            errorElement.textContent = "";
            const password2 = document.getElementById("confirmPassword").value.trim();
            if (password2 !== "" && value != password2) {
                document.getElementById("confirmPasswordError").textContent =
                    "Passwords do not match.";
            }
        }
    }

    // Validate confirm password
    if (fieldName === "confirmPassword") {
        const password = document.getElementById("password").value.trim();
        if (value === "") {
            errorElement.textContent = "Please re-enter the password.";
        } else if (value !== password) {
            errorElement.textContent = "Passwords do not match.";
        } else {
            errorElement.textContent = "";
        }
    }
    checkFormValidity();
}

function isValidEmail(email) {
    // Simple email validation using regular expression
    const emailRegex = /^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/;
    return emailRegex.test(email);
}

function checkFormValidity() {
    let allFieldsFilled = true;
    let allErrorsEmpty = true;

    // Add event listeners to the input fields
    const inputFields = document.querySelectorAll("#formReg input");
    inputFields.forEach(function(input) {
        input.addEventListener("input", checkFormValidity);
    });

    // Check if all input fields are filled
    inputFields.forEach(function(input) {
        if (input.value.trim() === "") {
            allFieldsFilled = false;
        }
    });

    // Check if any error messages are displayed
    const errorMessages = document.querySelectorAll(".error");
    errorMessages.forEach(function(error) {
        if (error.textContent.trim() !== "") {
            console.log(error.textContent);

            allErrorsEmpty = false;
        }
    });

    var registerBtn = document.getElementById("registerBtn");
    var isValid = allFieldsFilled && allErrorsEmpty;
    if (isValid) {
        registerBtn.removeAttribute("disabled");
        registerBtn.classList.remove("disabled");
    } else {
        registerBtn.setAttribute("disabled", "disabled");
        registerBtn.classList.add("disabled");
    }
}

function checkLogFormValidity() {
    let allFieldsFilled = true;
    let allErrorsEmpty = true;

    // Add event listeners to the input fields
    const inputFields = document.querySelectorAll("#formLog input");
    inputFields.forEach(function(input) {
        input.addEventListener("input", checkLogFormValidity);
    });

    // Check if all input fields are filled
    inputFields.forEach(function(input) {
        if (input.value.trim() === "") {
            allFieldsFilled = false;
        }
    });

    // Check if any error messages are displayed
    const errorMessages = document.querySelectorAll(".errorLogin");
    errorMessages.forEach(function(error) {
        if (error.textContent.trim() !== "") {
            allErrorsEmpty = false;
        }
    });

    var loginBtn = document.getElementById("loginBtn");
    var isValid = allFieldsFilled && allErrorsEmpty;

    if (isValid) {
        loginBtn.removeAttribute("disabled");
        loginBtn.classList.remove("disabled");
    } else {
        loginBtn.setAttribute("disabled", "disabled");
        loginBtn.classList.add("disabled");
    }
}

function validateUserInfoForm() {
    const username = document.getElementById("ChangeUserName").value.trim();
    const firstName = document.getElementById("ChangeFirstName").value.trim();
    const lastName = document.getElementById("ChangeLastName").value.trim();
    const email = document.getElementById("ChangeEmail").value.trim();
    const dateOfBirth = document.getElementById("ChangeDOB").value.trim();

    if (!username) {
        alert("Username is required.");
        return false;
    } else if (username.length < 4 || username.length > 20) {
        alert("Username must be between 4 and 20 characters.");
        return false;
    }
    if (!firstName) {
        alert("First Name is required.");
        return false;
    } else if (firstName.length < 3 || firstName.length > 20) {
        alert("First Name must be between 3 and 20 characters.");
        return false;
    }
    if (!lastName) {
        alert("Last Name is required.");
        return false;
    } else if (lastName.length < 3 || lastName.length > 20) {
        alert("Last Name must be between 3 and 20 characters.");
        return false;
    }
    if (!email) {
        alert("Email is required.");
        return false;
    } else if (!isValidEmail(email)) {
        alert("Invalid email format.");
        return false;
    }
    if (!dateOfBirth) {
        alert("Date of Birth is required.");
        return false;
    } else if (dateOfBirth > new Date()) {
        alert("Invalid date of birth.");
        return false;
    }
    // Additional validation can be added here (e.g., email format, date format)
    return true;
}