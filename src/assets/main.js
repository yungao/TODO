var todo={


    init:function(){
        todo.pagesize_update();
        $(window).resize(todo.pagesize_update);

        $("#usericon").click(function(){
            user.checklogin("login.html?index.html",function(data){
                var userinfotpl = $("#user-info-tpl").html();
                var userinfo = Mustache.compile(userinfotpl);
                var uinfomation = userinfo(data);
                $("#user-modal #user-info").empty().append(uinfomation);
            });
        });

        $("#user-modal #save-btn").click(function(){
            var options = {};
            $("#user-modal input").each(function(i,n){
                if(!$(this).attr("readonly")&&n.value)  options[n.getAttribute("name")]=n.value;
            });
            //user.updateuser(JSON.stringify(options));
            user.updateuser(options);
        });




    },

    pagesize_update:function(){
        var h = $(window).height();
        var w = $(window).width();

        if($("#user-modal").height()>$(window).height()) $("#user-modal").css("height",$(window).height());
            else    $("#user-modal").css("height",620);
        if($("#user-modal").width()>$(window).width()) $("#user-modal").css("width",$(window).width());
            else    $("#user-modal").css("width",480);
        h = $(window).height()-$("#user-modal").height();
        w = $(window).width()-$("#user-modal").width();
        $("#user-modal").css("top",h/2);
        $("#user-modal").css("left",w/2);

    },
















}
