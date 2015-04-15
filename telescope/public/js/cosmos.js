(function() {
    // Route configuration    
    function PolymerReady(callback) {
        if (window.Polymer) {
            console.log(window.Polymer);
            callback();
        } else {
            setTimeout(function() {
                PolymerReady(callback);
            }, 50);
        }
    }

    Route.match('/', function() {
        PolymerReady(function() {
            console.log(Polymer.dom);
            var content = Polymer.dom(document.querySelector('cosmos-content'));
            console.log(content);
            content.page = '';
        });
    });
    Route.match('/planets', function() {
        //            document.querySelector('cosmos-content').page = 'planet';
    });
    Route.match('/containers', function() {
        //            document.querySelector('cosmos-content').page = 'container';
    });
    Route.defaultRoute = function() {
        window.location.replace("/");
    };
    Route.regist();

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
