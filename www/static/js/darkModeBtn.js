/* for dark mode js */

// Get the checkbox element
const checkbox = document.querySelector('.darkmodebtn input[type="checkbox"]');

// Add event listener for the checkbox change event
checkbox.addEventListener('change', function () {
    // Check if the checkbox is checked
    if (this.checked) {
        // Remove the dark mode class from the body element
        document.body.classList.remove('dark-mode');
    } else {
        // Add the dark mode class to the body element
        document.body.classList.add('dark-mode');
    }
});

const checkbox2 = document.querySelector('.darkmodebtn2 input[type="checkbox"]');

// Add event listener for the checkbox change event
checkbox2.addEventListener('change', function () {
    // Check if the checkbox is checked
    if (this.checked) {
        // Remove the dark mode class from the body element
        document.body.classList.remove('dark-mode');
    } else {
        // Add the dark mode class to the body element
        document.body.classList.add('dark-mode');
    }
});