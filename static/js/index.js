"use strict"
let lay = document.querySelector(".layout");
let cre = document.querySelector(".create");
let join = document.querySelector(".join");
let spaces = document.querySelector(".spaces");
let errCre = document.getElementById("errCre");
let errJoin = document.getElementById("errJoin");
let popCre = document.getElementById("pop-create"); // The box that appears when Create is clicked
let popJoin = document.getElementById("pop-join");  //  The box that appears when Join is clicked
// popup tells what popup is active
let popup = null;
let creHandler = () => {
    // once the box is clicked, the only way out is to click exit
    if (popup) { return }
    popCre.style.display = "flex"; // Make the popup appear
    popup = popCre;
}
let joinHandler = () => {
    // once the box is clicked, the only way out is to click exit
    if (popup) { return }
    popJoin.style.display = "flex"; // Make the popup appear
    popup = popJoin;

}
let btnClicked = (e) => {
    const btnClicked = e.target.nodeName === "BUTTON";
    if (!btnClicked) return
    switch (e.target.className) {
        case "ext":
            popup.style.display = "none";
            popup = null;
            errCre.innerHTML = "";
            errJoin.innerHTML = "";
            break;
        case "ent":
            if (popup.id === "pop-join") {
                let room = document.getElementById("room").value;
                let pin = document.getElementById("pin").value;
                if (room === "" || pin === "") return
                fetch("/api/join/", {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ room: room, pin: pin, }),
                })
                    .then(async (res) => {
                        switch (res.status) {

                            case 406:
                                errJoin.innerHTML = "Pin is wrong"
                                break
                            case 403:
                                errJoin.innerHTML = "The room is Full"
                                break
                            case 404:
                                errJoin.innerHTML = "The room doesn't exist"
                                break
                            case 200:
                                location.href = '/game/' + await res.text()
                            default:
                                errJoin.innerHTML = "Waiting on the server"
                        }
                    })

            } else {
                // Creating a room
                let DOMName = document.getElementById("name");
                let room = DOMName.value;
                let players = document.getElementById("players").value;
                if (!room) {
                    errCre.innerHTML = "ENTER A ROOM NAME";
                    return
                }
                DOMName.value = ""
                fetch("/api/create/", {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ Players: players, Room: room }),
                    redirect: 'follow',
                    credentials: 'same-origin',
                })
                    .then(async (res) => {
                        switch (res.status) {
                            case 409:
                                errCre.innerHTML = "You didn't enter everything";
                                break
                            case 400:
                                errCre.innerHTML = "WHAT ARE YOU DOING?!";
                                break
                            case 200:
                                location.href = '/game/' + await res.text()
                                break
                            default:
                                errCre.innerHTML = "Waiting on the server"
                        }
                    })
            }
    }
}
window.onload = () => {
    spaces.addEventListener("click", btnClicked)
    cre.addEventListener("click", creHandler)
    join.addEventListener("click", joinHandler)
}

