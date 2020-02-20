function generate_board(row, col){

    for (let i = 0; i < 15; i++) {
        let tmp = ""
        for (let j = 0; j < 15; j++) {
            tmp += `<i class="i-nomal" id="go-${i}-${j}"></i>`
            $(".go-board").append(`<i id="go-${i}-${j}"></i>`)
        }
        // $(".go-board").append(`<div>${tmp}</div>`)
        $(".go-board").append(`<br>`)
    }
}

function Setup(x, y, color) {
    if (color == 1){
        $(`#go-${x}-${y}`).addClass("b");
    } else if(color==2){
        $(`#go-${x}-${y}`).addClass("w");
    }
}

function flushRoom(arr) {
    $("#room-list").empty();
    if (arr.length == 0) {
        $("#room-list").append(`<option value="0">空房间</option`)
    }
    for (let i = 0; i < arr.length; i++) {
        let msg = "可加入"
        if (arr[i]['is_full']){
            msg = "已满"
        }
        $("#room-list").append(`<option value="${arr[i]['room_number']}">房间号:${arr[i]['room_number']} ${msg}</option`)      
    }
}


$(document).ready(function(){

    $(window).on('load',function(){
        generate_board(15,15);
        $('#dialog').modal('show');
    });

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

    const ws = new WebSocket("ws://"+ document.location.host + "/v1/ws");
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
                    if (dic['content']['action'] == 'create') {
                        $("#room-number-info").html(dic['content']['room_number']);
                    }
                    else if (dic['content']['action'] == 'join'){
                        $("#room-number-info").html(dic['content']['room_number']);
                        $("#user-info").html(dic['content'].is_black == true?"先手":"后手");
                    }
                    $(".chess").removeClass("invisible");  
                }

                if (dic.msg != ""){
                    $("#toast-1 .toast-body").html(dic.msg);
                    $('.toast').toast('show');
                }
                
                break;
            case 1:
                if (dic.status == true) {
                    Setup(dic['content'].x, dic['content'].y,dic['content'].is_black == true?1:2);
                }
                if (dic.msg != ""){
                    $("#toast-1 .toast-body").html(dic.msg);
                    $('.toast').toast('show');
                }
                break;
            case 2:
                console.log(dic);
                flushRoom(dic['content'])
                break;
        }
    };

    $('.toast').on('hidden.bs.toast', function () {
        // do something...
    });

    //dialog
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
});