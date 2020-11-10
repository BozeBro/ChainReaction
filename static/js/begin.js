var chainReaction = function() {
  var board = new Array()
  for (h=0;h<height;h++) {
    board.push(new Array(width))
    for (w=0;w<width;w++) {
      board[h][w] = w
    }
  }
  return board
}
function move(board, x, y, p) {
  if (board[y][x][1] === p || board[y][x][1] === null) {
    board[y][x][0] += 1
  }
}
(function() {
  var board = initBoard(5, 5);
  console.log(board)

}
)();