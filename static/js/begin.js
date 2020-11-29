class chainReaction {
  constructor(rows = 8, cols = 8) {
    this.canvas = document.getElementById("chainReaction");
    this.ctx = this.canvas.getContext("2d");
    this.staticCanv = document.getElementById("static");
    this.staticCtx = this.staticCanv.getContext("2d");
    this.rows = rows;
    this.cols = cols;
    this.squareLength = Math.min(450 / rows, 450 / cols);
    this.radius = this.squareLength / 4
    this.squares = []
    this.exp = []
    this.staticCanv.onclick = (event) => {
      let canvasObj = event.target.getBoundingClientRect();
      let x = Math.floor((event.clientX - canvasObj.left) / this.squareLength);
      let y = Math.floor((event.clientY - canvasObj.top) / this.squareLength);
      this.clicked(x, y)
    }
  }
  initBoard() {
    this.staticCtx.width = this.ctx.canvas.width = this.rows * this.squareLength;
    this.staticCtx.height = this.ctx.canvas.height = this.cols * this.squareLength;
    this.ctx.fillStyle = "#fff"; // White

    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    for (let h = 0; h < this.cols; h++) {
      this.squares[h] = [];
      for (let l = 0; l < this.rows; l++) {
        this.staticCtx.beginPath();
        this.staticCtx.lineWidth = 1;
        this.staticCtx.strokeRect(l * this.squareLength, h * this.squareLength, this.squareLength, this.squareLength);
        this.staticCtx.closePath();
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
  animate(timestamp, x, y, dx, dy, i = 0) {
    start = timestamp ? start === undefined : start
    const elapsed = timestamp - start
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
    if (dx !== 0) {
      i += dx
      this.ctx.beginPath();
      this.ctx.arc(x + i, y, this.squareLength / 4, 0, 2 * Math.PI);
      this.ctx.stroke();
      this.ctx.closePath();
    } else {
      i += dy
      this.ctx.beginPath();
      this.ctx.arc(x, y + i, this.squareLength / 4, 0, 2 * Math.PI);
      this.ctx.stroke();
      this.ctx.closePath();
    }
    if (Math.abs(i) < this.squareLength) {
      requestAnimationFrame(() => this.animate(timestamp, x, y, dx, dy, i))
    } else { start = undefined }
  }
  explode(exp) {
    // 70 frames per second
    let d = 70 * this.squareLength / 1000 
    while (exp.length !== 0) {
      let moved = []
      for (const e of exp) {
        let x = e[0]; let y = e[1];
        let posX = loc(x, this.squareLength)
        let posY = loc(y, this.squareLength)
        if (x + 1 < this.rows) {
          moved.push([x + 1, y])
          this.animate(0, posX, posY, d, 0)
        }
        if (x - 1 >= 0) {
          moved.push([x - 1, y])
          this.animate(0, posX, posY, -d, 0)
        }
        if (y + 1 < this.cols) {
          moved.push([x, y + 1])
          this.animate(0, posX, posY, 0, d)
        }
        if (y - 1 >= 0) {
          moved.push([x, y - 1])
          this.animate(0, posX, posY, 0, -d)
        }
      }
      exp = this.move(moved)
    }
  }
  clicked(x, y) {
    let curSquare = this.squares[y][x];
    this.draw(x, y)
    if (curSquare[0] + 1 < curSquare[1]) {
      this.squares[y][x][0] += 1;
    } else {
      // Max capacity is reached
      this.squares[y][x][0] = 0;
      this.explode([x, y])
    }
  }
  move(neighbors) {
    let exp = []
    for (const neighbor of neighbors) {
      let x = neighbor[0]
      let y = neighbor[1]
      let curSquare = this.squares[y][x]
      if (curSquare[0] + 1 < curSquare[1]) {
        this.squares[y][x][0] += 1
      } else {
        this.squares[y][x][0] = 0
        exp.push([x, y])
      }
      this.draw(x, y)
    }
    return exp
  }
  draw(x, y) {
    let curSquare = this.squares[y][x];
    let circlePos = this.squareLength / 7.5;
    this.ctx.fillStyle = "#f00";
    this.ctx.lineWidth = 1;
    switch (curSquare[0]) {
      // Handles the current circle count in a square
      case 1:
        this.ctx.beginPath();
        this.ctx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
        break;
      case 2:
        this.ctx.beginPath();
        this.ctx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius, 0, 2 * Math.PI);
        this.ctx.stroke()
        this.ctx.fill();
        this.ctx.closePath();
        break;
      case 3:
        this.ctx.beginPath();
        this.ctx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), this.radius, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
        break;
      case 0:
        // Clear the circles
        let lw = this.radius / 4 //handle linewidth. 4 is arbitrary
        this.ctx.beginPath();
        this.ctx.globalCompositeOperation = "destination-out";
        this.ctx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius + lw, 0, 2 * Math.PI);
        this.ctx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), this.radius + lw, 0, 2 * Math.PI);
        this.ctx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), this.radius + lw, 0, 2 * Math.PI);
        this.ctx.fill();
        this.ctx.globalCompositeOperation = "source-over";
        this.ctx.closePath();
        break;
    }
  }
  check(x, y, f, ...args) {
    if (x + 1 < this.rows) { f(...args) }
    if (x - 1 >= 0) { f(-1, ...args) }
    if (y + 1 < this.cols) { f(...args) }
    if (y - 1 >= 0) { f(...args) }

  }
}
var loc = function (z, length, offset = 0) { return z * length + length / 2 + offset }
var check = function (x, y) {

}
let start;
chain = new chainReaction(10, 10);

chain.initBoard();


