angular.module('app')
    .directive('popover', function () {
        return {
            restrict: 'A',
            link: function (scope, element, attributes) {
                $(element).popover({
                    html: 'true'
                });
            }
        };
});
