package dao

import "github.com/cosmos-io/influxdbc"

var (
	Container *ContainerDao
	NewsFeed  *NewsFeedDao
	Planet    *PlanetDao
	Series    *SeriesDao
)

func Initialize(dbc *influxdbc.InfluxDB) {
	Container = &ContainerDao{dbc: dbc}
	NewsFeed = &NewsFeedDao{dbc: dbc}
	Planet = &PlanetDao{dbc: dbc}
	Series = &SeriesDao{dbc: dbc}
}
