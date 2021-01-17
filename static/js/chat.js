let typing = document.getElementById("chat-form")
typing.addEventListener("keypress", (event) => {
    event.preventDefault();
    if (event.key !== "Enter") {typing.value += event.key;return;}
	var message_element = document.getElementsByTagName("input")[0];
	var message = message_element.value;

	if (message.toString().length) {
        MessageAdd('<div class="message">' + message + '</div>')
		message_element.value = "";
	}
}, false);

function MessageAdd(message) {
	var chat_messages = document.getElementById("chat-messages");

	chat_messages.insertAdjacentHTML("beforeend", message);
	chat_messages.scrollTop = chat_messages.scrollHeight;
}

let chatForm = document.getElementById("chat-form")
let userChat = document.getElementById("user-chat")