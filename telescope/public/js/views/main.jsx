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

			var set = json['Stats.Cpu.TotalUtilization'];
			for (var i = set.length-1; i >= 0; i--) {
				cpuUsageData.push(set[i][2]);
				var d = new Date(set[i][0] * 1000)
				timeLabel.push(d.getMonth() + '/' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes());
			}

			set = json['Stats.Memory.Usage'];
			for (var i = set.length-1; i >= 0; i--) {
				memUsageData.push(set[i][2]);
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
			<MetricGraph containerName={this.props.data['Names.0'][3]} containerId={this.props.data['Id'][3]} planetName={this.props.planetName}/>,
			document.getElementById('graphBox')
		);
	}, 
	render: function() {
		var data = this.props.data;
		return (
			<tr className="clickable" onClick={this.handleClick}>
			   <td>
			      <a href="#">{data['Names.0'][3]}</a>
			   </td>
			   <td>{data['Status'][3]}</td>
			   <td>{data['Id'][3]}</td>
   			   <td>{data['Ports.0.PublicPort'][2] + ":" + data['Ports.0.PrivatePort'][2] + " " + data['Ports.0.Type'][3]}</td>
   			   <td>{data['Command'][3]}</td>
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
		var contIds = [];
		for (var key in self.props.data) {
			contIds.push(key);
		}

		var nodes = contIds.map(function (containerId) {
			return (
				<ContainerButton key={containerId} data={self.props.data[containerId]} planetName={self.props.planetName}/>
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

