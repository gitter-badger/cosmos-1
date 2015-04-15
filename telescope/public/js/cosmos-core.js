// Cosmos object initialize
(function() {
    window.Cosmos = {};
    Cosmos.API_VER = 'v1';

    Cosmos.loadScripts = function(scripts, type) {
        for (var i in scripts) {
            if (type) {
                $(document.body).append($('<script>', {
                    src: scripts[i],
                    type: type
                }));
            } else {
                $(document.body).append($('<script>', {
                    src: scripts[i],
                    type: 'text/javascript'
                }));
            }
        }
    };

    Cosmos.request = {
        getPlanets: function(done, fail, complete) {
            var xhr = $.ajax({
                url: '/' + Cosmos.API_VER + '/planets',
                method: 'GET',
                accept: 'application/json',
                dataType: 'json'
            });
            if (typeof done == 'function') {
                xhr.done(done);
            }
            if (typeof fail == 'function') {
                xhr.fail(fail);
            }
            if (typeof complete == 'function') {
                xhr.complete(complete)
            }
        },
        _convertContainerResponse: function(json) {
            var data = [];
            var keys = Object.keys(json);
            for (var i = 0; i < keys.length; i++) {
                var k = keys[i];
                var inKeys = Object.keys(json[k]);

                for (var j = 0; j < inKeys.length; j++) {
                    var inK = inKeys[j];
                    var newK = inK.replace(/\./g, "");
                    var val = json[k][inK];
                    delete(json[k][inK]);
                    json[k][newK] = val
                }
                json[k]['Key'] = k;
                var comps = k.split('.');
                json[k]['Planet'] = comps[0];
                json[k]['Container'] = comps[1];
                data.push(json[k]);
            }
            return data;
        },
        _convertContainerInfoResponse: function(json) {
            var keys = Object.keys(json);
            for (var i = 0; i < keys.length; i++) {
                var k = keys[i];
                var newK = k.replace(/\./g, "");
                var val = json[k];
                delete(json[k]);
                json[newK] = val
            }
            return json;
        },
        getContainers: function(planet, done, fail, complete) {
            var self = this;
            var url = '/' + Cosmos.API_VER;
            if (planet) {
                url = url + '/planets/' + planet + '/containers';
            } else {
                url = url + '/containers';
            }

            var xhr = $.ajax({
                url: url,
                method: 'GET',
                accept: 'application/json',
                dataType: 'json'
            });

            if (typeof done == 'function') {
                xhr.done(function(json) {
                    done(self._convertContainerResponse(json));
                });
            }
            if (typeof fail == 'function') {
                xhr.fail(fail);
            }
            if (typeof complete == 'function') {
                xhr.complete(complete)
            }
        },
        getContainerInfo: function(planet, containerName, done, fail, complete) {
            var self = this;
            var xhr = $.ajax({
                url: '/' + Cosmos.API_VER + '/planets/' + planet + '/containers/' + containerName,
                method: 'GET',
                accept: 'application/json',
                dataType: 'json'
            });

            if (typeof done == 'function') {
                xhr.done(function(json) {
                    done(self._convertContainerInfoResponse(json));
                });
            }
            if (typeof fail == 'function') {
                xhr.fail(fail);
            }
            if (typeof complete == 'function') {
                xhr.complete(complete)
            }
        }
    };

    Cosmos.drawGraph = function(selector, width, height, labels, data) {
        var dataset = {
            labels: [],
            datasets: [{
                label: "Metrics",
                fillColor: "rgba(33,133,117,0.2)",
                strokeColor: "rgba(33,133,117,1)",
                pointColor: "rgba(33,133,117,1)",
                pointStrokeColor: "#fff",
                pointHighlightFill: "#fff",
                pointHighlightStroke: "rgba(33,133,117,1)",
                data: []
            }]
        };
        if (labels) {
            dataset.labels = labels;
        }
        if (data) {
            dataset.datasets[0].data = data;
        }

        var container = $('<div/>').css({
            'text-align': 'center'
        });
        var chart = $('<canvas/>').css({
            'display': 'inline-block'
        }).attr('height', height);
        var opt = {
            bezierCurve: false,
            pointHitDetectionRadius: 4,
            scaleLabel: "<%= ' ' + value%>",
            scaleFontSize: 12,
            scaleFontColor: '#888'
        };

        if (isNaN(width)) {
            opt['responsive'] = true;
            opt['maintainAspectRatio'] = false;
        } else {
            chart.attr('width', width);
        }

        container.append(chart);
        $(selector).append(container);

        var ctx = chart[0].getContext("2d");
        new Chart(ctx).Line(dataset, opt);
    };
})();
