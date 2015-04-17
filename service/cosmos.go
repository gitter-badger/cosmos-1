package service

import (
	"fmt"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/influxdbc"
)

type CosmosService struct {
	dbc             *influxdbc.InfluxDB
	newsFeedService *NewsFeedService
	lifeTime        string
}

var (
	FEED_TYPE_ADD_CONTAINER    = 0
	FEED_TYPE_REMOVE_CONTAINER = 1
)

func NewCosmosService(dbc *influxdbc.InfluxDB) *CosmosService {
	return &CosmosService{dbc: dbc, lifeTime: "10m", newsFeedService: NewNewsFeedService(dbc)}
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
	savedContainers, err := this.GetContainersOfPlanet(token, planet, false)
	if err != nil {
		return err
	}
	newContainers := converter.ConvertFromContainerSeries(token, planet, series)
	feedSeries, err := this.newsFeedService.PostContainerNewsFeedIfNeeded(token, planet, savedContainers, newContainers)
	if err != nil {
		return err
	}

	series = append(series, feedSeries...)

	// Add Host metric series
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Name", token, planet), "txt_value", "num_value")
	planetSeries.AddPoint(planet, nil)
	_, err = this.dbc.WriteSeries([]*influxdbc.Series{planetSeries}, "s")
	if err != nil {
		return err
	}

	series = append(series, planetSeries)

	// Add Container series
	_, err = this.dbc.WriteSeries(series, "s")
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

func (this *CosmosService) GetNewsFeeds(token, time string) (interface{}, error) {
	cond := ""
	if time != "" {
		cond = fmt.Sprintf("WHERE time > %s", time)
	}
	series, err := this.dbc.Query(fmt.Sprintf("SELECT value FROM merge(/^NEWSFEED\\.%s\\..*/) %s LIMIT 30", token, cond), "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromNewsFeedSeries(series), nil
}
