var chainReaction = function (rows = 8, cols = 8) {
  this.board = document.getElementById("game");
  this.rows = rows;
  this.cols = cols;
  this.gameBoard = [];
}

chainReaction.prototype.makeBoard = function () {
  const squareLength = 450 / Math.max(this.rows, this.cols);
  for (c = 0; c < this.cols; c++) {
    var column = document.createElement("div");
    column.className = "column";
    for (r = 0; r < this.rows; r++) {
      var row = document.createElement("div");
      row.className = "row";
      row.style.height = String(squareLength) + "px";
      row.style.width = String(squareLength) + "px";
      row.id = String(r) + ":" + String(c);
      row.addEventListener("click", function () {
        squareClick(this.id)
      })
      column.appendChild(row);
    }
    this.board.appendChild(column);
  }
}
chainReaction.prototype.makeCodedBoard = function () {
  gameBoard = [];
  for (c = 0; c < this.cols; c++) {
    gameBoard[c] = []
    for (r = 0; r < this.rows; r++) {
      valid = 0;
      if (c - 1 >= 0) {
        valid += 1;
      }
      if (c + 1 < this.cols) {
        valid += 1;
      }
      if (r - 1 >= 0) {
        valid += 1;
      }
      if (r + 1 < this.rows) {
        valid += 1;
      }
      gameBoard[c][r] = valid;
    }
  }
  return gameBoard;
}
// Set move to null, so i can recursively call it
chainReaction.prototype.move = null
chainReaction.prototype.move = function (x, y, color) {
  // Player move
  let square = this.board.children[y].children[x];
  value = square.innerHTML;
  square.style.background = color;
  square.style.color = "white";
  if (value === "") {
    square.innerHTML = "1";
  } else if (parseInt(value) + 1 < this.gameBoard[y][x]) {
    square.innerHTML = String(parseInt(value) + 1);
  } else {
    if (x - 1 >= 0) {
      this.move(x - 1, y, color);
    }
    if (x + 1 < this.rows) {
      this.move(x + 1, y, color);
    }
    if (y - 1 >= 0) {
      this.move(x, y - 1, color);
    }
    if (y + 1 < this.cols) {
      this.move(x, y + 1, color);
    }
    square.style.background = "white";
    square.innerHTML = "";
  }
}
function squareClick(id) {
  // Will Try to make a move on the clicked square
  alert(id);
}
var chain = new chainReaction(8, 8);


chain.makeBoard();
chain.gameBoard = chain.makeCodedBoard();
chain.move(0, 0, "black");
chain.move(0, 0, "black");



