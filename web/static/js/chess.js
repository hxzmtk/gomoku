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
                ctx.moveTo(40, 40 + i * 40);
                ctx.lineTo(15 * 40, 40 + i * 40);
    
                //纵线
                ctx.moveTo(40 + i * 40, 40);
                ctx.lineTo(40 + i * 40, 15 * 40);

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
        if (who == 1) { //黑手
        } else if (who == 2) { //白手

        }
    };

    draw();
});