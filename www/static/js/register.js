/**
 * 
 * @param {SubmitEvent} formEvent
 */
function formChecker(formEvent) {
    formEvent.preventDefault();

    /** @type {HTMLFormElement} */
    const form = formEvent.target;

    // Create a FormData object to collect form data
    const formData = new FormData(form);

    // Convert FormData to JSON
    const jsonData = {};
    formData.forEach((value, key) => {
        jsonData[key] = value;
    });

    // Use the Fetch API to send a POST request with the JSON data
    fetch(form.action, {
        method: form.method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
    })
        .then(response => {
            if (!response.ok) {
                // If the response is not okay, reject the promise with the response text
                return response.text().then(text => {
                    throw new Error(text);
                });
            }
        })
        .then( ()=> {
            alert('Registration successful');
            // Redirect to home
            window.location.href = '/';
        })
        .catch(error => {
            // Handle errors
            console.error('Error:', error);
            alert(error.message);
        });
}
