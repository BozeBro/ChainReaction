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
    this.squares = [];
    this.canvas.onclick = (event) => {
      let canvasObj = event.target.getBoundingClientRect();
      let x = Math.floor((event.clientX - canvasObj.left) / this.squareLength);
      let y = Math.floor((event.clientY - canvasObj.top) / this.squareLength);
      this.move(x, y)
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
        this.ctx.beginPath();
        this.ctx.lineWidth = 1;
        this.ctx.strokeRect(l * this.squareLength, h * this.squareLength, this.squareLength, this.squareLength);
        this.ctx.closePath();
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
    if (start === undefined) {
      start = timestamp;
    }
    const elapsed = timestamp - start
    //ctx.clearRect(0, 0, canvas.width, canvas.height)
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
  move(x, y) {
    let curSquare = this.squares[y][x];
    this.ctx.save()
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
        this.move(x + 1, y)
      }
      if (x - 1 >= 0) { this.animate(0, posX, posY, -1 * this.squareLength / 60, 0); this.move(x - 1, y) }
      if (y + 1 < this.cols) { this.animate(0, posX, posY, 0, this.squareLength / 60); this.move(x, y + 1) }
      if (y - 1 >= 0) { this.animate(0, posX, posY, 0, -1 * this.squareLength / 60); this.move(x, y - 1) }
    }
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
}
var loc = function (z, length, push = 0) { return z * length + length / 2 + push }
let start;
chain = new chainReaction(10, 10);

chain.initBoard();


