const canvas = document.getElementById("chainReaction");
const ctx = canvas.getContext("2d");

class chainReaction {
  constructor(rows = 8, cols = 8) {
    this.rows = rows;
    this.cols = cols;
    this.squareLength = Math.min(450 / rows, 450 / cols);
    this.radius = this.squareLength / 4
    this.squares = [];
  }
  initBoard() {

    ctx.canvas.width = this.rows * this.squareLength;
    ctx.canvas.height = this.cols * this.squareLength;
    ctx.fillStyle = "#fff"; // White

    ctx.fillRect(0, 0, canvas.width, canvas.height);

    for (let h = 0; h < this.cols; h++) {
      this.squares[h] = [];
      for (let l = 0; l < this.rows; l++) {
        ctx.beginPath();
        ctx.lineWidth = 1;
        ctx.strokeRect(l * this.squareLength, h * this.squareLength, this.squareLength, this.squareLength);
        ctx.closePath();
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
  animate(timestamp, x, y, dx, dy, i=0) {
    if (start === undefined) {
      start = timestamp;
    }
    const elapsed = timestamp - start
    //ctx.clearRect(0, 0, canvas.width, canvas.height)
    if (dx !== 0) {
      i += dx
      ctx.beginPath();
      ctx.restore()
      ctx.arc(x + i, y, this.squareLength / 4, 0, 2 * Math.PI);
      ctx.stroke();
      ctx.closePath();
    } else {
      i += dy
      ctx.beginPath();
      ctx.arc(x, y + i, this.squareLength / 4, 0, 2 * Math.PI);
      ctx.stroke();
      ctx.closePath();
    }
    if (Math.abs(i) < this.squareLength) {
      requestAnimationFrame(() => this.animate(timestamp, x, y, dx, dy, i))
    } else { start = undefined }
  }
  move(x, y) {
    let curSquare = this.squares[y][x];
    ctx.save()
    if (curSquare[0] + 1 < curSquare[1]) {
      this.squares[y][x][0] += 1;
      this.draw(x, y);
    } else {
      this.squares[y][x][0] = 0;
      this.draw(x, y);
      let posX = loc(x, this.squareLength);
      let posY = loc(y, this.squareLength)
      if (x + 1 < this.rows) {
        this.animate(0, posX, posY, this.squareLength / 60, 0)
        this.move(x+1,y)
      }
      if (x - 1 >= 0) { this.animate(0, posX, posY, -1 * this.squareLength / 60, 0); this.move(x-1,y) }
      if (y + 1 < this.cols) { this.animate(0, posX, posY, 0, this.squareLength / 60);this.move(x,y+1) }
      if (y - 1 >= 0) { this.animate(0, posX, posY, 0, -1 * this.squareLength / 60);this.move(x,y-1) }
    }
  }
  draw(x, y) {
    let curSquare = this.squares[y][x];
    let circlePos = this.squareLength / 7.5;
    ctx.fillStyle = "#f00";
    ctx.lineWidth = 1;
    switch (curSquare[0]) {
      // Handles the current circle count in a square
      case 1:
        ctx.beginPath();
        ctx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius, 0, 2 * Math.PI);
        ctx.stroke();
        ctx.fill();
        ctx.closePath();
        break;
      case 2:
        ctx.beginPath();
        ctx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius, 0, 2 * Math.PI);
        ctx.stroke()
        ctx.fill();
        ctx.closePath();
        break;
      case 3:
        ctx.beginPath();
        ctx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), this.radius, 0, 2 * Math.PI);
        ctx.stroke();
        ctx.fill();
        ctx.closePath();
        break;
      case 0:
        // Clear the circles
        let lw = this.radius / 4 //handle linewidth. 4 is arbitrary
        ctx.beginPath();
        ctx.globalCompositeOperation = "destination-out";
        ctx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius + lw, 0, 2 * Math.PI);
        ctx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius + lw, 0, 2 * Math.PI);
        ctx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), this.radius + lw, 0, 2 * Math.PI);
        ctx.fill();
        ctx.globalCompositeOperation = "source-over";
        ctx.closePath();
        break;
    }
  }
}

let start;
betaX = 400;
betaY = 300;
//chainReaction.prototype.animate = function () {};


var loc = function (z, length, push = 0) { return z * length + length / 2 + push }
var ani;
function beginAni() {
  ani = "3"
}
chain = new chainReaction(10, 10);
canvas.onclick = (event) => {
  let canvasObj = event.target.getBoundingClientRect();
  let x = Math.floor((event.clientX - canvasObj.left) / chain.squareLength);
  let y = Math.floor((event.clientY - canvasObj.top) / chain.squareLength);
  chain.move(x, y)
}

chain.initBoard();


