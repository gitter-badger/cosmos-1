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
		Cosmos.request.getContainerInfo(self.props.planetName, self.props.containerId, '1m', function(json) {
			var cpuUsageData = [], memUsageData = [];
			var timeLabel = [];
			for (var i = json.length-1; i >= 0; i--) {
				cpuUsageData.push(parseInt(json[i].cpu_usage));
				memUsageData.push(parseInt(json[i].mem_usage));
				var d = new Date(json[i].time * 1000)				
				timeLabel.push(d.getMonth() + '/' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes());
			}

			Cosmos.drawGraph('#chart1', 400, 300, timeLabel, cpuUsageData)
			Cosmos.drawGraph('#chart2', 400, 300, timeLabel, memUsageData)
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
			<MetricGraph containerName={this.props.container.name} containerId={this.props.container.id} planetName={this.props.planetName}/>,
			document.getElementById('graphBox')
		);
	}, 
	render: function() {
		return (
			<tr className="clickable" onClick={this.handleClick}>
			   <td>
			      <a href="#">{this.props.container.name}</a>
			   </td>
			   <td>{this.props.container.status}</td>
			   <td>{this.props.container.id}</td>
   			   <td>{this.props.container.port}</td>
   			   <td>{this.props.container.command}</td>
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
				<ContainerButton key={container.id} container={container} planetName={self.props.planetName}/>
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
		               <th>Ports</th>
					   <th>Command</th>
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

