const canvas = document.getElementById("chainReaction");
const ctx = canvas.getContext("2d");
function chainReaction(x, y) {
  this.x = x;
  this.y = y;
  this.squareLength = Math.min(450 / x, 450 / y);
  this.squares = [];
}
chainReaction.prototype.initBoard = function () {

  ctx.canvas.width = this.x * this.squareLength;
  ctx.canvas.height = this.y * this.squareLength;
  ctx.fillStyle = "#fff"; // White

  ctx.fillRect(0, 0, canvas.width, canvas.height);

  for (h = 0; h < this.y; h++) {
    this.squares[h] = [];
    for (l = 0; l < this.x; l++) {
      ctx.beginPath();
      ctx.lineWidth = 1;
      ctx.strokeRect(l * this.squareLength, h * this.squareLength, this.squareLength, this.squareLength);
      ctx.closePath();
      val = (function () {
        valid = 0;
        if (l - 1 >= 0) { valid += 1 }
        if (l + 1 < this.x) { valid += 1 }
        if (h - 1 >= 0) { valid += 1 }
        if (h + 1 < this.y) { valid += 1 }
        return valid
      })();
      this.squares[h][l] = [0, val];
    }
  }
}
chainReaction.prototype.move = function(x, y) {
  ctx.beginPath();
  ctx.fillStyle = "#f00";
  ctx.lineWidth = 1;
  ctx.fillRect(x * this.squareLength + this.squareLength / 3, y * this.squareLength + this.squareLength / 3, this.squareLength / 3, this.squareLength / 3);
  //ctx.fillRect()
  ctx.closePath();
}
chain = new chainReaction(8, 8);
chain.initBoard(10, 15);
chain.move(1, 1);
ctx.clearRect(0, 0, chain.squareLength * 4, chain.squareLength * 4);