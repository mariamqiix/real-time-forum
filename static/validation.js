var textField = document.getElementById('newPostTitle');
var submitButton = document.getElementById('submitBtn');

textField.addEventListener('input', function() {
    if (textField.value.trim() !== "") {
        submitButton.style.cursor = 'pointer';
    } else {
        submitButton.style.cursor = 'not-allowed';
    }
});

textField.addEventListener('keyup', function() {
    if (textField.value.trim() !== "") {
        submitButton.style.cursor = 'pointer';
        submitButton.style.opacity = '1';
    } else {
        submitButton.style.cursor = 'not-allowed';
        submitButton.style.opacity = '0.5';

    }
});


//Registration validation 

// // Update button state and styles
// function updateButtonState() {
//     var registerBtn = document.getElementById('registerBtn');
//     var isValid = validateForm();

//     if (isValid) {
//         registerBtn.removeAttribute('disabled');
//         registerBtn.classList.remove('disabled');
//     } else {
//         registerBtn.setAttribute('disabled', 'disabled');
//         registerBtn.classList.add('disabled');
//     }
// }

// // Validation function for the entire form
// function validateForm() {
//     var firstName = document.getElementById('firstName');
//     var lastName = document.getElementById('lastName');
//     var email = document.getElementById('email');
//     var username = document.getElementById('username');
//     var dob = document.getElementById('dob');
//     var country = document.getElementById('country');
//     var password = document.getElementById('password');
//     var confirmPassword = document.getElementById('confirmPassword');

//     return ((firstName.length > 4) &&
//         (lastName.length > 4) &&
//         (email.length > 4) &&
//         (username.length > 4) &&
//         validateDate(dob !== null) &&
//         (country.length > 4) &&
//         (password.length > 4) &&
//         validatePasswordMatch());
// }

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
        if (value === "") {
            errorElement.textContent = "First name is required.";
        } else {
            errorElement.textContent = "";
        }
    }

    // Validate last name
    if (fieldName === "lastName") {
        if (value === "") {
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
        if (value === "") {
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
        if (value === "") {
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
                document.getElementById("confirmPasswordError").textContent = "Passwords do not match.";
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
        if (error.textContent !== "") {
            allErrorsEmpty = false;
        }
    });

    var registerBtn = document.getElementById('registerBtn');
    var isValid = (allFieldsFilled && allErrorsEmpty);

    if (isValid) {
        registerBtn.removeAttribute('disabled');
        registerBtn.classList.remove('disabled');
    } else {
        registerBtn.setAttribute('disabled', 'disabled');
        registerBtn.classList.add('disabled');
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
    const errorMessages = document.querySelectorAll(".error");
    errorMessages.forEach(function(error) {
        if (error.textContent !== "") {
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