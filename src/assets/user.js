var user = {
    user: [],

    checkLogin: function (url) {
        var name = $.cookie('username');
        var id = $.cookie('userid');
        if (typeof(name) == 'undefined' || name == "" || typeof(id) == 'undefined' || id == "") {
            window.location.href = url;
            return false;
        }

        user["name"] = name;
        user["id"] = id;
        return true;
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

    updateuser: function(options){
        utils.doPatch("/api/v1/user",options,function(data){
            console.log(data);
        },function(code,ermsg){
            console.log("HTTP " + code + ": " + errmsg);
        });
    },





}
