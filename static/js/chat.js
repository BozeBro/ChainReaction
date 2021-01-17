
function MessageAdd(message, username, color) {
	var chat_messages = document.getElementById("chat-messages");
    msg = `<div><span style="color:${color}">${username}</span> ${message}</div>`
	chat_messages.insertAdjacentHTML("beforeend", msg);
	chat_messages.scrollTop = chat_messages.scrollHeight;
}
