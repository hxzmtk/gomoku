$(document).ready(function(){
    // https://developer.mozilla.org/zh-CN/docs/Web/API/Canvas_API/Tutorial/Drawing_shapes
    const canvas = document.getElementById('chess');
    function draw() {
        if (canvas.getContext){
            var ctx = canvas.getContext('2d');
            ctx.strokeStyle = '#e5f1ed';
            for (var i=0;i < 15; i++) {
                /*
                x:横坐标 y:纵坐标
                moveTo(x,y) 设置起点
                lineTo(x,y) 绘制一条从当前位置到指定x以及y位置的直线。
                */

                //横线
                ctx.moveTo(20, 20 + i * 40);
                ctx.lineTo(15 * 40 -20, 20 + i * 40);
    
                //纵线
                ctx.moveTo(20 + i * 40, 20);
                ctx.lineTo(20 + i * 40, 15 * 40 - 20);

                ctx.stroke();
            }
        }
    };

    function arc(x, y , who) {
        if (canvas.getContext){
            ctx = canvas.getContext('2d');
            ctx.beginPath();
            var radius = 20; //半径
            ctx.arc(x, y, radius,0, Math.PI * 2,true);
            ctx.closePath();
            var gradient = ctx.createRadialGradient(x,y,radius,x,y,0);
            if (who == 1) {
                gradient.addColorStop(0, '#000');
		        gradient.addColorStop(1, '#343a40');
            } else {
                gradient.addColorStop(0,'#FFF');
	            gradient.addColorStop(1, '#f9eaea');
            }
            ctx.fillStyle = gradient;
            ctx.fill();
        }
    };
    //落子
    function Setup(row,col,who) {
        if (row >= 15 || col >=15){
            return
        }
        if (who == 1) { //黑手
            arc(20 + row * 40,20 + col *40, 1);
        } else if (who == 2) { //白手
            arc(20 + row * 40,20 + col *40, 2);
        }
    };

    $("#chess").on("click", function(e){
        var x = Math.floor(e.offsetX/40);
        var y = Math.floor(e.offsetY/40);
        // Setup(Math.floor(x/40),Math.floor(y/40),1);
        // Setup(Math.floor(x/40),Math.floor(y/40),2);
        let msg = {
            "m_type": 1,
            "content": {
                "x":x,
                "y":y,
                "room_number": parseInt($("#room-number-info").html()),
            }
        }
        ws.send(JSON.stringify(msg));
    });

    $("#room-create").on("click", function(w){
        let msg = {
            "m_type": 0,
            "content": {
                "action":"create"
            }
        }
        ws.send(JSON.stringify(msg));
    });

    $("#room-join").on("click", function(w){
        let msg = {
            "m_type": 0,
            "content": {
                "action":"join",
                "room_number":parseInt($("#modal-room-number").val())
            }
        }
        ws.send(JSON.stringify(msg));
        $("#modal-room-join").modal('hide');
    });

    draw();

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
                        console.log("加入成功");
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
        }
    };

    $('.toast').on('hidden.bs.toast', function () {
        // do something...
    });
});