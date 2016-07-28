angular.module('app')
    .directive('tooltip', function () {
        return {
            restrict: 'A',
            link: function (scope, element, attributes) {
                $(element).tooltip({
                    html: 'true'
                });
            }
        };
});
