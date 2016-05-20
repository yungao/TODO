var utils = {
    doPost: function(url, data, successFunc, errorFunc) {
        $.post(url, data)
            .success(function(response){
                successFunc(response);
            }).error(function(err) {
                errorFunc(err.status, err.responseText);
            });
    },

}
