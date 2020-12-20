"use strict"
let lay = document.querySelector(".layout");
let cre = document.querySelector(".create");
let join = document.querySelector(".join");
let spaces = document.querySelector(".spaces");

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
            break;
        case "ent":
            if (popup.id === "pop-join") {
                let room = document.getElementById("room");
                let pin = document.getElementById("pin");
            } else {
                let opt = document.getElementById("players");
            }
    }
}
window.onload = () => {
    spaces.addEventListener("click", btnClicked)
    cre.addEventListener("click", creHandler)
    join.addEventListener("click", joinHandler)
}

