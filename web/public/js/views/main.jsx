var PlanetButton = React.createClass({
	handleClick: function(e) {
		React.render(
			<ContainerButtonBox planetName={this.props.planetName}/>,
			document.getElementById('containerList')
		);
	}, 
	render: function() {
		return (
			<div className="clickable" style={{fontSize:"20px"}} onClick={this.handleClick}>
			   <a href="#" style={{"display": "block"}}>{this.props.planetName}</a>
			</div>
		);
	}
});

var PlanetButtonList = React.createClass({
	componentDidUpdate: function() {
		$(this.getDOMNode()).find('.clickable:first').click();
	},
	render: function() {
		var self = this;
		var nodes = self.props.data.map(function (planet) {
			return (
				<PlanetButton key={planet.name} planetName={planet.name} />
			);
		});
		return (
			<div>
			   {nodes}
			</div>
		);
	}
});

var PlanetButtonBox = React.createClass({
	getInitialState: function() {
		return {data: []};
	},
	componentDidMount: function() {
		var self = this;

		Cosmos.request.getPlanets(function(json, textStatus, jqXHR) {
			self.setState({data: json});
		});
	},
	render: function() {
		return (
			<div>
			   <h1>Planets</h1>
			   <PlanetButtonList data={this.state.data}/>
			</div>
		);
	}
});

var MetricGraph = React.createClass({
	getInitialState: function() {
		return {data: {}};
	},
	componentDidMount: function() {
		var self = this;
		Cosmos.request.getContainerInfo(self.props.planetName, self.props.containerName, '7d', function(json) {
			var cpuUsageData = [], memUsageData = [];
			var timeLabel = [];
			for (var i = 0; i < json.length; i++) {
				cpuUsageData.push(json[i].cpu_usage);
				memUsageData.push(json[i].mem_usage);
				var d = new Date(json[i].time * 1000)				
				timeLabel.push(d.getMonth() + '/' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes());
			}

			Cosmos.drawGraph('#chart1', 400, 300, timeLabel, cpuUsageData)
			Cosmos.drawGraph('#chart2', 400, 300, timeLabel, memUsageData)
//			self.setState({data:json})
		});
		
	},
	render: function() {
		return (
			<div>
			   <h1>{this.props.containerName}</h1>
			   <div className="row">
			      <div id="chart1" className="col-md-6 col-xs-6">
			         <h5>CPU usage</h5>
			      </div>
			      <div id="chart2" className="col-md-6 col-xs-6">
			         <h5>Memory usage</h5>
			      </div>
			   </div>
			   
			</div>
		);
	}			
});

var ContainerButton = React.createClass({
	handleClick: function(e) {		
		React.render(
			<MetricGraph containerName={this.props.containerName} planetName={this.props.planetName}/>,
			document.getElementById('graphBox')
		);
	}, 
	render: function() {
		return (
			<tr className="clickable" onClick={this.handleClick}>
			   <td>
			      <a href="#">{this.props.containerName}</a>
			   </td>
			   <td>...</td>
			   <td>...</td>
   			   <td>...</td>
   			   <td>...</td>
			</tr>
		);
	}	
});

var ContainerButtonList = React.createClass({
	componentDidMount: function() {
		$(this.getDOMNode()).find('.clickable:first').click();
	},	
	render: function() {
		var self = this;
		var nodes = self.props.data.map(function (container) {
			return (
				<ContainerButton key={container.name} containerName={container.name} planetName={self.props.planetName}/>
			);
		});
		return (			
		      <tbody>
				 {nodes}
		      </tbody>
		);
	}
});

var ContainerButtonBox = React.createClass({
	getInitialState: function() {
		return {data: []};
	}, 
	componentDidMount: function() {
		var self = this;
		Cosmos.request.getContainers(this.props.planetName, '7d', function(json, textStatus, jqXHR) {
			self.setState({data: json});
		});
	},
	componentDidUpdate: function() {
		$(this.getDOMNode()).find('.clickable:first').click();
	},
	render: function() {
		return (
			<div>
			   <h3>Containers</h3>
			   <table className="table">
		          <thead>
		             <tr>
		               <th>Name</th>
		               <th>Status</th>
		               <th>ImageID</th>
		               <th>StartAt</th>
					   <th>Description</th>
		             </tr>
		          </thead>
  		   	      <ContainerButtonList data={this.state.data} planetName={this.props.planetName}/>			   	  
			   </table>
			</div>
		);
	}
});

var MainPage = React.createClass({
	componentDidMount: function() {

	}, 
	render: function() {
		return (
		   <div className="row">
		      <div className="col-md-2 col-xs-2">
	  		     <PlanetButtonBox />
   			  </div>
		      <div className="col-md-10 col-xs-10">
			     <div id="graphBox"></div>
			     <div id="containerList"></div>
	     	  </div>
		   </div>
		);
	}
})

React.render(
	<MainPage />,
	document.getElementById('page')
);

