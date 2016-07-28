angular.module('app')
    .controller('loginCtrl', ['$scope', 'http', 'storage', 'logger', function($scope, http, storage, logger) {
        function login(user, pwd) {
            var furl = $.url('?furl');
            if (furl == null) {
                furl = "index.html";
            }

            http.post("/user/login", {name:user, pwd:pwd}, function(data) {
                storage.setUser(data);
                window.location.href = furl;
            }, function(code, errmsg) {
                if (code == 40001) {
                    $('#errorLogin').html("密码错误！").show();
                } else {
                    $('#errorLogin').html("用户不存在！").show();
                }
            });
        }

        $scope.login = function() {
            login($scope.user, $scope.pwd);
        };
    }]);
