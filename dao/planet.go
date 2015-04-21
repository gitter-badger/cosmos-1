package dao

import (
	"fmt"

	"github.com/cosmos-io/influxdbc"
)

type PlanetDao struct {
	dbc *influxdbc.InfluxDB
}

func (this *PlanetDao) GetPlanet(token, planet string, lifeTime int) ([]*influxdbc.Series, error) {
	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.PLANET\\.%s\\.%s\\.Name/ WHERE time > now() - %dm LIMIT 1", token, planet, lifeTime), "s")
}

func (this *PlanetDao) GetPlanets(token string, lifeTime int) ([]*influxdbc.Series, error) {
	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^MIN\\.PLANET\\.%s/ WHERE time > now() - %dm LIMIT 1", token, lifeTime), "s")
}

// func (this *PlanetDao) GetPlanetNamesPassLifeTimeOfUser(token string, lifeTime int) ([]*influxdbc.Series, error) {
// 	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^PLANET\\.%s\\.\\w+\\.Name$/ WHERE time < now() - %dm AND time > now() - %dm LIMIT 1", token, lifeTime, lifeTime+5), "s")
// }

func (this *PlanetDao) GetPlanetStatusesPassLifeTime(lifeTime int) ([]*influxdbc.Series, error) {
	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^PLANET\\.\\w+\\.\\w+\\.Status$/ WHERE time < now() - %dm AND time > now() - %dm LIMIT 1", lifeTime, lifeTime+5), "s")
}

func (this *PlanetDao) GetPlanetStatusesInLifeTime(lifeTime int) ([]*influxdbc.Series, error) {
	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^PLANET\\.\\w+\\.\\w+\\.Status$/ WHERE time > now() - %dm LIMIT 1", lifeTime), "s")
}

func (this *PlanetDao) GetPlanetStatusesInLifeTimeOfUser(token string, lifeTime int) ([]*influxdbc.Series, error) {
	return this.dbc.Query(fmt.Sprintf("SELECT num_value, txt_value FROM /^PLANET\\.%s\\.\\w+\\.Status$/ WHERE time > now() - %dm LIMIT 1", token, lifeTime), "s")
}
