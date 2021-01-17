
let hand = {
    "nilHand": 0,
    "blackHand": 1,
    "whiteHand": 2
}

let msgId = {
    'connect': 99999,
    'error': 0,
    'listRoom': 1,
    'createRoom' : 2,
    'joinRoom': 3,
    'chessboardWalk': 4
}

let msgAck = {}

let conn = undefined //保存websocket对象


function connect(){
    conn.send(JSON.stringify({
        "msgId":msgId.connect,
        "body": {
            "username": sessionStorage.getItem("un")
        }
    }))
}

function listRoom() {
    conn.send(JSON.stringify({
        "msgId":msgId.listRoom,
        "body": {}
    }))
}

function createRoom(){
    conn.send(JSON.stringify({
        "msgId":msgId.createRoom,
        "body": {}
    }))
}

//初始化二维数组
function initPlace(row, col) {
    place = Array(row).fill(0).map(x => Array(col).fill(0));
}

//使刚落下的棋子闪烁,用于提示
function remain(x,y){
    if (last_pos.x != -1  && last_pos.y != -1){
        $(`#go-${last_pos.x}-${last_pos.y}`).removeClass("chess-spinner");
        $(`#go-${x}-${y}`).addClass("chess-spinner");
    }
    $(`#go-${x}-${y}`).addClass("chess-spinner");
    last_pos.x = x, last_pos.y = y;
}

//生成棋盘
function generate_board(row, col){

    let board = ""
    for (let i = 0; i < row; i++) {
        let tmp = ""
        for (let j = 0; j < col; j++) {
            tmp += `<i class="i-nomal" id="go-${i}-${j}"></i>`
        }
        tmp +="<br>"
        board += tmp
    }
    document.getElementById("go-board").innerHTML =  board
    initPlace(row,col);
}

function handle(event) {
    try {
        let msg = JSON.parse(event.data)
        if (!msg.hasOwnProperty("msgId")){
            console.log("invalid msg:",msg)
            return
        }
        switch (msg.msgId) {
            case -msgId.connect:
                sessionStorage.setItem("un",msg.username)
                break;
            case -msgId.listRoom:
                let tmp = ""
                msg.data.forEach(element => {
                    element.enemy = element.enemy == "" ? "无":element.enemy
                    isDisabled = element.isFull == true ? "disabled": ""
                    tmp += `<tr>
                    <th scope="row">${element.roomId}</th>
                    <td>${element.master}</td>
                    <td>${element.enemy}</td>
                    <td><button type="button" class="btn btn-sm btn-primary"  ${isDisabled}>加入</button></td>
                  </tr>`
                });
                document.getElementById("dating-data").innerHTML = tmp
                break;
            case -msgId.createRoom:
                generate_board(15,15)
                document.getElementById("dating").classList.add("d-none")
                document.getElementById("room").classList.remove("d-none")
                document.getElementById("room-number-info").innerHTML = msg.roomId
                break;
            case -msgId.joinRoom:
                document.getElementById("room-number-info").innerHTML = msg.roomId
                break;
            case -msgId.chessboardWalk:
                break;
            default:
                break;
        }
    } catch (error) {
        console.log(error)
    }
}

function createWs() {
    return new WebSocket("ws://" + document.location.host + "/ws")
}

window.onload = function(){
    if (window["WebSocket"]) {
        let retryTimes = 5;
        conn = createWs();
        conn.onclose = function (evt) {
            if (retryTimes > 0) {
                setTimeout(function(){
                    conn = createWs()
                },3000)
                retryTimes--
            }
        };
        conn.onmessage = function (evt) {
            // var messages = evt.data.split('\n');
            console.log(evt)
            handle(evt)
        };
        conn.onopen = function (){
          connect()
          console.log("connected")
          listRoom()
        };
    }

    document.getElementById("go-board").onclick = function(e) {
        if (e.target.id.startsWith("go-")){
            let arr = e.target.id.split("-");
            let x = arr[1];
            let y = arr[2];
            conn.send(JSON.stringify({
                msgId: msgId.chessboardWalk,
                body: {
                    "x":parseInt(x),
                    "y":parseInt(y),
                    "roomId": parseInt(document.getElementById("room-number-info").innerHTML),
                }
            }));
        }
    }
}