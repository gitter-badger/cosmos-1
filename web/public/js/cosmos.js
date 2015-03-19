// Cosmos object initialize
(function(){
  window.Cosmos = {};
  Cosmos.API_VER = 'v1';

  Cosmos.loadScripts = function(scripts) {
    for (var i in scripts) {
      $(document.body).append($('<script>', {src: scripts[i], type: 'text/javascript'}));
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
    getContainers: function(planet, timeInterval, done, fail, complete) {
      var xhr = $.ajax({
        url: '/' + Cosmos.API_VER + '/' + planet + '/containers',        
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

  Cosmos.drawGraph = function(selector) {
    var data = {
      labels: ["January", "February", "March", "April", "May", "June", "July"],
      datasets: [
      {
        label: "My First dataset",
        fillColor: "rgba(220,220,220,0.2)",
        strokeColor: "rgba(220,220,220,1)",
        pointColor: "rgba(220,220,220,1)",
        pointStrokeColor: "#fff",
        pointHighlightFill: "#fff",
        pointHighlightStroke: "rgba(220,220,220,1)",
        data: [65, 59, 80, 81, 56, 55, 40]
      },
      {
        label: "My Second dataset",
        fillColor: "rgba(151,187,205,0.2)",
        strokeColor: "rgba(151,187,205,1)",
        pointColor: "rgba(151,187,205,1)",
        pointStrokeColor: "#fff",
        pointHighlightFill: "#fff",
        pointHighlightStroke: "rgba(151,187,205,1)",
        data: [28, 48, 40, 19, 86, 27, 90]
      }
      ]
    };

    var container = $('<div/>');
    var chart = $('<canvas/>', {id: 'chart'})
    chart.attr({width: '800'});
    chart.attr({height: '400'});
    container.append(chart);
    $(selector).append(container);

    var ctx = document.getElementById("chart").getContext("2d");
    new Chart(ctx).Line(data);
  };
})();


var Page = {};
Page.planetList = function() {
    Cosmos.request.getPlanets(function(json, textStatus, jqXHR) {

      var ul = $('<ul/>');
      for (var i = 0; i < json.length; i++) {
        var li = $('<li/>');
        var a = $('<a/>', {href: '/planets/' + json[i].name});
        a.text(json[i].name);
        li.append(a);
        ul.append(li);
      }

      $('#page').append(ul);

    }, function(jqXHR, textStatus, errorThrown) {
      console.log(jqXHR.responseText);
    });

};

Page.planetDetail = function(params) {
    var scripts = ['/js/bower_components/Chart.js/Chart.js'];
    Cosmos.loadScripts(scripts);

    $('#page').append($('<h4/>').text(params['planet']));
    Cosmos.request.getContainers(params['planet'], '7d', function(json, textStatus, jqXHR) {
      console.log(json)
    }, 
    function(jqXHR, textStatus, errorThrown) {
      console.log(errorThrown);
    })

    Cosmos.drawGraph('#page');
};


// Route configuration
(function(){
  Route.match('/', Page.planetList);
  Route.match('/planets/:planet', Page.planetDetail); 
  Route.regist();
})();