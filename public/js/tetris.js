// drawing
const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

ws = new WebSocket(`ws://${window.location.toString().substr(7)}/ws`);

ws.onopen = function (e) {
  console.log(e);
};

ws.onclose = function (e) {
  console.log(e);
  alert("connection closed");
};

ws.onmessage = function (e) {
  console.log("recieved board");
  var board = JSON.parse(e.data);
  var grid = board.Grid;

  var y = 0;
  var x = 0;

  for (; x <= 400 && y <= 800; x += 40) {
    if (x == 400) {
      x = 0;
      y += 40;
    }
    if (y == 800) {
      return;
    }

    console.log(y / 40, x / 40);

    switch (grid[y / 40][x / 40]) {
      case 0:
        ctx.fillStyle = "#FFFFFF";
        break;
      default:
        ctx.fillStyle = "#000000";
        break;
    }
    ctx.fillRect(x, y, 40, 40);
  }
};
