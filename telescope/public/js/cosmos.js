(function() {
    // Tab click event listener
    var tabs = document.querySelector('paper-tabs');
    var pages = document.querySelector('core-pages');

    window.CosmosTabs = tabs;

    tabs.addEventListener('core-select', function() {
        pages.selected = tabs.selected;

        $('core-pages div.page').each(function(i, elem) {
            if (i == pages.selected) {
                $(elem).show();
            } else if (i != pages.selected) {
                $(elem).hide();
            }
        });

        switch (pages.selected) {
            case 0:
                var pageCosmos = document.querySelector('page-cosmos');
                Cosmos.request.getPlanets(function(json) {
                    pageCosmos.planets = json;
                }, function(jqXHR) {
                    // Failed
                    alert('request failed - ' + jqXHR.responseText);
                });
                Cosmos.request.getContainers(null, function(json) {
                    // Succeed
                    pageCosmos.containers = json;
                }, function(jqXHR) {
                    // Failed
                    alert('request failed - ' + jqXHR.responseText);
                });
                break;

            case 1:
                var pagePlanet = document.querySelector('page-planet');
                Cosmos.request.getPlanets(function(json) {
                    // for (var key in json) {
                    //     json[key]['containers'] = {
                    //         'core.cosmos': {
                    //             'Name': 'cosmos'
                    //         }
                    //     };
                    // }
                    pagePlanet.planets = json;
                }, function(jqXHR) {
                    // Failed
                    alert('request failed - ' + jqXHR.responseText);
                });
                break;
            case 2:
                var pageContainer = document.querySelector('page-container');
                Cosmos.request.getContainers(null, function(json) {
                    // Succeed
                    pageContainer.containers = json;
                }, function(jqXHR) {
                    // Failed                    
                    alert('request failed - ' + jqXHR.responseText);
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
        Route.defaultRoute = function() {
            window.location.replace("/cosmos");
        };
        Route.regist();
    });
}());
