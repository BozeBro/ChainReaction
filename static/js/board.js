"use strict";
//  See game.html for information on the event handlers
//  canvas, ctx handle dynamic movement
//  stat handle objects that aren't moving
//  gr, grctx is just the grid
class chainReaction {
  constructor(rows = 8, cols = 8, color) {
    this.__ms = 200 // length of entire animation in milliseconds. Meant to be a constant
    this.mycolor = "" // Player color is initialized at "start" JSON
    this.color = color; // This is the color of current player's turn
    this.start = false // stops anyone from clicking the screen until the game starts
    this.canvas = document.getElementById("dynamic"); this.ctx = this.canvas.getContext("2d");
    this.statCtx = document.getElementById("static").getContext("2d");
    // grctx changes. Meant for constant display
    this.grctx = document.getElementById("grid").getContext("2d");
    this.rows = rows;
    this.cols = cols;
    this.squareLength = Math.min(screen.height * .80 / this.rows, screen.height * .80 / this.cols);
    this.squares = []; // Tells [number amount of circles, Exploding amount, cur color]
    this.msgHandler = {}; // map of functions that will handle a response based on resp type
    this.que = [] // que system to asynchronously handle responses from websockets. Incoming response should not occur after ongoing animation.
    this.state = true; // Tracks if an animation is taking place
  };
  initBoard() {
    // Make the visual board within boundary of 450px
    // Allows us to call initBoard() many times, for each time we start a game.
    this.statCtx.canvas.width = this.ctx.canvas.width = this.rows * this.squareLength;
    this.statCtx.canvas.height = this.ctx.canvas.height = this.cols * this.squareLength;
    this.grctx.canvas.width = this.ctx.canvas.width; this.grctx.canvas.height = this.ctx.canvas.height;
    this.ctx.fillStyle = "#fff"; // White
    this.grctx.lineWidth = 1;
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    for (let h = 0; h < this.cols; h++) {
      for (let l = 0; l < this.rows; l++) {
        this.grctx.beginPath();
        this.grctx.strokeRect(
          l * this.squareLength,
          h * this.squareLength,
          this.squareLength,
          this.squareLength
        );
        this.grctx.closePath();
      }
    }
  }
  //  Animates each frame recursively.
  //  animations - Instructs program how to animate (animation data)
  //  unmoving - Tells what to draw on static Canvas. Position after each index of animation.
  //  ts - timestamp given by requestAnimationFrame.
  //  start - starting point of animation. Allows us to know how much time has passed relative to the starting point.
  //  ind - tells index of animation and unmoving array.
  //  color - tells color the circles should be.
  animate(animations, unmoving, ts, start, ind, color) {
    const d = this.squareLength / this.__ms;
    const elapsed = ts - start;
    const i = d * elapsed;
    this.ctx.fillStyle = this.color;
    this.ctx.lineWidth = 1;
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    for (let [x, y, dx, dy] of animations[ind]) {
      x = loc(x, this.squareLength);
      y = loc(y, this.squareLength);
      /*
      x - int : Current square's position relative to canvas
      y - int : Current square's position relative to canvas
      dx - int : unit x vector. Either 1, -1, or 0 (for no movement)
      dy - int : unit y vector. Either 1, -1, or 0 (for no movement)
      -----
      Movement will either be horizontal or vertical. NO diagonals
      */
      if (dx !== 0) {
        // If dx is 0, then dy must be moving
        this.ctx.beginPath();
        this.ctx.arc(x + Math.min(i, this.squareLength) * dx, y, this.squareLength / 4, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
      } else {
        this.ctx.beginPath();
        this.ctx.arc(x, y + Math.min(i, this.squareLength) * dy, this.squareLength / 4, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
        this.ctx.closePath();
      }
    }
    if (elapsed < this.__ms) {
      // Complete rest of animation
      requestAnimationFrame((ts) => this.animate(animations, unmoving, ts, start, ind, color))
    } else if (ind + 1 < animations.length) {
      // Complete next level of explosion / animation
      for (let [x, y, v] of unmoving[ind+1]) {
        this.draw(x, y, v)
      }
      this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
      requestAnimationFrame((ts) => this.animate(animations, unmoving, ts, ts, ind+1, color))
    } else {
      // Draw the last unmoving square. Clear screen.
      for (let [x, y, v] of unmoving[ind+1]) {
        this.draw(x, y, v)
      }
      this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
      this.color = color;
      changeTurn(color);
      changeBarC(color);
      this.state = true
      // First element should be the move response. Get the next one.
      // delay the shift till the very end of async, so we do not start any new functions.
      chain.que.shift()
      if (this.que.length > 0) {
        let data = this.que[0]

        chain.msgHandler[data.type.toLowerCase()](data)
      }
      cancelAnimationFrame(ts)
    }
  }

  // Draws certain amount of circles on the board
  //  x - y coordinate of the square
  //  y - X coordinate of the square
  //  v - Tells how many circles are in a square
  draw(x, y, v) {
    let circlePos = this.squareLength / 7.5;
    let radius = this.squareLength / 4;
    this.statCtx.fillStyle = this.color;
    this.statCtx.lineWidth = 1;
    const draw1 = () => {
      this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength, -1 * circlePos), loc(y, this.squareLength, -1 * circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
    }
    const draw2 = () => {
      draw1()
      this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength, circlePos), loc(y, this.squareLength, -1 * circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
    }
    const draw3 = () => {
      draw2()
      this.statCtx.beginPath();
        this.statCtx.arc(loc(x, this.squareLength), loc(y, this.squareLength, circlePos), radius, 0, 2 * Math.PI);
        this.statCtx.stroke();
        this.statCtx.fill();
        this.statCtx.closePath();
    }
    const erase = () => {
      this.statCtx.beginPath();
        this.statCtx.clearRect(x * this.squareLength, y * this.squareLength, this.squareLength, this.squareLength);
        this.statCtx.closePath();
    }
    switch (v) {
      // Handles the current circle count in a square
      case 1:
        erase()
        draw1()
        break;
      case 2:
        erase()
        draw2()
        break;
      case 3:
        erase()
        draw3()
        break;
      default:
        // Clear the circles
        erase()
    }
  }
}

// Helper function to find position on the board
const loc = function (z, length, offset = 0) { return z * length + length / 2 + offset }
// The bar is the element above the board
let bar = document.getElementById("bar").getContext("2d");
// Tells if a game has started
let display = document.getElementById("turn-display");
function changeTurn(username, color) {
	display.innerHTML = `It is ${username}'s turn. The color is ${color}.`
  }
let changeBarC = (color) => {
    bar.fillStyle = color;
    bar.fillRect(0, 0, bar.canvas.width, bar.canvas.height);
}
