//全局变量
let hand = {
    "nilHand": 0,
    "blackHand": 1,
    "whiteHand": 2
}

let player = hand.nilHand;  //记录 “我”是 '黑子'还是'白子';
let playing = hand.nilHand; //记录当前该"谁"落子了
let place = undefined //二维数组；存放棋子
let last_pos = {x:-1,y:-1} //存放上一个"落子"的位置

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
            $("#user-info").html('先手');
            break;
        case hand.whiteHand:
            $("#user-info").html('后手');
            break;
        default:
            $("#user-info").html('无');
            break;
    }
};

//更新状态
function updateStatus(who){
    playing = who;
    if (playing == hand.nilHand) {
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

    $(".go-board i").removeClass("w b");
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
                    "m_type": 0,
                    "content": {
                        "action":"start"
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
                    "m_type": 0,
                    "content": {
                        "action":"leave",
                        "room_number":parseInt($("#room-number-info").html())
                    }
                }))
            }
        }
    });
}

$(document).ready(function(){

    $(".go-board").on("click", function(e){
        if (e.target.id.startsWith("go-")){
            let arr = e.target.id.split("-");
            let x = arr[1];
            let y = arr[2];

            let msg = {
                "m_type": 1,
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
            "m_type": 0,
            "content": {
                "action":"create"
            }
        }
        ws.send(JSON.stringify(msg));
        $('#dialog').modal('hide');
        $(".container").removeClass("d-none");
    });

    $("#room-join").on("click", function(e){
        let msg = {
            "m_type": 0,
            "content": {
                "action":"join",
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
    };

    ws.onclose = function(){
        console.log("DISCONNECT");
    };

    ws.onmessage = function(event){
        console.log(event.data);
        let dic = JSON.parse(event.data);
        switch (dic['m_type']) {
            case 0:
                console.log(dic);
                if (dic.status == true) {
                    let action = dic['content']['action'];
                    if (action == 'create') {
                        $("#room-number-info").html(dic['content']['room_number']);
                    } else if (action == 'join'){
                        $("#room-number-info").html(dic['content']['room_number']);
                        ConfirmGameStart();
                        // $("#user-info").html(dic['content'].is_black == true?"先手":"后手");
                        updateIdentity(dic['content'].is_black == true?hand.blackHand:hand.whiteHand)

                    } else if (action == 'leave'){
                        ResetAll();
                        $(".container").addClass("d-none");
                        $('#dialog').modal('show');

                    }  else if (action == "start") {
                        $("#room-number-info").html(dic['content']['room_number']);
                        updateIdentity(dic['content'].is_black == true?hand.blackHand:hand.whiteHand);
                    }
                }

                if (dic.msg != ""){
                    alertMsg(dic.msg);
                }
                
                break;
            case 1:
                if (dic.status == true) {
                    Setup(dic['content'].x, dic['content'].y,dic['content'].is_black == true?1:2);
                    updateStatus(dic['content'].is_black == true?hand.blackHand:hand.whiteHand);
                    remain(dic['content'].x, dic['content'].y);
                }
                if (dic.msg != ""){
                    alertMsg(dic.msg);
                }
                break;
            case 2:
                console.log(dic);
                flushRoom(dic['content'])
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
                "m_type": 2,
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