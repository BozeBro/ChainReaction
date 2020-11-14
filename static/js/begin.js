const canvas = document.getElementById("chainReaction");
const ctx = canvas.getContext("2d");
var squares = [];
var squareLength = null

function initBoard(x, y) {
  height = 450;
  width = 450;

  squareLength = Math.min(width / x, height / y);

  ctx.canvas.width = x * squareLength;
  ctx.canvas.height = y * squareLength;
  ctx.fillStyle = "#fff";

  ctx.fillRect(0, 0, canvas.width, canvas.height);

  for (h = 0; h < y; h++) {
    squares[h] = [];
    for (l = 0; l < x; l++) {
      ctx.beginPath();
      ctx.lineWidth = 1;
      ctx.strokeRect(l * squareLength, h * squareLength, squareLength, squareLength);
      ctx.closePath();
      squares[h][l] = 0;
    }
  }
  return squareLength
};
function makeCircle() {
  ctx.beginPath();
  ctx.fillStyle = "#f00";
  ctx.lineWidth = 1;
  ctx.fillRect(1*squareLength + squareLength / 3, 1 * squareLength + squareLength / 3, squareLength/3, squareLength/3);
  //ctx.fillRect()
  ctx.closePath();
}
squareLength = initBoard(20, 20);
makeCircle();