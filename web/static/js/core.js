
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
    'chessboardWalk': 4,
    'startGame': 5,
    'restartGame': 6,
    'leaveRoom': 7,
    'watchGame': 8,
    'askRegret': 9,
    'agreeRegret': 10,



    'ntfJoinRoom': 1001,
    'ntfStartGame': 1002,
    'ntfWalk': 1003,
    'ntfGameOver': 1004,
    'ntfRestartGame': 1005,
    'ntfLeaveRoom': 1006,
    'ntfWalkWatchingUser': 1007,
    'ntfAskRegret': 1008,
    'ntfAgreeRegret': 1009,
    'ntfSyncWalk': 1010,
    'ntfCommonMsg': 1011
}

let msgAck = {}

let conn = undefined //保存websocket对象
let user = {
    "myhand": hand.nilHand,
    "name": "",
    "isMaster": false,
    "lastPos": {x:-1,y:-1}
}

function resetGame() {
    user.myhand = hand.nilHand
    user.lastPos.x = -1
    user.lastPos.y = -1
    updateStatus(hand.nilHand)
    document.getElementById("go-board").innerHTML = ""
    generate_board(15,15)
}


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

function joinRoom(roomId) {
    conn.send(JSON.stringify({
        "msgId":msgId.joinRoom,
        "body": {
            "roomId": parseInt(roomId)
        }
    }))
}

function startGame() {
    conn.send(JSON.stringify({
        "msgId":msgId.startGame,
        "body": {
            "roomId": parseInt(document.getElementById("room-number-info").innerHTML)
        }
    }))
}

function restartGame() {
    conn.send(JSON.stringify({
        "msgId":msgId.restartGame,
        "body": {
            "roomId": parseInt(document.getElementById("room-number-info").innerHTML)
        }
    }))
}

function watchGame(roomId) {
    conn.send(JSON.stringify({
        "msgId":msgId.watchGame,
        "body": {
            "roomId": parseInt(roomId)
        }
    }))
}


function askRegret() {
    conn.send(JSON.stringify({
        "msgId": msgId.askRegret,
        "body":{}
    }))    
}

function agreeRegret(agree) {
    conn.send(JSON.stringify({
        "msgId": msgId.agreeRegret,
        "body":{
            "agree": (agree >= 1)
        }
    }))
}

//初始化二维数组
function initPlace(row, col) {
    place = Array(row).fill(0).map(x => Array(col).fill(0));
}

//使刚落下的棋子闪烁,用于提示
function mark(x,y){
    if (user.lastPos.x != -1  && user.lastPos != -1){
        document.getElementById(`go-${user.lastPos.x}-${user.lastPos.y}`).classList.remove("chess-spinner")
    }
    document.getElementById(`go-${x}-${y}`).classList.add("chess-spinner");
    user.lastPos.x = x, user.lastPos.y = y;
}

// 落棋
function walk(x, y, h) {
    if (h == hand.blackHand){
        document.getElementById(`go-${x}-${y}`).classList.add("b")
    } else if(h == hand.whiteHand){
        document.getElementById(`go-${x}-${y}`).classList.add("w")
    } else {
        document.getElementById(`go-${x}-${y}`).classList.remove("w", "b", "chess-spinner")
    }
}

function resetWalk(x, y) {
    document.getElementById(`go-${x}-${y}`).classList.remove("w", "b", "chess-spinner")
}

function updateStatus(h){
    let content = ""
    let style = ""
    if (h == hand.nilHand){
        document.getElementById("chess-status").innerText = "无"
        return
    } else if (user.lastPos.x == -1 && h == hand.blackHand ){ // user.lastPos.x == -1 ，代表游戏刚刚开始，还没有棋子
        style = "spinner-grow"
        content = "轮到你了"
    } else if (user.lastPos.x == -1 && h == hand.whiteHand ) { 
        style = "spinner-border"
        content = "对方思考中"
    } else if (h == user.myhand){
        style = "spinner-border"
        content = "对方思考中"
    } else {
        style = "spinner-grow"
        content = "轮到你了"
    }
    let elm = `<span>
                <span class="${style} ${style}-sm text-primary" role="status" aria-hidden="true"></span>
                <span style="font-size:0.5rem">${content}</span>
            </span>`
    document.getElementById("chess-status").innerHTML = elm
}

function getEnemyHand(h) {
    switch (h){        
        case hand.blackHand:
            return hand.whiteHand
        case hand.whiteHand:
            return hand.blackHand
        default:
            return hand.nilHand
    }
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
            case -msgId.error:
                modalSystemMessage(msg.msg)
                break;
            case -msgId.connect:
                sessionStorage.setItem("un",msg.username)
                user.name = msg.username
                document.getElementById("myname").innerText = user.name
                if (msg.roomId > 0) {
                    document.getElementById("room-number-info").innerHTML = msg.roomId
                    document.getElementById("dating").classList.add("d-none")
                    document.getElementById("room").classList.remove("d-none")
                    msg.walks.forEach(element =>{
                        walk(element.x,element.y,element.hand)
                    })
                    if (!msg.isWatcher){
                        user.myhand = msg.myhand
                        user.lastPos.x = msg.latest.x
                        user.lastPos.y = msg.latest.y
                        updateStatus(msg.latest.hand)
                    }
                    if (msg.latest.x >= 0){
                        mark(msg.latest.x,msg.latest.y)
                    }
                } else {
                    listRoom()
                }
                break;
            case -msgId.listRoom:
                let tmp = ""
                msg.data.forEach(element => {
                    element.master = element.master == "" ? "无":element.master
                    element.enemy = element.enemy == "" ? "无":element.enemy
                    isDisabled = element.isFull == true ? "disabled": ""
                    funcName = element.isFull == true ? "watchGame": "joinRoom"
                    info = element.isFull == true ? "观战": "加入"
                    tmp += `<tr>
                    <th scope="row">${element.roomId}</th>
                    <td>${element.master}</td>
                    <td>${element.enemy}</td>
                    <td><button type="button" class="btn btn-sm btn-primary" onclick="${funcName}(${element.roomId})">${info}</button></td>
                  </tr>`
                });
                document.getElementById("dating-data").innerHTML = tmp
                break;
            case -msgId.createRoom:
                user.isMaster = true
                generate_board(15,15)
                document.getElementById("dating").classList.add("d-none")
                document.getElementById("room").classList.remove("d-none")
                document.getElementById("room-number-info").innerHTML = msg.roomId
                break;
            case -msgId.joinRoom:
                generate_board(15,15)
                document.getElementById("dating").classList.add("d-none")
                document.getElementById("room").classList.remove("d-none")
                document.getElementById("room-number-info").innerHTML = msg.roomId
                break;
            case -msgId.chessboardWalk:
                break;
            case -msgId.startGame:
                break;
            case -msgId.restartGame:
                modalSystemMessage("游戏已重开")
                break;
            case -msgId.leaveRoom:
                document.getElementById("dating").classList.remove("d-none")
                document.getElementById("room").classList.add("d-none")
                listRoom()
                resetGame()
                break;
            case -msgId.watchGame:
                document.getElementById("dating").classList.add("d-none")
                document.getElementById("room").classList.remove("d-none")
                document.getElementById("room-number-info").innerHTML = msg.roomId
                break;
            case -msgId.askRegret:
                break;
            case -msgId.agreeRegret:
                break;

            case msgId.ntfJoinRoom:
                if (msg.username == user.name) {
                    modalSystemMessage("您成为房主了")
                }else {
                    modalStartGame(msg.username)
                }
                break;
            case msgId.ntfStartGame:
                user.myhand = msg.hand
                if (!user.isMaster){modalSystemMessage("游戏开始了")}
                updateStatus(msg.hand)
                break;
            case msgId.ntfWalk:
                walk(msg.x,msg.y,msg.hand)
                mark(msg.x,msg.y)
                updateStatus(msg.hand)
                break;
            case msgId.ntfGameOver:
                modalSystemMessage(msg.msg)
                break;
            case msgId.ntfRestartGame:
                resetGame()
                user.myhand = msg.hand
                if (!user.isMaster){modalSystemMessage("房主重开了游戏")}
                updateStatus(msg.hand)
                break;
            case msgId.ntfLeaveRoom:
                if (!user.isMaster){
                    modalSystemMessage("对方离开了,您已成为房主")
                    user.isMaster = true
                }else {
                    modalSystemMessage("对方离开了该房间")
                }
                resetGame()
                break;
            case msgId.ntfWalkWatchingUser:
                msg.walks.forEach(element =>{
                    walk(element.x,element.y,element.hand)
                })
                mark(msg.latest.x,msg.latest.y)
                break;
            case msgId.ntfAskRegret:
                modalAskRegret()
                break;
            case msgId.ntfAgreeRegret:
                if (msg.agree){
                    modalSystemMessage("对方同意了您的悔棋")
                } else{
                    modalSystemMessage("对方拒绝悔棋了")
                }
                break
            case msgId.ntfSyncWalk:
                generate_board(15,15)
                msg.walks.forEach(element =>{
                    walk(element.x,element.y,element.hand)
                })
                mark(msg.latest.x,msg.latest.y)
                break;
            case msgId.ntfCommonMsg:
                modalSystemMessage(msg.msg)
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
            conn = createWs()
        };
        conn.onmessage = function (evt) {
            // var messages = evt.data.split('\n');
            console.log(evt)
            handle(evt)
        };
        conn.onopen = function (){
          connect()
          console.log("connected")
        };
    }
    generate_board(15,15)
    document.getElementById("go-board").onclick = function(e) {
        if (e.target.id.startsWith("go-")){
            let arr = e.target.id.split("-");
            if (arr.length !== 3) {return}
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

function modalStartGame(username) {
    let modalEl = document.getElementById('modalStartGame')
    let body = modalEl.getElementsByClassName("modal-body")[0]
    if (body != undefined && username =="") {
        body.innerHTML = "请开始游戏"
    }
    body.textContent = `玩家:${username}加入游戏,可以开始游戏了`
    let modal = new bootstrap.Modal(modalEl,{keyboard: false,backdrop:"static"})
    modal.show()
}

function modalSystemMessage(message) {
    let modalEl = document.getElementById('modalSystemMessage')
    let body = modalEl.getElementsByClassName("modal-body")[0]
    if (body == undefined || message == "" || message == undefined) {
        return
    }
    body.textContent = message
    let modal = new bootstrap.Modal(modalEl,{keyboard: false,backdrop:"static"})
    modal.toggle()
}

function modalAskRegret() {
    let modalEl = document.getElementById('modalAskRegret')
    let modal = new bootstrap.Modal(modalEl,{keyboard: false,backdrop:"static"})
    modal.show()
}

function btnGameStart() {
    startGame()
}

function btnLeaveRoom() {
    conn.send(JSON.stringify({
        "msgId":msgId.leaveRoom,
        "body": {
            "roomId": parseInt(document.getElementById("room-number-info").innerHTML)
        }
    }))
}

function btnGameRestart() {
    restartGame()
}

function btnRegret() {
    askRegret()
}