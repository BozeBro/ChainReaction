
function MessageAdd(message, username, color) {
	var chat_messages = document.getElementById("chat-messages");
    msg = `<div><span style="color:${escapeHtml(color)}">${escapeHtml(username)}</span> ${escapeHtml(message)}</div>`
	chat_messages.insertAdjacentHTML("beforeend", msg);
	chat_messages.scrollTop = chat_messages.scrollHeight;
}

function escapeHtml(text) {
	var map = {
	  '&': '&amp;',
	  '<': '&lt;',
	  '>': '&gt;',
	  '"': '&quot;',
	  "'": '&#039;'
	};

	return text.replace(/[&<>"']/g, (m) => map[m]);
  }