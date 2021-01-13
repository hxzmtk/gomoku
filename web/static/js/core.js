
let conn = undefined //保存websocket对象

function createRoom(){
    conn.send(JSON.stringify({
        "msgId":1,
        "body": {
            "test":"test"
        }
    }))
}

window.onload = function(){
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            // var messages = evt.data.split('\n');
            console.log(evt)
        };
        conn.onopen = function (){
          console.log("connected")
        };
    }
}