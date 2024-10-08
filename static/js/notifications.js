async function fetchNotifications() {
    try {
        const response = await fetch("/notifications");
        if (!response.ok) {
            handleErrorResponse(response);
            return; // Ensure we don't proceed if the response is not ok
        }
        const notifications = await response.json();

        // Log the notifications to debug
        console.log("Fetched notifications:", notifications);

        const notificationList = document.querySelector(".notification-list");
        notificationList.innerHTML = ""; // Clear existing notifications

        if (!notifications || !Array.isArray(notifications) || notifications.length === 0) {
            const noNotificationsParagraph = document.createElement("p");
            noNotificationsParagraph.textContent = "No notifications";
            notificationList.appendChild(noNotificationsParagraph);
            toggleVisibility("notifications");
            return;
        }

        notifications.forEach((notification) => {
            let message = "";

            if (notification.is_react) {
                message = `${notification.ReactionNotifi.username} ${notification.ReactionNotifi.reaction}d to your post: ${notification.ReactionNotifi.PostResponse.title}`;
            } else if (notification.is_comment) {
                message = `${notification.CommentNotifi.username} commented on your post: ${notification.CommentNotifi.PostResponse.title}`;
            } else if (notification.is_report) {
                message = `Your report on ${
                    notification.ReportRequestNotifi.reported_post_id
                } was ${notification.ReportRequestNotifi.accepted ? "accepted" : "rejected"}`;
            } else if (notification.is_promote_request) {
                message = `Your promote request was ${
                    notification.PromoteRequestNotification.accepted ? "accepted" : "rejected"
                }`;
            }

            const listItem = document.createElement("li");
            listItem.className = "notification-item";
            listItem.textContent = message;
            if (notification.is_react) {
                listItem.setAttribute(
                    "onclick",
                    `PostPageHandler(${JSON.stringify(notification.ReactionNotifi.PostResponse)})`
                );
            } else if (notification.is_comment) {
                listItem.setAttribute(
                    "onclick",
                    `PostPageHandler(${JSON.stringify(notification.CommentNotifi.PostResponse)})`
                );
            }
            notificationList.appendChild(listItem);
        });
        toggleVisibility("notifications");
    } catch (error) {
        console.error("Failed to fetch notifications:", error);
    }
}