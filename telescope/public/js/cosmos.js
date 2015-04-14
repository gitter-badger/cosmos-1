(function() {
    // Route configuration
    $(function() {
        Route.match('/', function() {});
        Route.match('/planets', function() {});
        Route.match('/containers', function() {});
        Route.defaultRoute = function() {
            window.location.replace("/");
        };
        Route.regist();
    });

    $(window).on('pushstate', function(e, param) {
        if (param) {
            console.log(param.page);
        }
    });
    window.addEventListener('popstate', function(e) {
        if (e.state) {
            console.log(e.state.page);
        }
    });
}());
