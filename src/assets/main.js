var todo={


    init:function(){
        todo.pagesize_update();
        $(window).resize(todo.pagesize_update);


        user.checklogin("login.html?index.html",function(data){
            var sss = "<span class='glyphicon glyphicon-user'></span>"+data.name;
            $("#usericon").html(sss);
        },function(){});


    },

    pagesize_update:function(){
        var h = $(window).height() - 80;
        var w = $(window).width() - 160;
        $("#bodybox").css("width",w);
    },
















}
