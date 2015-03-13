(function(){
  function getPlanets() {
    var container = $('<div>');
    $('#page').append(container);
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
})();