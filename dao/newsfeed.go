package dao

import (
	"fmt"

	"github.com/cosmos-io/influxdbc"
)

type NewsFeedDao struct {
	dbc *influxdbc.InfluxDB
}

func (this *NewsFeedDao) GetNewsFeeds(token, time string) ([]*influxdbc.Series, error) {
	cond := ""
	if time != "" {
		cond = fmt.Sprintf("WHERE time > %s", time)
	}
	series, err := this.dbc.Query(fmt.Sprintf("SELECT value FROM merge(/^NEWSFEED\\.%s\\..*/) %s LIMIT 30", token, cond), "s")
	if err != nil {
		return nil, err
	}

	return series, nil
}
