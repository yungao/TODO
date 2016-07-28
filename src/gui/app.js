var app = angular
    .module('app', [])
    .service('storage', storage)
    .factory('logger', logger)
    .factory('http', http);

var urlbase = '/api/v1';

function storage(logger, http) {
    this.user = {};

    this.setUser = function(user) {
        this.user = user;
        // logger.log("Set user: " + JSON.stringify(this.user));
    };

    this.getUser = function() {
        // logger.log("Get user: " + JSON.stringify(this.user));
        return this.user;
    };

    this.checkLogin = function(url) {
        var resp = http.getSync("/user/login", {});
        if (resp.status == '401') {
            logger.error("Do not login, " + resp.status + "- " + resp.statusText); 
            window.location.href = url;
            return false;
        } else {
            this.user = eval("(" + resp.responseText + ")");
            return true;
        }
    };
};

function http() {
    var http = {
        get: function (url, data, successFunc, errorFunc) {
            url = urlbase + url;
            $.get(url, data).success(function(response) {
                successFunc(response);
            }).error(function(err) {
                errorFunc(err.status, err.responseText);
            });
        },

        getSync: function(url, data, successFunc, errorFunc) {
            url = urlbase + url;
            return $.ajax({ 
                type : "GET", 
                url : url, 
                data : data, 
                async : false, 
                contentType: "application/json",
            }); 
        },

        post: function (url, data, successFunc, errorFunc) {
            url = urlbase + url;
            $.post(url, data).success(function(response) {
                successFunc(response);
            }).error(function(err) {
                errorFunc(err.status, err.responseText);
            });
        },

        patch: function (url, data, successFunc, errorFunc) {
            url = urlbase + url;
            $.ajax({
                type:   "PATCH",
                url:    url,
                data:   data,
                success:function (response) {
                    successFunc(response);
                },
                error:  function (err) {
                    errorFunc(err.status, err.responseText);
                }
            });
        },

        postJson: function (url, data, successFunc, errorFunc) {
            url = urlbase + url;
            $.ajax({
                type:           "POST",
                url:            url,
                contentType:    "application/json",
                dataType:       "json",
                data:           data,
                success:        function (response) {
                    successFunc(response);
                },
                error:          function (err) {
                    errorFunc(err.status, err.responseText);
                }
            });
        },
    };

    return http;
}

function logger() {
    var TAGS = {
        log:    'LOG:\t',
        info:   'INFO:\t',
        debug:  'DEBUG:\t',
        warn:   'WARNING:\t',
        error:  'ERROR:\t',
    };

    var logger = {
        log: function (msg) {
            console.log(TAGS['log'], msg);
        },

        info: function(msg) { 
            console.log(TAGS['info'], msg);
        },

        debug: function(msg) {
            console.debug(TAGS['debug'], msg);
        },

        warn: function(msg) {
            console.warn(TAGS['warn'], msg);
        },

        error: function(msg) {
            console.error(TAGS['error'], msg);
        },
    };

    return logger;
};
