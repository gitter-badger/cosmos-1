package dao

import "github.com/cosmos-io/influxdbc"

type SeriesDao struct {
	dbc *influxdbc.InfluxDB
}

func (this *SeriesDao) WriteSeries(series []*influxdbc.Series) (string, error) {
	// Insert series
	return this.dbc.WriteSeries(series, "s")
}
