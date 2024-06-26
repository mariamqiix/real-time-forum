var div = document.getElementById('messagesBar');
var mediaQuery = window.matchMedia('(min-width: 1200px)');

function checkWidth(mediaQuery) {
  if (mediaQuery.matches) {
    div.style.display = 'block';
  } else {
    div.style.display = 'none';
  }
}

mediaQuery.addListener(checkWidth); // Add listener for changes in screen width

// Initial check when the page loads
checkWidth(mediaQuery);


function toggleDiv(divName) {
  const contentDiv = document.getElementById(divName);
  const overlayDiv = document.getElementById('overlay');

  if (contentDiv.style.display === 'none') {
    // Show the content div
    contentDiv.style.display = 'block';
    overlayDiv.style.display = 'block';
  } else {
    // Hide the content div
    contentDiv.style.display = 'none';
    overlayDiv.style.display = 'none';
  }
}

async function submitSelection(divName) {
    // Get the selected categories
    const categoryCheckboxes = document.querySelectorAll(
        '.categoriesList input[type="checkbox"]:checked'
    );

    const selectedCategories = Array.from(categoryCheckboxes).map((checkbox) => checkbox.value);

    // Get the title and topic
    const title = document.getElementById("newPostTitle").value;
    const topic = document.getElementById("topic").value;
    console.log("Selected option:", selectedCategories[0]);
    // Perform further actions with the selected option

    // Create the request body
    const formData = new FormData();
    formData.append("title", title);
    formData.append("topic", topic);
    formData.append("selectedCategories", selectedCategories);

        try {
        // Send the POST request
        const response = await fetch("http://localhost:8080/post/add", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            console.log("post added successful");
            toggleDiv(divName); // Hide the content div after submission

        } else {
            const error = await response.text();
            console.error("Signup failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }

}

var mediaQuery2 = window.matchMedia('(min-width: 1800px)');

//naviagtionBar resize
function checkSize(mediaQuery2) {
  const navigationBar = document.querySelector('.navigationBar');
  const navigationBarBtns = document.querySelectorAll('.navigationBarBtns');
  const logoName = document.querySelectorAll('.logoName')[0];

  if (!mediaQuery2.matches) {
    navigationBarBtns.forEach(btn => {
      btn.style.display = 'none';
    });

    logoName.style.display = 'none';
    // navigationBar.classList.add("center-content");
  } else {
    navigationBarBtns.forEach(btn => {
      btn.style.display = 'block';
    });

    logoName.style.display = 'block';
    // navigationBar.classList.remove("center-content");
  }
}

mediaQuery2.addListener(checkSize); // Add listener for changes in screen width

// Initial check when the page loads
checkSize(mediaQuery2);


var textField = document.getElementById('newPostTitle');
var submitButton = document.getElementById('submitBtn');

textField.addEventListener('input', function () {
  if (textField.value.trim() !== "") {
    submitButton.style.cursor = 'pointer';
  } else {
    submitButton.style.cursor = 'not-allowed';
  }
});

textField.addEventListener('keyup', function () {
  if (textField.value.trim() !== "") {
    submitButton.style.cursor = 'pointer';
    submitButton.style.opacity = '1';
  } else {
    submitButton.style.cursor = 'not-allowed';
    submitButton.style.opacity = '0.5';

  }
});

function showHome() {
  var homeDiv = document.getElementById("home");

  if (homeDiv.style.display === "none") {
    homeDiv.style.display = "block";
  } else {
    homeDiv.style.display = "none";
  }
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

function validateField(input) {
  const errorId = input.id + "Error";
  const errorElement = document.getElementById(errorId);
  const value = input.value.trim();

  if (value === "") {
    errorElement.textContent = "Field is required";
  } else {
    errorElement.textContent = "";
  }
}


async function submitForm() {
      const inputs = document.querySelectorAll("#registration-form input");
      for (let i = 0; i < inputs.length; i++) {
          validateField(inputs[i]);
      }

    // Get the form values
    const firstName = document.getElementById("firstName").value;
    const lastName = document.getElementById("lastName").value;
    const email = document.getElementById("email").value;
    const username = document.getElementById("username").value;
    const dob = document.getElementById("dob").value;
    const password = document.getElementById("password").value;

    // Create the request body
    const formData = new FormData();
    formData.append("firstName", firstName);
    formData.append("lastName", lastName);
    formData.append("email", email);
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
                toggleVisibility("home");

        } else {
            const error = await response.text();
            console.error("Signup failed:", error);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// display posts in home

function HomePageRequest() {
    // Send the form data to the Go server
    fetch("http://localhost:8080/", {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error("Error: " + response.status);
            }
            return response.json();
        })
        .then((data) => {
            displayPost(data); // Pass the response data to the displayScores function
        })
        .catch((error) => {
            console.error("Error:", error);
        });
}



function displayPost(data) {
    console.log("Received data:", data); // Log the received data to check its structure

    if (Array.isArray(data.Posts)) {
        const homeNavigationContent = document.getElementById("home navigationContent");

        data.Posts.forEach((post) => {
            const postBox = document.createElement("div");
            postBox.classList.add("postBox");

            const postUserPic = document.createElement("div");
            postUserPic.classList.add("postUserPic");
            postBox.appendChild(postUserPic);

            const postTitle = document.createElement("div");
            postTitle.classList.add("postTitle");

            const titleElement = document.createElement("span");
            titleElement.classList.add("title");
            titleElement.setAttribute(
                "style",
                '--font-selector:R0Y7S3VsaW0gUGFyay1yZWd1bGFy;--framer-font-family:"Kulim Park";--framer-font-size:40px'
            );
            titleElement.textContent = post.title;
            postTitle.appendChild(titleElement);

            const postUserName = document.createElement("span");
            postUserName.classList.add("postUserName");

            postUserName.textContent = post.author.username;
            postTitle.appendChild(postUserName);

            const postContent = document.createElement("span");
            postContent.classList.add("postContent");
            postContent.setAttribute(
                "style",
                '--font-selector:R0Y7S3VmYW0tcmVndWxhcg==;--framer-font-family:"Kufam";--framer-font-size:15px;--framer-text-color:rgba(91, 91, 91, 1)'
            );
            postContent.textContent = post.message;
            postTitle.appendChild(postContent);
            postBox.appendChild(postTitle);
            homeNavigationContent.appendChild(postBox);
        });
    } else {
        console.error("Invalid data format. Expected an array of posts.");
    }
}


function fetchAndAppendCategories() {
    // Fetch the categories from the server
      const categoriesContainer = document.querySelector(".categoriesList");
      categoriesContainer.innerHTML = "";

    fetch("http://localhost:8080/category")
        .then((response) => response.json())
        .then((data) => {
            // Iterate over the categories and append them to the container
            data.forEach((category) => {
                appendCategory(category.Name,category.Id);
            });
        })
        .catch((error) => {
            console.error("Error fetching categories:", error);
        });
}

function appendCategory(name,Id) {
    // Get the container element where the categories will be added
    const categoriesContainer = document.querySelector(".categoriesList");

    // Create a new label element
    const label = document.createElement("label");

    // Create a new checkbox input element
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.name = name;
    checkbox.value = Id ;
    checkbox.classList.add("checkBoxCateg");

    // Create the category name text node
    const categoryText = document.createTextNode(name);

    // Append the checkbox and category name to the label
    label.appendChild(checkbox);
    label.appendChild(categoryText);

    // Append the label to the container
    categoriesContainer.appendChild(label);
    categoriesContainer.appendChild(document.createElement("br"));
}

fetchAndAppendCategories();


HomePageRequest();
