"use strict"
class chainReaction {
  constructor(rows = 8, cols = 8, color) {
    /*
    canvas, ctx handle dynamic movement
    stat handle objects that aren't moving
    gr, grctx is just the grid
    state tracks if an animation is taking place
      - blocks clicking event if false
    */
    this.canvas = document.getElementById("dynamic");
    this.ctx = this.canvas.getContext("2d");
    this.statCanv = document.getElementById("static");
    this.statCtx = this.statCanv.getContext("2d");
    this.gr = document.getElementById("grid");
    this.grctx = this.gr.getContext("2d");
    this.socket = new WebSocket("ws://" + document.location.host + "/ws");
    this.rows = rows;
    this.cols = cols;
    this.squareLength = Math.min(450 / rows, 450 / cols);
    this.squares = [];
    this.state = true; // Tracks if an animation is taking place
    this.color = color;
    this.gr.onclick = (event) => {
      if (this.state === true) {
        let canvasObj = event.target.getBoundingClientRect();
        // Get the square coords clicked relative to
        let x = Math.floor((event.clientX - canvasObj.left) / this.squareLength);
        let y = Math.floor((event.clientY - canvasObj.top) / this.squareLength);
        this.socket.send(JSON.stringify({x:x, y:y, color:this.color}))
      }
    }
  }
  initBoard() {
    this.statCtx.canvas.width = this.ctx.canvas.width = this.rows * this.squareLength;
    this.statCtx.canvas.height = this.ctx.canvas.height = this.cols * this.squareLength;
    this.grctx.canvas.width = this.ctx.canvas.width; this.grctx.canvas.height = this.ctx.canvas.height;
    this.ctx.fillStyle = "#fff"; // White
    this.grctx.lineWidth = 1;
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    for (let h = 0; h < this.cols; h++) {
      this.squares[h] = [];
      for (let l = 0; l < this.rows; l++) {
        this.grctx.beginPath();
        this.grctx.strokeRect(
          l * this.squareLength,
          h * this.squareLength,
          this.squareLength,
          this.squareLength
        );
        this.grctx.closePath();
        let val = (function (rows, cols) {
          let valid = 0;
          if (l - 1 >= 0) { valid += 1 }
          if (l + 1 < rows) { valid += 1 }
          if (h - 1 >= 0) { valid += 1 }
          if (h + 1 < cols) { valid += 1 }
          return valid;
        })(this.rows, this.cols);
        this.squares[h][l] = [0, val];
      }
    }
  }
  explode(exp) {
    // exp is the items that will explode. A stack
    const d = 70 * this.squareLength / 1000 // Distance to move per frame
    let toAnimate = {
      "moved": [],
      "animations": [],
    }
    while (exp.length !== 0) {
      let info = {
        "moved": [],
        "animations": [],
      }
      for (let [x, y] of exp) {
        // Call move() on neighbors
        const posX = loc(x, this.squareLength)
        const posY = loc(y, this.squareLength)
        if (x + 1 < this.rows) {
          info.animations.push([posX, posY, 1, 0])
        }
        if (x - 1 >= 0) {
          info.animations.push([posX, posY, -1, 0])
        }
        if (y + 1 < this.cols) {
          info.animations.push([posX, posY, 0, 1])
        }
        if (y - 1 >= 0) {
          info.animations.push([posX, posY, 0, -1])
        }
      }
      exp = this.move(exp, info)
      toAnimate.animations.push(info.animations)
      toAnimate.moved.push(info.moved)
    }
    const ani = async () => {
      return new Promise(() =>
        requestAnimationFrame(() => this.animate(toAnimate, -d, d, 0)))
    }
    return ani()
  }
  clicked(x, y) {
    const curSquare = this.squares[y][x];
    const d = curSquare[0] + 1 < curSquare[1] ? 1 : 0;
    this.squares[y][x][0] *= d;
    this.squares[y][x][0] += d;
    this.draw(x, y, this.squares[y][x][0]);
    if (d === 0) { this.state = false; this.explode([[x, y]]) }

  }
  move(exp, info) {
    // Check if neighbors will explode.
    // Add coords and amount of circles of each square (For animation)
    let expN = [];
    for (let [x, y] of exp) {
      for (let [dx, dy] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
        let nx = x + dx, ny = y + dy;
        if (0 <= nx && nx < this.rows && 0 <= ny && ny < this.cols) {
          let curSquare = this.squares[ny][nx];
          const d = curSquare[0] + 1 < curSquare[1] ? 1 : 0;
          this.squares[ny][nx][0] *= d;
          this.squares[ny][nx][0] += d;
          info.moved.push([nx, ny, this.squares[ny][nx][0]]);
          if (d === 0) { expN.push([nx, ny]) }
        }
      }
    }
    return expN
  }
  animate(toAnimate, i, d, ind) {
    i += d
    this.ctx.fillStyle = this.color
    this.ctx.lineWidth = 1
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
    for (let [x, y, dx, dy] of toAnimate.animations[ind]) {
      if (dx !== 0) {
        this.ctx.beginPath();
        this.ctx.arc(x + i * dx, y, this.squareLength / 4, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
      } else {
        this.ctx.beginPath();
        this.ctx.arc(x, y + i * dy, this.squareLength / 4, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
      }
    }
    if (Math.abs(i) < this.squareLength) {
      // Complete animation 
      return new Promise(() => requestAnimationFrame(() => this.animate(toAnimate, i, d, ind)))
    }
    else if (ind + 1 < toAnimate.moved.length) {
      // Go to next set to animate
      for (let [x, y, v] of toAnimate.moved[ind]) {
        this.draw(x, y, v)
      }
      return new Promise(() =>
        requestAnimationFrame(() => this.animate(toAnimate, -d, d, ind + 1)))
    } else {
      // Everything's settled now. Redraw stat screen
      for (let [x, y, v] of toAnimate.moved[ind]) {
        this.draw(x, y, v)
      }
      this.state = true;
      return new Promise(() =>
        requestAnimationFrame(() => this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)))
    }
  }
  draw(x, y, v) {
    let circlePos = this.squareLength / 7.5;
    let radius = this.squareLength / 4;
    this.statCtx.fillStyle = this.color;
    this.statCtx.lineWidth = 1;
    switch (v) {
      // Handles the current circle count in a square
      case 1:
        this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
        break;
      case 2:
        this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
        break;
      case 3:
        this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
        break;
      default:
        // Clear the circles
        this.statCtx.beginPath();
        this.statCtx.clearRect(x * this.squareLength, y * this.squareLength, this.squareLength, this.squareLength);
        this.statCtx.closePath();
        break;
    }
  }
}
const loc = function (z, length, offset = 0) { return z * length + length / 2 + offset }
//let chain = new chainReaction(15, 15, "red");
//chain.initBoard();


