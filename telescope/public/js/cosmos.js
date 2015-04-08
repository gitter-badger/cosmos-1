(function() {
    // Tab click event listener
    var tabs = document.querySelector('paper-tabs');
    var pages = document.querySelector('core-pages');

    tabs.addEventListener('core-select', function() {
        pages.selected = tabs.selected;

        switch (pages.selected) {
            case 0:
	            var pageCosmos = document.querySelector('page-cosmos');
                Cosmos.request.getPlanets(function(json) {
                	pageCosmos.planets = json;
                	console.log(pageCosmos.planets);
                });
                Cosmos.request.getContainers(null, function(json) {
                    // Succeed
                    pageCosmos.containers = json;
                }, function() {
                    // Failed
                    alert('request failed');
                });
                break;

            case 1:
                break;
            case 2:
                Cosmos.request.getContainers(null, function(json) {
                    // Succeed
                    var pageContainer = document.querySelector('page-container');
                    pageContainer.containers = json;
                }, function() {
                    // Failed                    
                    alert('request failed');
                })
                break;
        }
    });

    // Route configuration
    $(function() {
        Route.match('/', function() {
            window.location.replace("/cosmos");
        });
        Route.match('/cosmos', function() {
            document.querySelector('#tabs').selected = 0;
        });
        Route.match('/planets', function() {
            document.querySelector('#tabs').selected = 1;
        });
        Route.match('/containers', function() {
            document.querySelector('#tabs').selected = 2;
        });
        Route.regist();
    });
}());
