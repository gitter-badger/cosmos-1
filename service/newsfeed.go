package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/influxdbc"
)

type NewsFeedService struct {
}

var (
	FEED_TYPE_ADD_CONTAINER    = 0
	FEED_TYPE_REMOVE_CONTAINER = 1
	FEED_TYPE_ADD_PLANET       = 2
	FEED_TYPE_REMOVE_PLANET    = 3
)

func (this *NewsFeedService) PostNewsFeedIfNeeded(token, planet string, savedPlanetSeries, savedContainerSeries, newContainerSeries []*influxdbc.Series) ([]*influxdbc.Series, error) {
	newContainers := converter.ConvertFromContainerSeries(token, planet, newContainerSeries)
	savedContainers := converter.ConvertFromContainerSeries(token, planet, savedContainerSeries)

	addSeries, err := this.checkContainersAdded(token, savedContainers, newContainers)
	if err != nil {
		return nil, err
	}
	removeSeries, err := this.checkContainersRemoved(token, savedContainers, newContainers)
	if err != nil {
		return nil, err
	}

	planetSeries, err := this.checkPlanetAdded(token, planet, savedPlanetSeries)
	if err != nil {
		return nil, err
	}

	if planetSeries != nil {
		addSeries = append(addSeries, planetSeries)
	}
	if len(removeSeries) > 0 {
		addSeries = append(addSeries, removeSeries...)
	}

	return addSeries, nil
}

func (this *NewsFeedService) checkPlanetAdded(token, planet string, savedPlanets []*influxdbc.Series) (*influxdbc.Series, error) {
	added := true
	for _, series := range savedPlanets {
		if strings.HasPrefix(series.Name, fmt.Sprintf("PLANET.%s.%s", token, planet)) {
			status := series.Points[0][3].(string)
			if strings.Contains(status, "Up") {
				added = false
				break
			}
		}
	}
	if added {
		series, err := this.MakeNewsFeed(token, planet, FEED_TYPE_ADD_PLANET)
		if err != nil {
			return nil, err
		}
		return series, nil
	}

	return nil, nil
}

func (this *NewsFeedService) checkContainersAdded(token string, savedContainers, newContainers map[string]map[string]interface{}) ([]*influxdbc.Series, error) {
	series := make([]*influxdbc.Series, 0)
	for key, cont := range newContainers {
		added := false

		statusRow := cont["Status"].([][]interface{})[0]
		status := statusRow[0].(string)
		if strings.Contains(status, "Up") == false {
			continue
		}

		savedCont, exist := savedContainers[key]

		if exist {
			savedStatusRow := savedCont["Status"].([][]interface{})[0]
			status = savedStatusRow[3].(string)
			if strings.Contains(status, "Up") == false {
				added = true
			}
		} else {
			added = true
		}

		if added {
			// Add new one
			feedSeries, err := this.MakeNewsFeed(token, key, FEED_TYPE_ADD_CONTAINER)
			if err != nil {
				return nil, err
			}
			series = append(series, feedSeries)
		}
	}

	return series, nil
}

func (this *NewsFeedService) checkContainersRemoved(token string, savedContainers, newContainers map[string]map[string]interface{}) ([]*influxdbc.Series, error) {
	series := make([]*influxdbc.Series, 0)
	for key, savedCont := range savedContainers {
		removed := false

		savedStatusRow := savedCont["Status"].([][]interface{})[0]
		savedStatus := savedStatusRow[3].(string)
		if strings.Contains(savedStatus, "Up") == false {
			continue
		}

		cont, exist := newContainers[key]
		if exist {
			statusRow := cont["Status"].([][]interface{})[0]
			status := statusRow[0].(string)
			if strings.Contains(status, "Up") == false {
				removed = true
			}
		} else {
			removed = true

			// if curiosity can't send exited container info, then make one explicitly
			statusSeries := influxdbc.NewSeries(fmt.Sprintf("CONTAINER.%s.%s.Status", token, key), "txt_value", "num_value")
			statusSeries.AddPoint("Exited", nil)
			series = append(series, statusSeries)
		}
		if removed {
			// Removed one
			feedSeries, err := this.MakeNewsFeed(token, key, FEED_TYPE_REMOVE_CONTAINER)
			if err != nil {
				return nil, err
			}

			series = append(series, feedSeries)
		}
	}

	return series, nil
}

func (this *NewsFeedService) MakeNewsFeed(token, key string, feedType int) (*influxdbc.Series, error) {
	var (
		data map[string]interface{}
	)

	data = make(map[string]interface{})
	data["User"] = token

	switch feedType {
	case FEED_TYPE_REMOVE_CONTAINER:
		comps := strings.Split(key, ".")
		data["Planet"] = comps[0]
		data["Container"] = comps[1]
	case FEED_TYPE_ADD_CONTAINER:
		comps := strings.Split(key, ".")
		data["Planet"] = comps[0]
		data["Container"] = comps[1]
	case FEED_TYPE_ADD_PLANET:
		data["Planet"] = key
	case FEED_TYPE_REMOVE_PLANET:
		data["Planet"] = key
	}
	data["Type"] = feedType

	raw, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	jsonData := string(raw)
	feedSeries := influxdbc.NewSeries(fmt.Sprintf("NEWSFEED.%s.%s", token, key), "type", "value")
	feedSeries.AddPoint(feedType, jsonData)

	return feedSeries, nil
}
