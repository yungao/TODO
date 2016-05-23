var utils = {
    doPost: function(url, data, successFunc, errorFunc) {
        $.post(url, data)
            .success(function(response){
                successFunc(response);
            }).error(function(err) {
                errorFunc(err.status, err.responseText);
            });
    },
    
    doGet: function(url,successFunc,errorFunc) {
        $.ajax({
            url:url,
            type:"get",
            dataType:"json",
            success:function(data){
                successFunc(data);
            },
            error:function(data){
                errorFunc(data.status,data.responseText);
            },
        });
    },

    doDelete: function(url,successFunc,errorFunc) {
         $.ajax({
            url:url,
            type:"delete",
            dataType:"json",
            success:function(data){
                successFunc(data);
            },
            error:function(err){
                errorFunc(err.status,err.responseText);
            },
        });
    },

    doPatch: function(url,successFunc,errorFunc) {
         $.ajax({
            url:url,
            type:"patch",
            dataType:"json",
            success:function(data){
                successFunc(data);
            },
            error:function(err){
                errorFunc(err.status,err.responseText);
            },
        });
    },

}
