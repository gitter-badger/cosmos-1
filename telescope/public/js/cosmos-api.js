// Cosmos object initialize
(function() {
    window.Cosmos = {
        API_VER: 'v1',
        getPlanets: function(done, fail, complete) {
            var self = this;
            var xhr = $.ajax({
                url: '/' + Cosmos.API_VER + '/planets',
                method: 'GET',
                accept: 'application/json',
                dataType: 'json'
            });
            if (typeof done == 'function') {
                xhr.done(function(json) {
                    done(self._convertPlanetResponse(json));
                });
            }
            if (typeof fail == 'function') {
                xhr.fail(fail);
            }
            if (typeof complete == 'function') {
                xhr.complete(complete)
            }
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
        },
        getNewsFeeds: function(done, fail, complete) {
            var self = this;
            var xhr = $.ajax({
                url: '/' + Cosmos.API_VER + '/newsfeeds',
                method: 'GET',
                accept: 'application/json',
                dataType: 'json'
            });

            if (typeof done == 'function') {
                xhr.done(function(json) {
                    done(self._convertNewsFeedsResponse(json));
                });
            }
            if (typeof fail == 'function') {
                xhr.fail(fail);
            }
            if (typeof complete == 'function') {
                xhr.complete(complete)
            }
        },
        _convertNewsFeedsResponse: function(json) {
            var data = [];
            for (var i = 0; i < json.length; i++) {
                var j = JSON.parse(json[i][2])
                j['Time'] = json[i][0];
                j['Key'] = j.Planet;
                if (j.Container) {
                    j['Key'] += '.' + j.Container;
                } 
                if (j['Type'] == 2 || j['Type'] == 3) {
                    // Planet NewsFeed
                    // set hidden property to TRUE
                    j['Hidden'] = true;
                }
                data.push(j);
            }
            console.log(data);
            return data;
        },
        _convertPlanetResponse: function(json) {
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
                json[k]['Planet'] = k;
                data.push(json[k]);
            }
            return data;
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
        }
    };
})();
