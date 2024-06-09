/**
 * 
 * @param {HTMLButtonElement} element 
 */
function selectedCategory(element) {
    element.classList.toggle('selected-category');
}

/**
 * 
 * @param {SubmitEvent} formEvent
 */
function processPostForm(formEvent) {
    formEvent.preventDefault()

    /** @type {HTMLFormElement} */
    const form = formEvent.target;

    // Create a FormData object to collect form data
    const formData = new FormData(form);

    // Convert FormData to JSON
    const jsonData = {};
    formData.forEach((value, key) => {
        jsonData[key] = value;
    });

    // Check if the file is over 20 mbs
    if (jsonData.hasOwnProperty('image')) {
        const file = jsonData['image'];
        if (file.size > 20971520) {
            alert('File size is too large, max size is 20MB');
            return;
        }
    }


    // set the categories

    // check if the url in this format '/post/{post_id}/comment' with regex
    const url = window.location.pathname;
    const urlRegex = /^\/post\/\d+\/comment$/;
    const editRegex = /^\/post\/\d+\/edit$/;
    if (!urlRegex.test(url) && !editRegex.test(url)) {
        const activeCategoriesButtons = document.querySelectorAll('.selected-category');
        if (activeCategoriesButtons.length === 0) {
            alert('Please select at least one category');
            return;
        }
        activeCategoriesButtons.forEach(element => {
            formData.append('categories', element.textContent)
        });
    }


    fetch(form.action, {
        method: form.method,
        body: formData,
    })
        .then(response => {
            if (!response.redirected) {
                // If the response is not okay, reject the promise with the response text
                return response.text().then(text => {
                    throw new Error(text);
                });
            }
            // follow the redirect
            window.location.href = response.url;
        })
        .catch(error => {
            // Handle errors
            console.error('Error:', error);
            alert(error.message);
        });
}