var user = {
    session: null,

    login: function() {
        var furl = $.url('?furl');
        if (furl == null) furl = "index.html";
        var name = $('#id-input-name').val();
        var pwd = $('#id-input-pwd').val();
        utils.doPost("/api/v1/user/login", {name:name, pwd:pwd}, function(data) {
            window.location.href = furl;
//            $('#id-errmsg').html("Hi " + data.nickname + ", welcome TODO!");
        }, function(code, errmsg) {
            $('#id-errmsg').html("HTTP " + code + ": " + errmsg);
        });
    },

    checklogin:function(url,successFuc){
        utils.doGet("/api/v1/user/login",function(data){
            if(typeof(successFuc) != "undefined") successFuc(data);
            console.log(data);
        },function(code,errmsg){
            console.log("HTTP " + code + ": " + errmsg);
            window.location.href = url;
        });
    },

    logout: function() {
        utils.doGet("/api/v1/user/logout",function(data){
            console.log(data);
            window.location.href = "login.html";
        },function(code,errmsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },

    createuser: function() {
        utils.doPost("/api/v1/user",{name:"lq123",pwd:"123456",email:"lq@todo.com"},function(data){
            console.log(data);
        },function(code,ermsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },

    deleteuser: function() {
        utils.doDelete("/api/v1/user/34",function(data){
            console.log(data);
        },function(code,ermsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },

    getuser: function(id){
        if (typeof(id) == "undefined") id="";
        else id = "/" + id;
        utils.doGet("/api/v1/user"+id,function(data){
            console.log(data);
        },function(code,ermsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },

    updateuser: function(id){
        if (typeof(id) == "undefined") id="";
        else id = "/" + id;
        utils.doPatch("/api/v1/user"+id,function(data){
            console.log(data);
        },function(code,ermsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },





}
