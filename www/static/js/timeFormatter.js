document.addEventListener('DOMContentLoaded', () => {
    const timeElements = document.getElementsByTagName('time');
    for (const timeElement of timeElements) {
        const date = new Date(timeElement.dateTime);
        // maybe format the date to something more readable later
        timeElement.textContent = timeElement.title = date.toLocaleString();
    }
});