function _reloadTheme() {
    if (localStorage.getItem('theme') === 'dark') {
        document.documentElement.style.setProperty('--bg-color', '#141318');
        document.documentElement.style.setProperty('--border-color', 'white');
    } else {
        // MARK: - PRIMARY COLORS (BG, Buttons, Links, etc.)
        document.documentElement.style.setProperty('--primary-color', '#FFB428');
        document.documentElement.style.setProperty('--primary-color-hover', 'orange');
        document.documentElement.style.setProperty('--secondary-color', '#ededed');
        document.documentElement.style.setProperty('--secondary-color-hover', '#e0e0e0');
        // MARK: - TEXT COLORS
        document.documentElement.style.setProperty('--text-color', 'black');
        document.documentElement.style.setProperty('--secondary-text-color', '#808080');
        document.documentElement.style.setProperty('--secondary-text-color-hover', 'var(--primary-color)');
        document.documentElement.style.setProperty('--disabled-color', '#B4B4B4');
        document.documentElement.style.setProperty('--bg-color', 'white');
        document.documentElement.style.setProperty('--border-color', 'black');
    }
}

function reloadTheme() {
    if (localStorage.getItem('theme') === 'dark') {
        localStorage.setItem('theme', 'light')
    } else {
        localStorage.setItem('theme', 'dark')
    }
    _reloadTheme()
}

/**
 * 
 * @param {HTMLButtonElement} element the button which was clicked
 */
function showDropdown(element) {
    // get the ul item under the button
    /** @type {HTMLUListElement} */
    const dropdownElement = element.nextElementSibling;
    // check if the ul is hidden
    if (dropdownElement.style.display === 'none' || dropdownElement.style.display === '') {
        // show the ul
        dropdownElement.style.display = 'block';
        document.body.addEventListener('click', (event) => {
            if (event.target !== element) {
                dropdownElement.style.display = 'none';
            }
        });
    } else {
        // hide the ul
        dropdownElement.style.display = 'none';
    }
}

/**
 * 
 * @param {HTMLButtonElement} button the dropdown button that was clicked 
 */
function dropdownSort(button) {
    window.location.href = `/?sort=${button.dataset.type}`
}

document.addEventListener('DOMContentLoaded', () => {
    _reloadTheme()
})

/**
 * 
 * @param {HTMLButtonElement} element 
 */
function handleReaction(element) {
    const reactionType = element.dataset.type;

    // get the post id from teh containing article element
    const postContainer = element.closest('[data-post]');
    const postId = postContainer.dataset.post;

    console.info(reactionType, postId);

    const url = `/post/${postId}/${reactionType}`;
    let method = 'POST';
    if (element.classList.contains('reaction-button-reacted')) {
        method = 'DELETE';
    }

    fetch(url, {
        method: method,
    }).then(response => {
        if (response.ok) {
            // reload the page
            let url = new URL(window.location.href);
            url.hash = postId;
            window.location.href = url.href;
            window.location.reload();
        } else {
            throw new Error('Network response was not ok.');
        }
    }).catch(error => {
        console.error('Error:', error);
        alert('An error occurred, please try again later.');
    });

}