
document.addEventListener('DOMContentLoaded', () => {
    /** @type {HTMLCollectionOf<HTMLTextAreaElement>} */
    const elements = document.getElementsByClassName('post-data-message')

    for (const element of elements) {
        element.style.height = element.scrollHeight + 'px';
    }
});