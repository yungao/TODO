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
            success:function(data){
                successFunc(data);
            },
            error:function(data){
                errorFunc(data);
            },
        });
    },

}
