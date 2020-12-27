"use strict"
/*
  See game.html for information on the event handlers
*/
class chainReaction {
  constructor(rows = 8, cols = 8, color, socket) {
    /*
    canvas, ctx handle dynamic movement
    stat handle objects that aren't moving
    gr, grctx is just the grid
    state tracks if an animation is taking place
      - blocks clicking event if false
    */
    this.__ms = 200 // length of entire animation in milliseconds. Meant to be a constant
    this.mycolor = ""
    this.color = color; // This is the color of the player's turn
    this.start = false
    this.socket = socket;
    this.canvas = document.getElementById("dynamic"); this.ctx = this.canvas.getContext("2d");
    this.statCtx = document.getElementById("static").getContext("2d");
    // grctx changes. Meant for constant display
    this.grctx = document.getElementById("grid").getContext("2d");
    this.rows = rows;
    this.cols = cols;
    this.squareLength = Math.min(450 / this.rows, 450 / this.cols); 
    this.squares = []; // Tells [number amount of circles, Exploding amount, cur color]
    this.state = true; // Tracks if an animation is taking place
  }
  initBoard() {
    // Make this.squares proper sizing
    // Make the visual board
    // Allows us to call initBoard() many times
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
        this.squares[h][l] = [0, val, ""];
      }
    }
  }
  clicked(x, y) {
    /* 
    x - int : x position of the square clicked
    y - int : y position of the square clicked
    -----
    Will Call explosion if square capacity is full
     */
    const curSquare = this.squares[y][x];
    const newVal = curSquare[0] + 1 < curSquare[1] ? 1 : 0;
    curSquare[0] *= newVal;
    curSquare[0] += newVal;
    this.draw(x, y, curSquare[0]);
    if (newVal === 0) { this.state = false; curSquare[2] = ""; this.explode([[x, y]]) }

  }
  async explode(exp) {
    /*
    exp - [[x, y]] : the squares that are going to explode
    -----
    Recursive function that animates each level of explosion
    Edit this.squares
    Disable this.state so no one can click
    Makes temporary global called this.exp to use in this.animate
    // Check if neighbors will explode.
    // Add coords and amount of circles of each square (For animation)
    // Change the color of neighboring squares (on this.squares)
    */
    if (exp.length === 0) {
      this.state = true
      return
    }
    const d = this.squareLength / this.__ms // Distance to move per frame
    let toAnimate = {
      "moved": [],
      "animations": [],
    }
    let expN = [];
    // Exploding Neighbors
    for (let [x, y] of exp) {
      for (let [dx, dy] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
        const posX = loc(x, this.squareLength)
        const posY = loc(y, this.squareLength)
        let nx = x + dx, ny = y + dy;
        if (0 <= nx && nx < this.rows && 0 <= ny && ny < this.cols) {
          // Add each neighbor of exploded to animation
          toAnimate.animations.push([posX, posY, dx, dy]);
          let curSquare = this.squares[ny][nx];
          curSquare[2] = this.color
          curSquare[0] += 1
          if (curSquare[0] >= curSquare[1]) {
            curSquare[0] = 0
            curSquare[2] = ""
            expN.push([nx, ny])
          }
          toAnimate.moved.push([nx, ny, curSquare[0]]);
        }
      }
    }
    this.exp = expN
    await new Promise(() =>
      requestAnimationFrame((ts) => this.animate(toAnimate, d, ts, ts)))
  }
  async animate(toAnimate, d, ts, start) {
    /*
    toAnimate - [{"moved": [], "animations": []}] ; Contains animation data
    toAnimate.moved : Tells what to draw on static Canvas
    toAnimate.animations : Instructs program how to animate (animation data)
    i - int : Tells the next frame of the animation
    d - int : How far each distance apart cicle should be drawn
    ----
    Animates each frame recursively.
    this.exp is used here and ONLY here.
    */
    const elapsed = ts - start
    const i = d * elapsed
    this.ctx.fillStyle = this.color
    this.ctx.lineWidth = 1
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
    for (let [x, y, dx, dy] of toAnimate.animations) {
      /*
      x - int : Current square's original x position
      y - int : Current square's original y position
      dx - int : unit x vector. Either 1, -1, or 0 (for no movement)
      dy - int : unit y vector. Either 1, -1, or 0 (for no movement)
      -----
      Movement will either be horizontal or vertical. NO diagonal
      */
      if (dx !== 0) {
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
      await new Promise(() =>
        requestAnimationFrame((ts) => this.animate(toAnimate, d, ts, start)))
    } else {
      // Everything's finished now. Redraw static screen
      for (let [x, y, v] of toAnimate.moved) {
        this.draw(x, y, v)
      }
      this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
      this.explode(this.exp)
      this.state = true

    }
  }
  draw(x, y, v) {
    /*
    x - int : y coordinate of the square
    y - int : X coordinate of the square
    v - int : Tells how many circles are in a square
    */
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
