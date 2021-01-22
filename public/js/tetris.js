// drawing
const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

ws = new WebSocket(`ws://${window.location.toString().substr(7)}/ws`);

var images = [];

for (let i = 0; i < 8; i++) {
  let img = document.createElement("img");
  img.src = `../../images/pieces/tetriscube${i}.png`;
  images.push(img);
}

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

    ctx.drawImage(images[grid[y / 40][x / 40]], x, y);
  }
};

window.addEventListener("keypress", function (e) {
  console.log(e.code);
  var response = {
    Misc: "",
    Key: e.code,
  };

  ws.send(JSON.stringify(response));
});
