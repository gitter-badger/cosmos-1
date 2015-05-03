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

var (
    DEFAULT_USER = "default"
)

func NewCosmosService(lifeTime int) *CosmosService {
	return &CosmosService{LifeTime: lifeTime, newsFeedService: &NewsFeedService{}}
}

func (this *CosmosService) GetContainers() (map[string]map[string]interface{}, error) {
	series, err := dao.Container.GetContainers(DEFAULT_USER, this.LifeTime)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(DEFAULT_USER, "", series), nil
}

func (this *CosmosService) AddContainersOfPlanet(planet string, data []byte) error {
	// Container metrics
	series, err := converter.ConvertToContainerSeries(DEFAULT_USER, planet, data)
	if err != nil {
		return err
	}
	savedContainerSeries, err := dao.Container.GetContainersOfPlanet(DEFAULT_USER, planet, false, this.LifeTime)
	if err != nil {
		return err
	}
	savedPlanetSeries, err := dao.Planet.GetPlanetStatusesInLifeTimeOfUser(DEFAULT_USER, this.LifeTime)
	if err != nil {
		return err
	}
	feedSeries, err := this.newsFeedService.PostNewsFeedIfNeeded(DEFAULT_USER, planet, savedPlanetSeries, savedContainerSeries, series)
	if err != nil {
		return err
	}
	series = append(series, feedSeries...)

	// Add Host metric series
	planetSeries := influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Name", DEFAULT_USER, planet), "txt_value", "num_value")
	planetSeries.AddPoint(planet, nil)
	series = append(series, planetSeries)

	planetSeries = influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Status", DEFAULT_USER, planet), "txt_value", "num_value")
	planetSeries.AddPoint("Up", nil)
	series = append(series, planetSeries)

	// Add Container series
	_, err = dao.Series.WriteSeries(series)
	if err != nil {
		return err
	}
	return nil
}

func (this *CosmosService) GetContainersOfPlanet(planet string, useRollup bool) (map[string]map[string]interface{}, error) {
	series, err := dao.Container.GetContainersOfPlanet(DEFAULT_USER, planet, useRollup, this.LifeTime)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(DEFAULT_USER, planet, series), nil
}

func (this *CosmosService) GetContainerMetrics(planetName, containerName string, metric []string, period string) (map[string]interface{}, error) {
	series, err := dao.Container.GetContainerMetrics(DEFAULT_USER, planetName, containerName, metric, period)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerMetricSeries(DEFAULT_USER, planetName, containerName, series), nil
}

func (this *CosmosService) GetPlanets() (interface{}, error) {
	series, err := dao.Planet.GetPlanets(DEFAULT_USER, this.LifeTime)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromPlanetSeries(DEFAULT_USER, series), nil
}

func (this *CosmosService) GetPlanetMetrics(planetName string, metric []string) (interface{}, error) {
	series, err := dao.Container.GetContainerMetrics(DEFAULT_USER, planetName, "", metric, "10m")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(DEFAULT_USER, planetName, series), nil
}

func (this *CosmosService) GetNewsFeeds(time string) (interface{}, error) {
	series, err := dao.NewsFeed.GetNewsFeeds(DEFAULT_USER, time)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromNewsFeedSeries(series), nil
}
