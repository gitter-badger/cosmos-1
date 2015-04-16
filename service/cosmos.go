package service

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/influxdbc"
)

type CosmosService struct {
	dbc      *influxdbc.InfluxDB
	lifeTime string
}

var (
	FEED_TYPE_ADD_CONTAINER    = 0
	FEED_TYPE_REMOVE_CONTAINER = 1
)

func NewCosmosService(dbc *influxdbc.InfluxDB) *CosmosService {
	return &CosmosService{dbc: dbc, lifeTime: "10m"}
}

func (this *CosmosService) postContainerNewsFeedIfNeeded(token, planet string, series []*influxdbc.Series) error {
	savedContainers, err := this.GetContainersOfPlanet(token, planet, false)
	if err != nil {
		return err
	}

	newContainers := converter.ConvertFromContainerSeries(token, planet, series)

	for key, _ := range newContainers {
		if _, ok := savedContainers[key]; !ok {
			// Add new one
			var msg = "New CONTAINER is added! - " + key
			fmt.Println(msg)
			data := make(map[string]interface{})
			data["content"] = msg
			this.AddNewsFeedOfContainer(token, key, FEED_TYPE_ADD_CONTAINER, data)
		}
	}

	for key, _ := range savedContainers {
		if _, ok := newContainers[key]; !ok {
			// Removed one
			// FeedType 1
			var msg = "CONTAINER is removed - " + key
			fmt.Println(msg)
			data := make(map[string]interface{})
			data["content"] = msg
			this.AddNewsFeedOfContainer(token, key, FEED_TYPE_REMOVE_CONTAINER, data)
		}
	}

	return nil
}

func (this *CosmosService) AddNewsFeedOfContainer(token, key string, feedType int, data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonData := string(raw)

	feedSeries := influxdbc.NewSeries(fmt.Sprintf("NEWSFEED.%s.%s", token, key), "type", "value")
	feedSeries.AddPoint(feedType, jsonData)

	series := make([]*influxdbc.Series, 1)
	series[0] = feedSeries

	err = this.dbc.WriteSeries(series, "s")
	if err != nil {
		return err
	}

	return nil
}

func (this *CosmosService) GetContainers(token string) (map[string]map[string]interface{}, error) {
	series, err := this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\./ WHERE time > now() - %s LIMIT 1", token, this.lifeTime), "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(token, "", series), nil
}

func (this *CosmosService) AddContainersOfPlanet(token, planet string, data []byte) error {
	// Container metrics
	series, err := converter.ConvertToContainerSeries(token, planet, data)
	if err != nil {
		return err
	}

	// Post newsfeed for new container if needed
	this.postContainerNewsFeedIfNeeded(token, planet, series)

	// Host metrics
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Name", token, planet), "txt_value", "num_value")
	planetSeries.AddPoint(planet, nil)
	series = append(series, planetSeries)

	err = this.dbc.WriteSeries(series, "s")
	if err != nil {
		return err
	}

	return nil
}

func (this *CosmosService) GetContainersOfPlanet(token, planet string, useRollup bool) (map[string]map[string]interface{}, error) {
	var dbQuery string
	if useRollup {
		dbQuery = fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\.%s\\./ WHERE time > now() - %s LIMIT 1", token, planet, this.lifeTime)
	} else {
		dbQuery = fmt.Sprintf("SELECT num_value, txt_value FROM /^CONTAINER\\.%s\\.%s\\./ WHERE time > now() - %s LIMIT 1", token, planet, this.lifeTime)
	}

	series, err := this.dbc.Query(dbQuery, "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(token, planet, series), nil
}

func (this *CosmosService) GetContainerInfo(token, planetName, containerName string) (map[string]interface{}, error) {
	seriesName := converter.MakeContainerSeriesName(token, planetName, containerName)

	dbQuery := fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\./ WHERE time > now() - %s LIMIT 10", seriesName, this.lifeTime)
	series, err := this.dbc.Query(dbQuery, "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerInfoSeries(token, planetName, containerName, series), nil
}

func (this *CosmosService) GetPlanets(token string) (interface{}, error) {
	series, err := this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.PLANET\\.%s\\./ WHERE time > now() - %s LIMIT 1", token, this.lifeTime), "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromPlanetSeries(token, series), nil
}
