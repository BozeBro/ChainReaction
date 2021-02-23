
function MessageAdd(message, username, color) {
	var chat_messages = document.getElementById("chat-messages");
	var date = new Date();
	var data = changeHours(Number(date.getHours()));
	var hours = data[0];
	var meridian = data[1];
    msg = `<div><span style="color:purple">${"| "+hours+":"+date.getMinutes()+" "+meridian+" | "}</span><span style="color:${escapeHtml(color)}">${escapeHtml(username)}</span> ${escapeHtml(message)}</div>`
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

  function changeHours(hours) {
	  if (hours > 12) {
		  return [(hours-12).toString(), "PM"]
		}
		return [hours.toString(), "AM"]
  }
var beginning;
var curTime;
function timing(timestamp) {
	curTime = timestamp
    requestAnimationFrame((timestamp) => timing(timestamp))
}