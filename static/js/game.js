"use strict"
let path = document.location.pathname
let id = path.slice(path.length - 8)
let waiting = document.getElementById("waitingRoom");
let socket = new WebSocket("ws://" + document.location.host + "/ws/" + id);
let chain = new chainReaction(15, 15, "red", socket);
// The bar above the players
let bar = document.getElementById("bar").getContext("2d");
// start is the only new variable
// true when game is ok to start.
let start = false
bar.canvas.width = chain.canvas.width;
// height is 10% of width
bar.canvas.height = bar.canvas.width * 0.05;

let changeBarC = (color) => {
    bar.fillStyle = color;
    bar.fillRect(0, 0, bar.canvas.width, bar.canvas.height);
}

chain.gr.onclick = (event) => {
    if (chain.state === true && start === true && chain.color === chain.mycolor) {
        let canvasObj = event.target.getBoundingClientRect();
        // Get the square coords clicked relative to
        let x = Math.floor((event.clientX - canvasObj.left) / chain.squareLength);
        let y = Math.floor((event.clientY - canvasObj.top) / chain.squareLength);
        if (chain.squares[y][x][-1] === this.mycolor || chain.squares[y][x][-1] === "") {
            // Prevent players from clicking other people's squares
            chain.socket.send(JSON.stringify({ type: "move", x: x, y: y, color: chain.color }))
        }
    }
}
chain.socket.onmessage = function (e) {
    let data = JSON.parse(e.data);
    switch (data.type) {
        case "start":
            start = true
            waiting.style.display = "none";
            chain.color = data.next
            changeBarC(data.next)
            break;
        case "move":
            chain.clicked(data.x, data.y);
            changeBarC(data.next);
            chain.color = data.next
            break;
        case "color":
            chain.mycolor = data.color
    }
};
chain.initBoard();
window.onload = () => {
    let start = document.querySelector(".ent");
    if (start !== null && start !== true) {
        start.onclick = () => {
            socket.send(JSON.stringify({ type: "start", val: true }));
        }
    }
}