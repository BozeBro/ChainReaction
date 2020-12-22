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
console.log("NEW CAHNGE NEVER SEEN BEFORE")
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
            console.log("exit")
            popup.style.display = "none";
            popup = null;
            errCre.innerHTML = "";
            errJoin.innerHTML = "";
            break;
        case "ent":
            if (popup.id === "pop-join") {
                console.log("join")
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
                console.log("sent")
                // Creating a room
                let DOMName = document.getElementById("name");
                let room = DOMName.value;
                let players = document.getElementById("players").value;
                if (!room) {
                    errCre.innerHTML = "ENTER A ROOM NAME";
                    return
                }
                DOMName.value = ""
                fetch("http://" + document.location.host + "/api/create", {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ Players: players, Room: room }),
                })
                    .then(res => {
                        switch (res.status) {
                            case 409:
                                errCre.innerHTML = "ROOM IS TAKEN";
                            case 400:
                                errCre.innerHTML = "WHAT ARE YOU DOING?";
                        }
                        //return res
                    })
                /*.then(res => {
                    if (!res.redirected) return;
                    window.location.replace(res.url);

                }) */
            }
    }
}
window.onload = () => {
    spaces.addEventListener("click", btnClicked)
    cre.addEventListener("click", creHandler)
    join.addEventListener("click", joinHandler)
}

