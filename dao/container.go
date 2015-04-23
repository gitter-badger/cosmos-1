package dao

import (
	"fmt"

	"github.com/cosmos-io/cosmos/converter"
	"github.com/cosmos-io/influxdbc"
)

type ContainerDao struct {
	dbc *influxdbc.InfluxDB
}

func (this *ContainerDao) GetContainers(token string, lifeTime int) ([]*influxdbc.Series, error) {
	series, err := this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\./ WHERE time > now() - %dm LIMIT 1", token, lifeTime), "s")
	if err != nil {
		return nil, err
	}

	return series, nil
}

func (this *ContainerDao) GetContainersOfPlanet(token, planet string, useRollup bool, lifeTime int) ([]*influxdbc.Series, error) {
	var dbQuery string
	if useRollup {
		dbQuery = fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s\\.%s\\./ WHERE time > now() - %dm LIMIT 1", token, planet, lifeTime)
	} else {
		dbQuery = fmt.Sprintf("SELECT num_value, txt_value FROM /^CONTAINER\\.%s\\.%s\\./ WHERE time > now() - %dm LIMIT 1", token, planet, lifeTime)
	}

	series, err := this.dbc.Query(dbQuery, "s")
	if err != nil {
		return nil, err
	}

	return series, nil
}

func (this *ContainerDao) GetContainerMetrics(token, planetName, containerName, metric string) ([]*influxdbc.Series, error) {
	seriesName := converter.MakeContainerSeriesName(token, planetName, containerName)
	if metric != "all" {
		seriesName = seriesName + "\\." + metric
	} else {
		seriesName = seriesName + "\\."
	}

	dbQuery := fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.CONTAINER\\.%s/ LIMIT 10", seriesName)
	series, err := this.dbc.Query(dbQuery, "s")

	if err != nil {
		return nil, err
	}
	return series, nil
}
