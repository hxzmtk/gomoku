//全局变量
let hand = {
    "nilHand": 0,
    "blackHand": 1,
    "whiteHand": 2
}

let msgType = {
     "clientInfoMsg": 0,
    "roomMsg": 1,
    "chessWalkMsg": 2,
    "roomList": 3,
}
let roomAction = {
    "create": 0,
    "join": 1,
    "start": 2,
    "leave": 3,
    "restart": 4,
    "reset": 5,
    "watch": 6,
    "watchChessWalk": 7,
    "roomRegret": 8,
    "roomRegretAgree": 9,
    "roomRegretReject": 10,
}

let player = hand.nilHand;  //记录 “我”是 '黑子'还是'白子';
let playing = hand.nilHand; //记录当前该"谁"落子了
let place = undefined //二维数组；存放棋子
let last_pos = {x:-1,y:-1} //存放上一个"落子"的位置
let is_master = false; //是否是房主

let ws = undefined //保存websocket对象

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

    for (let i = 0; i < row; i++) {
        let tmp = ""
        for (let j = 0; j < col; j++) {
            tmp += `<i class="i-nomal" id="go-${i}-${j}"></i>`
            $(".go-board").append(`<i id="go-${i}-${j}"></i>`)
        }
        // $(".go-board").append(`<div>${tmp}</div>`)
        $(".go-board").append(`<br>`)
    }

    initPlace(row,col);
}

// 落棋
function Setup(x, y, color) {
    if (color == 1){
        $(`#go-${x}-${y}`).addClass("b");
    } else if(color==2){
        $(`#go-${x}-${y}`).addClass("w");
    } else {
        $(`#go-${x}-${y}`).removeClass("w b chess-spinner");
    }
}


// 获取房间列表
function flushRoom(arr) {
    $("#room-list").empty();
    if (arr.length == 0) {
        $("#room-list").append(`<option value="0">空房间,请创建房间</option`)
    }
    for (let i = 0; i < arr.length; i++) {
        let msg = "可加入"
        if (arr[i]['is_full']){
            msg = "已满"
        }
        $("#room-list").append(`<option value="${arr[i]['room_number']}">房间号:${arr[i]['room_number']} ${msg}</option`)      
    }
}


//提示消息
function alertMsg(msg) {
    let elm = $(".alert-msg");
    elm.empty();
    elm.html(`<div class="col d-flex justify-content-center">${msg}</div>`);
    elm.fadeTo(2000, 500).slideUp(500, function(){
        $(".alert").slideUp(500);
    });
}

//更新身份
function updateIdentity(who){
    player = who
    switch (player) {
        case hand.blackHand:
            $("#myName").addClass("btn-dark");
            $("#targetName").removeClass("btn-dark");
            break;
        case hand.whiteHand:
            $("#targetName").addClass("btn-dark");
            $("#myName").removeClass("btn-dark");
            break;
        default:
            $("#myName").removeClass("btn-dark");
            $("#targetName").removeClass("btn-dark");
            break;
    }
};

//更新状态
function updateStatus(who){
    playing = who;
    if (playing === hand.nilHand) {
        $("#chess-status").empty();
        $("#chess-status").append("无")
        return
    }
    let content = "";
    let style = "";
    if (playing !== player) {
        style = "spinner-grow";
        content = "轮到你了"
    }else {
        style = "spinner-border";
        content = "对方思考中"
    }

    let elm = `<span>
                <span class="${style} ${style}-sm text-primary" role="status" aria-hidden="true"></span>
                <span style="font-size:0.5rem">${content}</span>
            </span>`
    $("#chess-status").empty()
    $("#chess-status").append(elm);
}

//重置
function ResetAll(){
    player = hand.nilHand;
    playing = hand.nilHand;
    place = undefined;
    last_pos = {x:-1,y:-1}

    $("#targetName").html("待加入");

    $(".go-board i").removeClass("w b chess-spinner");
    initPlace(15,15);
    updateIdentity(hand.nilHand);
    updateStatus(hand.nilHand);
}

//重置棋盘渲染
function ResetGrid(){
    for (let i = 0; i < 15; i++) {
        for (let j = 0; j < 15; j++) {
            $(`#go-${i}-${j}`).removeClass("b w");
        }
    }
}

//重新开局
function Restart(){
    ResetGrid();
    $(".go-board i").removeClass("w b chess-spinner");
    initPlace(15,15);
    updateStatus(hand.nilHand);
}

//确认开始游戏
function  ConfirmGameStart(){
    bootbox.confirm({
        message: "对方已加入,请开始游戏",
        buttons: {
            confirm: {
                label: 'Yes',
                className: 'btn-success'
            },
            cancel: {
                label: 'No',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.start
                    }
                }))
            }
        }
    });
}

//确认是否离开房间
function  ConfirmLeaveHome(){
    bootbox.confirm({
        message: "真的要离开房间吗?",
        buttons: {
            confirm: {
                label: 'Yes',
                className: 'btn-success'
            },
            cancel: {
                label: 'No',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg, 
                    "content": {
                        "action":roomAction.leave,
                        "room_number":parseInt($("#room-number-info").html())
                    }
                }))
                is_master = false;
                ResetAll();
                $(".container").addClass("d-none");
                $('#dialog').modal('show');
            }
        }
    });
}

//重开
function ConfirmGameRestart(){
    bootbox.confirm({
        message: "确认重新开局吗",
        buttons: {
            confirm: {
                label: 'Yes',
                className: 'btn-success'
            },
            cancel: {
                label: 'No',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.restart,
                    }
                }))
            } else {

            }
        }
    });
}

//对方离线,提示房主重置
function ConfirmGameReset(){

    bootbox.alert({ 
        message: "对方离开房间或掉线,请重置", 
        closeButton: false,
        callback: function (result) {
            ws.send(JSON.stringify({
                "m_type": msgType.roomMsg,
                "content": {
                    "action":roomAction.reset,
                }
            }))
        }
    })
}

//观战
function ConfirmWatch(){
    bootbox.confirm({
        message: "确认要观战吗",
        buttons: {
            confirm: {
                label: 'Yes',
                className: 'btn-success'
            },
            cancel: {
                label: 'No',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.watch,
                        "room_number":parseInt($("#room-list :selected")[0].value)
                    }
                }))
            } else {
                $(".container").addClass("d-none");
                $('#dialog').modal('show');
            }
        }
    });
}

//悔棋
function  ConfirmRegret(){
    bootbox.confirm({
        message: "您真的要悔棋吗",
        buttons: {
            confirm: {
                label: 'Yes',
                className: 'btn-success'
            },
            cancel: {
                label: 'No',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.roomRegret
                    }
                }));
            }
        }
    });
}

// 对方请求悔棋
function  AgreeRegret(){
    bootbox.confirm({
        message: "对方请求悔棋",
        buttons: {
            confirm: {
                label: '同意',
                className: 'btn-success'
            },
            cancel: {
                label: '拒绝',
                className: 'btn-danger'
            }
        },
        callback: function (result) {
            if (result){
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.roomRegretAgree
                    }
                }))
            } else {
                ws.send(JSON.stringify({
                    "m_type": msgType.roomMsg,
                    "content": {
                        "action":roomAction.roomRegretReject
                    }
                }))
            }
        }
    });
}

//消息提示
function BootboxAlert(msg){
    bootbox.alert(msg);
}

//处理RoomMsg消息
function handleRoomMsg(msg){
    if (msg.status) {
        switch (msg['content']['action']) {
            case roomAction.create:
                $("#room-number-info").html(msg['content']['room_number']);
                updateStatus(hand.nilHand);
                is_master = true;
                break;
            case roomAction.join:
                $("#room-number-info").html(msg['content']['room_number']);

                console.log(msg['content'],$("#myName").html());
                if (msg['content']['is_master']) {
                    ConfirmGameStart();
                }
                $("#targetName").html(msg['content'].hasOwnProperty('name')?msg['content']['name']:"待加入");

                break;
            case roomAction.start:
                $("#room-number-info").html(msg['content']['room_number']);
                updateIdentity(msg['content'].is_black == true?hand.blackHand:hand.whiteHand);
                player = msg['content'].is_black == true?hand.blackHand:hand.whiteHand;
                $("#chess-status").empty();
                if (player == hand.blackHand){
                    $("#chess-status").append("您是先手");
                } else if (player == hand.whiteHand){
                    $("#chess-status").append("您是后手");
                }
                break;
            case roomAction.leave:
                if (msg['content'].hasOwnProperty('is_master')) {
                    if (msg['content']['is_master']){
                        is_master = true;
                    }else {
                        is_master = false;
                    }
                    ConfirmGameReset();
                     
                } else {
                    ResetAll();
                    $(".container").addClass("d-none");
                    $('#dialog').modal('show');
                }
                break;
            case roomAction.restart:
                Restart();
                break;
            case roomAction.reset:
                ResetAll();
                BootboxAlert("游戏重置成功");
                break;
            case roomAction.watch:
                $("#room-number-info").html(msg['content']['room_number']);
                $("#chess-status").empty().append("观战中")
                BootboxAlert("观战中");
                break;
            case roomAction.watchChessWalk:
                let xy = msg['content']['xy']
                console.log(xy)
                for (let i = 0;i<xy.length;i++) {
                    Setup(xy[i].x, xy[i].y,xy[i].hand);
                }
                remain(msg['content']["now_walk"].x, msg['content']["now_walk"].y)
                break;
            case roomAction.roomRegret:
                AgreeRegret()
                break;
            case roomAction.roomRegretAgree:
                $('#request-regret').modal('hide');
                BootboxAlert(msg["msg"])
                let xy1 = msg['content']['xy']
                for (let i = 0;i<xy1.length;i++) {
                    Setup(xy1[i].x, xy1[i].y,hand.nilHand);
                }
                break;
            default:
                break;
        }
    } else {
        switch (msg['content']['action']) {
            case roomAction.watch:
                ConfirmWatch()
                break;
            case roomAction.roomRegret:
                BootboxAlert(msg["msg"]);
                break;
        }
    }
    if (msg.msg != ""){
        alertMsg(msg.msg);
    }
};

//处理下棋消息
function handleChessWalkMsg(msg){
    if (msg.status === true) {
        Setup(msg['content'].x, msg['content'].y,msg['content']['is_black'] === true?1:2);
        updateStatus(msg['content']['is_black'] === true?hand.blackHand:hand.whiteHand);
        remain(msg['content'].x, msg['content'].y);
    } else  {
        BootboxAlert(msg["msg"]);
    }
    if (msg.msg != ""){
        alertMsg(msg.msg);
    }
};

$(document).ready(function(){

    $(".go-board").on("click", function(e){
        if (e.target.id.startsWith("go-")){
            let arr = e.target.id.split("-");
            let x = arr[1];
            let y = arr[2];

            let msg = {
                "m_type": msgType.chessWalkMsg,
                "content": {
                    "x":parseInt(x),
                    "y":parseInt(y),
                    "room_number": parseInt($("#room-number-info").html()),
                }
            }
            ws.send(JSON.stringify(msg));
        }
    });

    $("#room-create").on("click", function(e){
        let msg = {
            "m_type": msgType.roomMsg,
            "content": {
                "action":roomAction.create,
            }
        }
        console.log(msg);
        ws.send(JSON.stringify(msg));
        $('#dialog').modal('hide');
        $(".container").removeClass("d-none");
    });

    $("#room-join").on("click", function(e){
        let msg = {
            "m_type": msgType.roomMsg,
            "content": {
                "action":roomAction.join,
                "room_number":parseInt($("#room-list :selected")[0].value)
            }
        }
        ws.send(JSON.stringify(msg));
        $("#modal-room-join").modal('hide');
        $("#dialog").modal('hide');
        $(".container").removeClass("d-none");
    });

    ws = new WebSocket("ws://"+ document.location.host + "/ws/human");
    ws.onopen = function(){
        console.log("CONNECT");
        ws.send(JSON.stringify({
            "m_type": msgType.clientInfoMsg,
        }));
    };

    ws.onclose = function(){
        console.log("DISCONNECT");
    };

    ws.onmessage = function(event){
        console.log(event.data);
        let dic = JSON.parse(event.data);
        switch (dic['m_type']) {
            case msgType.roomMsg:
                handleRoomMsg(dic);
                break;
            case msgType.chessWalkMsg:
                handleChessWalkMsg(dic);
                break;
            case msgType.roomList:
                console.log(dic);
                flushRoom(dic['content'])
                break;
            case msgType.clientInfoMsg:
                console.log(dic);
                $("#myname").html(dic.content.name)
                break;
        }
        alertMsg(dic.msg);
    };

    $('.toast').on('hidden.bs.toast', function () {
        // do something...
    });

});


//dialog
$(document).ready(function(){
    $("#choice-enemy-1").on("click", function(e){
        if ($(this).is(":checked")){
            $(".choice-level").addClass("d-none");
            $(".choice-action").removeClass("d-none");
        }
    });

    $("#choice-action-1").on("click", function(e){
        if ($(this).is(":checked")){
            $(".group-choice-room").addClass("d-none");

            $("#room-join").addClass("d-none");
            $("#room-create").removeClass("d-none");
        }
    });
    $("#choice-action-2").on("click", function(e){
        if ($(this).is(":checked")){
            $(".choice-level").addClass("d-none");
            $(".choice-action").removeClass("d-none");
            $(".group-choice-room").removeClass("d-none");

            $("#room-create").addClass("d-none");
            $("#room-join").removeClass("d-none");

            ws.send(JSON.stringify({
                "m_type": msgType.roomList,
                "content": {
                }
            }))
        }
    });
});


$(window).on('load', function(){
    generate_board(15,15);
    $('#dialog').modal('show');
});