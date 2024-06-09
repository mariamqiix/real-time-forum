document.addEventListener("DOMContentLoaded", () => {
    getNotifications();
});

function getNotifications() {
    console.log('notification triggered!');
    fetch('/notifications')
        .then(response => response.json())
        .then(data => {
            const notificationsSideMenu = document.getElementById('mySidepanel');
            notificationsSideMenu.innerHTML = `
                <a href="javascript:void(0)" class="closebtn" onclick="closeNav()">×</a>
            `; // Clear previous notifications and add close button

            data.forEach(notification => {
                const notificationElement = document.createElement('div');
                notificationElement.classList.add('notification');
                if (notification.read) {
                    notificationElement.classList.add('read'); // Add a class for read notifications
                }
                notificationElement.innerHTML = `
                <div class="noti">
                  <h3 style="color: white;">${notification.title}</h3>
                  <p style="color: white;">${notification.body}</p>
                  <a href="${notification.link}" data-id="${notification.id}" class="view-notification">View</a>
                </div>
                `;
                notificationsSideMenu.appendChild(notificationElement);
            });

            // Add event listeners to the "View" links
            document.querySelectorAll('.view-notification').forEach(link => {
                link.addEventListener('click', markNotificationRead);
            });
        })
        .catch(error => {
            console.error('Error retrieving notifications:', error);
        });
}

function markNotificationRead(event) {
    event.preventDefault();
    const notificationId = this.getAttribute('data-id');
    fetch(`/notifications/${notificationId}/read`, {
        method: 'POST',
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            console.log('Notification marked as read:', data);
            // Redirect to the notification link
            window.location.href = this.getAttribute('href');
        })
        .catch(error => {
            console.error('Error marking notification as read:', error);
        });
}

function closeNav() {
    document.getElementById("mySidepanel").style.width = "0";
}
