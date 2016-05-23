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

    checklogin:function(url){
        utils.doGet("/api/v1/user/login",function(data){
            console.log(data);
        },function(code,errmsg){
            window.location.href = url;
        });
    },

    logout: function() {
        utils.doGet("/api/v1/user/logout",function(data){
            window.location.href = "login.html";
        },function(code,errmsg){
            console.log(errmsg);
        });
    },


}
