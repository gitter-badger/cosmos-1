package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos-io/influxdbc"
)

type NewsFeedService struct {
	dbc *influxdbc.InfluxDB
}

func NewNewsFeedService(dbc *influxdbc.InfluxDB) *NewsFeedService {
	return &NewsFeedService{dbc: dbc}
}

func (this *NewsFeedService) PostContainerNewsFeedIfNeeded(token, planet string, savedContainers, newContainers map[string]map[string]interface{}) ([]*influxdbc.Series, error) {
	addSeries, err := this.checkContainersAdded(token, planet, savedContainers, newContainers)
	if err != nil {
		return nil, err
	}
	removeSeries, err := this.checkContainersStopped(token, planet, savedContainers, newContainers)
	if err != nil {
		return nil, err
	}

	return append(addSeries, removeSeries...), nil
}

func (this *NewsFeedService) checkContainersAdded(token, planet string, savedContainers, newContainers map[string]map[string]interface{}) ([]*influxdbc.Series, error) {
	series := make([]*influxdbc.Series, 0)
	for key, cont := range newContainers {
		added := false

		status := cont["Status"].([]interface{})[0].(string)
		if strings.Contains(status, "Up") == false {
			continue
		}

		savedCont, exist := savedContainers[key]

		if exist {
			status = savedCont["Status"].([]interface{})[3].(string)
			if strings.Contains(status, "Up") == false {
				added = true
			}
		} else {
			added = true
		}

		if added {
			// Add new one
			feedSeries, err := this.makeNewsFeed(token, key, FEED_TYPE_ADD_CONTAINER)
			if err != nil {
				return nil, err
			}
			series = append(series, feedSeries)
		}
	}

	return series, nil
}

func (this *NewsFeedService) checkContainersStopped(token, planet string, savedContainers, newContainers map[string]map[string]interface{}) ([]*influxdbc.Series, error) {
	series := make([]*influxdbc.Series, 0)
	for key, savedCont := range savedContainers {
		removed := false

		savedStatus := savedCont["Status"].([]interface{})[3].(string)
		if strings.Contains(savedStatus, "Up") == false {
			continue
		}

		cont, exist := newContainers[key]
		if exist {
			status := cont["Status"].([]interface{})[0].(string)
			if strings.Contains(status, "Up") == false {
				removed = true
			}
		} else {
			removed = true
		}
		if removed {
			// Removed one
			// FeedType 1

			feedSeries, err := this.makeNewsFeed(token, key, FEED_TYPE_REMOVE_CONTAINER)
			if err != nil {
				return nil, err
			}

			series = append(series, feedSeries)
		}
	}

	return series, nil
}

func (this *NewsFeedService) makeNewsFeed(token, key string, feedType int) (*influxdbc.Series, error) {
	var (
		msg  string
		data map[string]interface{}
	)

	switch feedType {
	case FEED_TYPE_REMOVE_CONTAINER:
		msg = "CONTAINER is removed - " + key
	case FEED_TYPE_ADD_CONTAINER:
		msg = "New CONTAINER is added - " + key
	}

	data = make(map[string]interface{})
	data["Content"] = msg
	data["User"] = token
	comps := strings.Split(key, ".")
	data["Planet"] = comps[0]
	data["Container"] = comps[1]

	raw, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	jsonData := string(raw)

	feedSeries := influxdbc.NewSeries(fmt.Sprintf("NEWSFEED.%s.%s", token, key), "type", "value")
	feedSeries.AddPoint(feedType, jsonData)

	return feedSeries, nil
}
