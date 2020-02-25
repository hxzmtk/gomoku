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
    if (playing == hand.nilHand) {
        $("#chess-status").empty();
        $("#chess-status").append("无")
        return
    }
    let content = "";
    let style = "";
    if (playing != player) {
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

    $(".go-board i").removeClass("w b chess-spinner");
    initPlace(15,15);
    updateIdentity(hand.nilHand);
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
            }
        }
    });
}

//消息提示
function BootboxAlert(msg){
    bootbox.alert(msg);
    window.setTimeout(function(){
        bootbox.hideAll();
    },3000);
}

//处理RoomMsg消息
function handleRoomMsg(msg){
    if (msg.status = true) {
        switch (msg['content']['action']) {
            case roomAction.create:
                $("#room-number-info").html(msg['content']['room_number']);
                is_master = true;
                break;
            case roomAction.join:
                console.log(msg['content'],$("#myName").html());
                if (is_master) {
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
                ResetAll();
                $(".container").addClass("d-none");
                $('#dialog').modal('show');
                break;
            default:
                break;
        }
    }
    if (msg.msg != ""){
        alertMsg(msg.msg);
    }
};

//处理下棋消息
function handleChessWalkMsg(msg){
    if (msg.status == true) {
        Setup(msg['content'].x, msg['content'].y,msg['content'].is_black == true?1:2);
        updateStatus(msg['content'].is_black == true?hand.blackHand:hand.whiteHand);
        remain(msg['content'].x, msg['content'].y);
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

    ws = new WebSocket("ws://"+ document.location.host + "/v1/ws");
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
                $("#myname").html(dic.content.name);
                $("#myName").html(dic.content.name);
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
    $("#choice-enemy-2").on("click", function(e){
        if ($(this).is(":checked")){
            $(".choice-action").addClass("d-none");
            $(".choice-level").removeClass("d-none");
            $(".group-choice-room").addClass("d-none");
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