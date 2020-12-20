"use strict"
let lay = document.querySelector(".layout");
let cre = document.querySelector(".create");
let join = document.querySelector(".join");
let spaces = document.querySelector(".spaces");
let errCre = document.getElementById("errCre");
let errJoin = document.getElementById("errJoin");
let popCre = document.getElementById("pop-create");
let popJoin = document.getElementById("pop-join");
// popup tells what popup is active
let popup = null;

let creHandler = () => {
    if (popup) { return }
    popCre.style.display = "flex";
    popup = popCre;
}
let joinHandler = () => {
    if (popup) { return }
    popJoin.style.display = "flex";
    popup = popJoin;

}
let btnClicked = (e) => {
    const btnClicked = e.target.nodeName === "BUTTON";
    if (!btnClicked) { return }
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
                if (!(room && pin)) { return }
                fetch("http://" + document.location.host + "/api/join", {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ room: room, pin: pin, }),
                })
                    .then((res) => {
                        console.log(res)
                    })

            } else {
                // Creating a room
                let DOMName = document.getElementById("name");
                let name = DOMName.value;
                let players = document.getElementById("players").value;
                if (!name) errCre.innerHTML = "ENTER A ROOM NAME";
                DOMName.value = ""
                fetch("http://" + document.location.host + "/api/create", {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ players: players, name: name }),
                })
                    .then((res) => {
                        m = res;
                        switch (res.status) {
                            case 409:
                                errCre.innerHTML = "ROOM IS TAKEN";
                        }
                    })
            }
    }
}
var m;
window.onload = () => {
    spaces.addEventListener("click", btnClicked)
    cre.addEventListener("click", creHandler)
    join.addEventListener("click", joinHandler)
}

