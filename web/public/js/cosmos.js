(function(){
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

  function getPlanets() {
    var container = $('<div>');
    container.append($('<script>', {src: '/js/bower_components/Chart.js/Chart.js'}));
    var chart = $('<canvas>', {id: 'myChart'})
    chart.attr({width: '800'});
    chart.attr({height: '400'});
    container.append(chart);
    $('#page').append(container);
    graph();
  }

  $.ajax({
    type: 'GET',
    accept: 'application/json',
    url: '/planets',
    success: function(response){
      console.log('success');
      getPlanets();
    },
    error: function(response){
      console.log('error');
      getPlanets();
    }
  });

  function graph() {
    var ctx = document.getElementById("myChart").getContext("2d");
    var myChart = new Chart(ctx).Line(data);
  }
})();