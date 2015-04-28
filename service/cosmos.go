package service

import (
	"fmt"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/cosmos/dao"
	"github.com/cosmos-io/influxdbc"
)

type CosmosService struct {
	newsFeedService *NewsFeedService
	LifeTime        int
}

func NewCosmosService(lifeTime int) *CosmosService {
	return &CosmosService{LifeTime: lifeTime, newsFeedService: &NewsFeedService{}}
}

func (this *CosmosService) GetContainers(token string) (map[string]map[string]interface{}, error) {
	series, err := dao.Container.GetContainers(token, this.LifeTime)
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
	savedContainerSeries, err := dao.Container.GetContainersOfPlanet(token, planet, false, this.LifeTime)
	if err != nil {
		return err
	}
	savedPlanetSeries, err := dao.Planet.GetPlanetStatusesInLifeTimeOfUser(token, this.LifeTime)
	if err != nil {
		return err
	}
	feedSeries, err := this.newsFeedService.PostNewsFeedIfNeeded(token, planet, savedPlanetSeries, savedContainerSeries, series)
	if err != nil {
		return err
	}
	series = append(series, feedSeries...)

	// Add Host metric series
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Name", token, planet), "txt_value", "num_value")
	planetSeries.AddPoint(planet, nil)
	series = append(series, planetSeries)

	planetSeries = influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Status", token, planet), "txt_value", "num_value")
	planetSeries.AddPoint("Up", nil)
	series = append(series, planetSeries)

	// Add Container series
	_, err = dao.Series.WriteSeries(series)
	if err != nil {
		return err
	}
	return nil
}

func (this *CosmosService) GetContainersOfPlanet(token, planet string, useRollup bool) (map[string]map[string]interface{}, error) {
	series, err := dao.Container.GetContainersOfPlanet(token, planet, useRollup, this.LifeTime)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(token, planet, series), nil
}

func (this *CosmosService) GetContainerMetrics(token, planetName, containerName string, metric []string, period string) (map[string]interface{}, error) {
	series, err := dao.Container.GetContainerMetrics(token, planetName, containerName, metric, period)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerMetricSeries(token, planetName, containerName, series), nil
}

func (this *CosmosService) GetPlanets(token string) (interface{}, error) {
	series, err := dao.Planet.GetPlanets(token, this.LifeTime)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromPlanetSeries(token, series), nil
}

func (this *CosmosService) GetPlanetMetrics(token, planetName string, metric []string) (interface{}, error) {
	series, err := dao.Container.GetContainerMetrics(token, planetName, "", metric, "10m")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(token, planetName, series), nil
}

func (this *CosmosService) GetNewsFeeds(token, time string) (interface{}, error) {
	series, err := dao.NewsFeed.GetNewsFeeds(token, time)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromNewsFeedSeries(series), nil
}
