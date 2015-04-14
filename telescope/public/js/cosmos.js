(function() {
    // Route configuration
    $(function() {
        Route.match('/', function() {
        });
        Route.match('/planets', function() {
        });
        Route.match('/containers', function() {
        });
        Route.defaultRoute = function() {
            window.location.replace("/");
        };
        Route.regist();
    });


    window.onpopstate = history.onpushstate = function(e) { 
        console.log(e);
    };

}());
