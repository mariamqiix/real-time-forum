function changeContent(column, element) {
    // Remove the 'selected' class from all the columns
    var columns = document.getElementsByTagName("th");
    for (var i = 0; i < columns.length; i++) {
        columns[i].classList.remove("selected");
    }
    // Add the 'selected' class to the clicked column
    element.classList.add("selected");
}

function confirmAction() {
    if (confirm("Are you sure you want to proceed?")) {
        // User clicked "OK"
        // Perform the desired action here
        // For example, submit a form or execute a function
        console.log("Proceeding with the action...");
    } else {
        // User clicked "Cancel" or closed the dialog
        // Cancel the action or do nothing
        console.log("Action canceled.");
    }
}


function sendPromotionRequest() {
    var answer = document.getElementById("answer").value;
    // Add code here to handle sending the answer
    if (answer.trim() === "") {
        alert("Please provide an answer before sending.");
        return;
    }

    const formData = new FormData();
    formData.append("answer", answer);

    fetch("http://localhost:8080/PromotionRequest", {
            method: "POST",
            body: formData,
        })
        .then((response) => {
            if (response.ok) {
                alert("successful");
            } else {
                alert("failed");
            }
        })
        .catch((error) => {
            console.error("Error:", error);
        });
    toggleDiv("request-moderator");
    document.getElementById("answer").value = "";
}
// Clear any existing posts
searchoutput.innerHTML = "";


function EditPost() {}

// // Fetch the JSON file
// fetch('./countries.json')
//   .then(response => response.json())
//   .then(data => {
//     const countrySelect = document.getElementById('countrySelect');

//     // Populate the select element with options
//     data.countries.forEach(country => {
//       const option = document.createElement('option');
//       option.text = country;
//       option.value = country;
//       countrySelect.add(option);
//     });
//   })
//   .catch(error => console.error('Error loading countries:', error));


// function loadCountries() {
//     fetch('countries.json')
//         .then(response => {
//             if (!response.ok) {
//                 throw new Error('Failed to load the JSON file.');
//             }
//             return response.json();
//         })
//         .then(data => {
//             const countries = data.countries;
//             console.log(countries); // You can process the countries array as needed
//         })
//         .catch(error => {
//             console.error('Error loading countries:', error);
//         });
// }

// // Call the function to load the countries
// loadCountries();

