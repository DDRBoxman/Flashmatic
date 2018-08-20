let mouseUpFunction = (id) => {
  nanoajax.ajax({
    url: '/key',
    method: 'POST',
    body: JSON.stringify({"key_id": id})
  }, function (code, responseText, request) {

  })
};

let buttonsDiv = document.getElementById('buttons');


for (i = 0; i < 15; i++) {
  let button = document.createElement("img");
  button.id = "button_" + i;
  button.src = "data:image/gif;base64,R0lGODlhAQABAAD/ACwAAAAAAQABAAACADs=";
  let id = i;
  button.addEventListener("mouseup", () => {
    mouseUpFunction(id);
  });
  buttonsDiv.appendChild(button);
}

nanoajax.ajax({
  url: '/display', method: 'GET'
}, function (code, responseText, request) {
  const display = JSON.parse(responseText);
  display.forEach((icon) => {
    const button = document.getElementById("button_" + icon.key_id);
    button.src = "/icon/" + icon.icon_path;
  });
});