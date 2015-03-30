// Cosmos object initialize
(function(){
  window.Cosmos = {};
  Cosmos.API_VER = 'v1';

  Cosmos.loadScripts = function(scripts, type) {
    for (var i in scripts) {
      if (type) {
        $(document.body).append($('<script>', {src: scripts[i], type: type}));
      } else {
        $(document.body).append($('<script>', {src: scripts[i], type: 'text/javascript'}));
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
    getContainers: function(planet, timeToLive, done, fail, complete) {
      var xhr = $.ajax({
        url: '/' + Cosmos.API_VER + '/planets/' + planet + '/containers',
        method: 'GET',
        accept: 'application/json',
        dataType: 'json',
        data: {ttl: timeToLive}
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
    getContainerInfo: function(planet, container, timeInterval, done, fail, complete) {
      var xhr = $.ajax({
        url: '/' + Cosmos.API_VER + '/planets/' + planet + '/containers/' + container,
        method: 'GET',
        accept: 'application/json',
        dataType: 'json',
        data: {interval: timeInterval}
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
    }
  };

  Cosmos.drawGraph = function(selector, width, height, labels, data) {
    var dataset = {
      labels: [],
      datasets: [
      {
        label: "container metric",
        fillColor: "rgba(220,220,220,0.2)",
        strokeColor: "rgba(220,220,220,1)",
        pointColor: "rgba(220,220,220,1)",
        pointStrokeColor: "#fff",
        pointHighlightFill: "#fff",
        pointHighlightStroke: "rgba(220,220,220,1)",
        data: []
      }
      ]
    };
    if (labels) {
      dataset.labels = labels;
    }
    if (data) {
      dataset.datasets[0].data = data;
      console.log(dataset);
    }


    var container = $('<div/>');
    var chart = $('<canvas/>')
    chart.attr({width: width});
    chart.attr({height: height});
    container.append(chart);
    $(selector).append(container);

    var ctx = chart[0].getContext("2d");
    new Chart(ctx).Line(dataset);
  };
})();


var Page = {};
Page.planetList = function() {
    Cosmos.request.getPlanets(function(json, textStatus, jqXHR) {
      var page = $('#page');
      page.addClass('planet-list');
      page.append($('<h4/>').text('Planets'));

      var divRow = $('<div/>').addClass('row');
      
      for (var i = 0; i < 12; i++) {
        var divCol = $('<div/>').addClass('col-md-1 col-xs-2');
        var a = $('<a/>', {href: '#'});
        a.text(json[0].name).click(function(e) {          
          return false;
        });
        divCol.append(a);
        divRow.append(divCol);        
      }
      page.append(divRow);
    }, function(jqXHR, textStatus, errorThrown) {
      console.log(jqXHR.responseText);
    });

};

Page.planetDetail = function() {
    

    var page = $('#page');
    page.addClass('planet-detail');

    page.append($('<h4/>').text(Route.params['planet']));
    Cosmos.request.getContainers(Route.params['planet'], '7d', function(json, textStatus, jqXHR) {
      var divRow = $('<div/>').addClass('row');
      var divLeft = $('<div/>').addClass('col-md-4');
      var divRight = $('<div/>').addClass('col-md-8');

      divRow.append(divLeft).append(divRight);

      var ul = $('<ul />');
      ul.addClass('list');
      for (var i = 0;i < json.length; i++) {
        var li = $('<li/>');
        var a = $('<a/>', {href: '#'});
        a.text(json[i].name).click(function(e) {
          Cosmos.drawGraph(divRight);
          return false;
        });
        li.append(a);
        ul.append(li);
      }

      divLeft.append(ul);
      page.append(divRow);
    }, 
    function(jqXHR, textStatus, errorThrown) {
      console.log(errorThrown);
    })

    //Cosmos.drawGraph('#page');
};


// Route configuration
(function(){
  // Route.match('/', Page.planetList);
  Route.match('/', function() {
    Cosmos.loadScripts(['/vendor/chartjs-1.0.2/chart.min.js'], "text/javascript");
    Cosmos.loadScripts(['/js/views/main.jsx'], "text/jsx");
  });
  Route.match('/planets/:planet', Page.planetDetail);
  Route.regist();
})();