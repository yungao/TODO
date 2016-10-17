angular.module('app')
    .controller('todoCtrl', ['$scope', "$http", 'http', 'storage', 'logger', function($scope, $http, http, storage, logger) {
        // check login
        if (!storage.checkLogin("login.html?index.html")) {
            return;
        }

        $scope.user = {};
        $scope.users = [];
        $scope.todos = [];
        $scope.newTodo = {}; // for store add todo infos
        $scope.curTodo = null; // for store current processing todo infos
        $scope.curIndex = 0; // for store current processing todo index
        $scope.update = null; // for store current update todo infos

        // yyyy-MM-dd
        function formatDate (unixtm) {
            function pad(n){return (n<10 ? '0'+n : n);}
            var d = new Date(unixtm * 1000);
            return d.getFullYear() + '-' + pad(d.getMonth()+1) + '-' + pad(d.getDate());
        }
        // yyyy-MM-dd HH:mm:ss
        function formatDateTime (time) {
            function pad(n){return (n<10 ? '0'+n : n);}
            var d = new Date(time * 1000);
            return d.getFullYear() + '-' + pad(d.getMonth()+1)+'-'
                + pad(d.getDate()) + ' ' + pad(d.getHours()) + ':'
                + pad(d.getMinutes());
        }

        function parseDate (obj) {
            if (obj == null) return 0;
            var time = obj.replace(new RegExp("-", "gm"), "/");
            var date = new Date(time);
            return date.getTime();
        }

        function dateToLimit(flimit) {
            return parseDate(flimit) / 1000 + 86399;  // 86399 = 24*60*60 - 1
        }

        function initICheck() {
            $('input.icheck-square-blue').iCheck({
                checkboxClass: 'icheckbox_square-blue',
                radioClass: 'iradio_square-blue',
            });
            $('input.icheck-flat-red').iCheck({
                checkboxClass: 'icheckbox_flat-red',
                radioClass: 'iradio_flat-red'
            });
            $('input.icheck-flat-pink').iCheck({
                checkboxClass: 'icheckbox_flat-pink',
                radioClass: 'iradio_flat-pink'
            });
        }

        function updateUserInfo() {
            if (typeof $scope.user.profile == 'undefined' || $scope.user.profile.length == 0) {
                $scope.user.profile = '/assets/img/ic_profile_default.png';
                // $scope.$apply();
            }
        }

        $scope.formatTime = function (time) {
            var now = (new Date().getTime()) / 1000;
            var interval = now - time;
            //计算出相差天数
            var days = Math.floor(interval/ (24*3600));
            if (days <= 2) {
                if (days == 2) {
                    return '前天';
                } else if (days == 1) {
                    return '昨天';
                } else {
                    //计算出小时数
                    var h = interval % (24*3600);
                    var hours = Math.floor(h / 3600);
                    if (hours >= 1) {
                        return hours + '小时前';
                    }

                    //计算小时数后剩余的毫秒数
                    var m = h % 3600;
                    var minutes = Math.floor(m/60);

                    if (minutes >= 1) {
                        return minutes + ' 分钟前';
                    } else {
                        return '刚刚';
                    }
                }
            } else {
                return formatDateTime(time);
            }
        };

        $scope.aboutTodo = function () {
            $('#aboutTodo').modal('show');
        };

        $scope.userProfile = function(user) {
            if (user != null && user.profile.length > 0) {
                return user.profile;
            }
            return "/assets/img/ic_profile_default.png";
        };

        $scope.userNickname = function(user) {
            if (user.nickname.length > 0) {
                return user.nickname;
            }
            return "-";
        };

        $scope.userEmail = function(user) {
            if (user.email.length > 0) {
                return user.email;
            }
            return "-";
        };

        $scope.userTooltip = function(user) {
            return $scope.userNickname(user) + "&nbsp;:&nbsp;&nbsp;" + $scope.userEmail(user);
        };

        $scope.isPartner = function(user) {
            if (user.id == $scope.user.id) {
                return true;
            }

            if ($scope.curTodo != null && $scope.curTodo.hasOwnProperty('partners')) {
                var partners = $scope.curTodo.partners;
                for (i = 0; i < partners.length; i++) {
                    if (partners[i].pid == user.id) {
                        return true;
                    }
                }
            }

            return false;
        };

        $scope.userInfo = function(uid) {
            for (i = 0; i < $scope.users.length; i++) {
                if ($scope.users[i].id == uid) {
                    return $scope.users[i];
                }
            }

            return null;
        }

        $scope.allUsers = function() {
            http.get("/user", {}, function(data) {
                $scope.users = data;
                // logger.info("Users: " + JSON.stringify(data));
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.createUser = function() {
            var user = {name:'papa', pwd:'123456'};
            http.post("/user", user, function(data) {
                logger.info(data);
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.allTodos = function () {
            http.get("/todo", {}, function(data) {
                $scope.todos = data;
                // logger.info("All TODO: " + JSON.stringify(data));
                $scope.$apply();
                if ($scope.todos.length > 0) {
                    $scope.todoDetails($scope.todos[0], 0);
                }
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.todoDetails = function (todo, index) {
            http.get("/todo/" + todo.id, {}, function(data) {
                $scope.curTodo = data;
                $scope.curIndex = index;
                {
                    // $scope.todos[index] = data; // can not replace ??
                    $scope.todos[index].limit = data.limit;
                    $scope.todos[index].priority = data.priority;
                }
                $scope.curTodo.flimit = formatDate($scope.curTodo.limit);
                logger.info("Cur TODO: " + JSON.stringify(data));
                $scope.$apply();
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.addTodo = function () {
            $('#addTodo').modal('show');
			$('#addTodo #limit').datetimepicker({
				language: 'zh-CN',
				weekStart: 1,
				todayBtn: 1,
				autoclose: 1,
				todayHighlight: 1,
				startView: 2,
				//forceParse: 0,
				showMeridian: 0,
				timepicker : false,
				minView: 'month',
				startDate: new Date(new Date().getTime() + 24*60*60*1000),
			});
            initICheck();

            // $scope.newTodo = {};
            $scope.newTodo.flimit = formatDate(new Date().getTime() / 1000 + 3600 * 24 * 2);
            $scope.newTodo.type = 0;
        };

        $scope.saveNewTodo = function() {
            $('#addTodo').modal('hide');

            var seconds = dateToLimit($scope.newTodo.flimit);
            $scope.newTodo.limit = seconds;
            delete $scope.newTodo['flimit'];

            http.post("/todo", $scope.newTodo, function(data) {
                $scope.allTodos();
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.updateTodo = function() {
            http.patch("/todo/" + $scope.curTodo.id, {action:7, content:$scope.update.content}, function(data) {
                $scope.update.content = null;
                $scope.todoDetails($scope.curTodo, $scope.curIndex);
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.modifyTodo = function(action) {
            logger.error("action: " + action);
            var param;
            if (action == 6) {
                param = {action:action, priority: $scope.curTodo.priority};
            } else if (action == 5) {
                $scope.curTodo.limit = dateToLimit($scope.curTodo.flimit);
                param = {action:action, limit: $scope.curTodo.limit};
            }

            http.patch("/todo/" + $scope.curTodo.id, param, function(data) {
                $scope.curTodo = data;
                $scope.todoDetails($scope.curTodo, $scope.curIndex);
            }, function(code, errmsg) {
                logger.error(code + ": " + errmsg);
            });
        };

        $scope.addPartner = function (user) {
            if (user.id != $scope.user.id) {
                var partner = {todoid:$scope.curTodo.id, pid:user.id, duty:1};
                http.post("/todo/partner", partner, function(data) {
                    logger.info(data);
                    $scope.todoDetails($scope.curTodo, $scope.curIndex);
                }, function(code, errmsg) {
                    logger.error(code + ": " + errmsg);
                });
            } else {
                logger.error("can not modify Todo creator");
            }
        };

        $scope.processContent = function (process, pos) {
            var content;
            switch (process.action) {
                // -1： 撤销完成
                // 0： 完成
                // 1： 创建
                case 1:
                    if (pos == 0) {
                        content = '创建了新任务';
                    }
                    break;
                // 2： 添加参与者
                case 2:
                    if (pos == 0) {
                        content = '邀请了用户：';
                    } else {
                        content = $scope.userInfo(process.content).name;
                    }
                    break;
                // 3： 指派执行者
                // 4： 添加TAG
                // 5:  修改结束时间
                case 5:
                    if (pos == 0) {
                        content = '更新截止时间为：';
                    } else {
                        content = formatDate(process.content);
                    }
                    break;
                // 6:  修改紧急度
                case 6:
                    if (pos == 0) {
                        content = '更新优先级为：';
                    } else {
                        content = (process.content == '1' ? '紧急任务' : '普通任务');
                    }
                    break;

                default:
                    if (pos == 0) {
                        content = process.content;
                    }
                    break;
            }

            return content;
        };

        $scope.processActionIcon = function (process) {
            var icon;

            switch (process.action) {
                case -1: // -1： 撤销完成
                    icon = 'glyphicon glyphicon-retweet';
                    break;
                case 0: // 0： 完成
                    icon = 'glyphicon glyphicon-ok';
                    break;
                case 1: // 1： 创建
                    icon = 'glyphicon glyphicon-plus';
                    break;
                case 2: // 2： 添加参与者
                    icon =  'glyphicon glyphicon-plus-sign'; 
                    break;
                case 3: // 3： 指派执行者
                    icon = 'glyphicon glyphicon-hand-right';
                    break;
                case 4: // 4： 添加TAG
                    icon = 'glyphicon glyphicon-tag';
                    break;
                case 5: // 5:  修改结束时间
                    icon = 'glyphicon glyphicon-calendar';
                    break;
                case 6: // 6:  修改紧急度
                    icon = 'glyphicon glyphicon-flash';
                    break;
                // 7： 更新
                default:
                    icon = 'glyphicon glyphicon-comment';
                    break;
            }

            return icon;
        };

        $scope.selectUsers = function () {
            $("#selectUser").modal('show');
        };

        $scope.selectTags = function () {
            $("#selectTag").modal('show');
        };

        $scope.user = storage.getUser();
        $scope.allUsers();
        updateUserInfo();
        $scope.allTodos();
        // init datetimepicker
        $('#curTodo #limit').datetimepicker({
            language: 'zh-CN',
            weekStart: 1,
            todayBtn: 1,
            autoclose: 1,
            todayHighlight: 1,
            startView: 2,
            //forceParse: 0,
            showMeridian: 0,
            timepicker : false,
            minView: 'month',
            startDate: new Date(new Date().getTime() + 24*60*60*1000),
        });

    }]);


