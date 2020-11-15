var board = document.getElementById("game");
function makeBoard(rows, cols) {
  const squareLength = 450 / Math.max(rows, cols);
  for (c = 0; c < cols; c++) {
    var column = document.createElement("div");
    column.className = "column";
    for (r = 0; r < rows; r++) {
      var row = document.createElement("div");
      row.className = "row";
      row.id = "test";
      row.style.height = String(squareLength) + "px";
      row.style.width = String(squareLength) + "px";
      column.appendChild(row);
    }
    board.appendChild(column);
  }
}
function makeGameBoard(rows, cols) {
  gameBoard = [];
  for (c = 0; c < cols; c++) {
    gameBoard[c] = []
    for (r = 0; r < rows; r++) {
      valid = 0;
      if (c - 1 >= 0) {
        valid += 1;
      }
      if (c + 1 < cols) {
        valid += 1;
      }
      if (r - 1 >= 0) {
        valid += 1;
      }
      if (r + 1 < rows) {
        valid += 1;
      }
      gameBoard[c][r] = valid;
    }
  }
  return gameBoard;
}
const move = function (x, y, color) {
  // Player move
  let square = board.children[y].children[x];
  square.style.background = color;
  square.style.color = "white";
  if (square.innerHTML === "") {
    square.innerHTML = "1";
  } else if (notFilled(x, y, square)) {
    square.innerHTML = String(parseInt(square.innerHTML) + 1);
  } else {
    explode(x, y, color);
    square.style.background = "white";
    square.innerHTML = "";
  }
}
const notFilled = function (x, y, square) {
  if (parseInt(square.innerHTML) + 1 >= gameBoard[y][x]) {
    return false;
  }
  return true;
}
const explode = function (x, y, color) {
  if (x - 1 >= 0) {
    move(x - 1, y, color);
  }
  if (x + 1 < rows) {
    move(x + 1, y, color);
  }
  if (y - 1 >= 0) {
    move(x, y - 1, color);
  }
  if (y + 1 < cols) {
    move(x, y + 1, color);
  }
}
document.getElementById("game").onclick = function (event) {
  alert(event.region)
  event.preventDefault();
}
const rows = 25, cols = 25;
makeBoard(rows, cols);
var gameBoard = makeGameBoard(rows, cols);
move(0, 0, "black");
move(1, 0, "black");
move(1, 0, "black");
move(1, 0, "black");

document.getElementById("test").onclick = function (event) {
  alert(event.region)
  event.preventDefault();
}