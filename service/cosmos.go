package service

import (
	"fmt"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/influxdbc"
)

type CosmosService struct {
	dbc      *influxdbc.InfluxDB
	lifeTime string
}

func NewCosmosService(dbc *influxdbc.InfluxDB) *CosmosService {
	return &CosmosService{dbc: dbc, lifeTime: "10m"}
}

func (this *CosmosService) GetContainers(token string) (interface{}, error) {
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

func (this *CosmosService) GetContainersOfPlanet(token, planet string) (interface{}, error) {
	dbQuery := fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\.%s\\./ WHERE time > now() - %s LIMIT 1", token, planet, this.lifeTime)
	series, err := this.dbc.Query(dbQuery, "s")
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromContainerSeries(token, planet, series), nil
}

func (this *CosmosService) GetContainerInfo(token, planetName, containerName string) (interface{}, error) {
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
