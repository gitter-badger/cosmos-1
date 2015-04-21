package worker

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos-io/cosmos/dao"
	"github.com/cosmos-io/cosmos/service"
	"github.com/cosmos-io/influxdbc"
)

type NewsFeedWorker struct {
	newsFeedService *service.NewsFeedService
	lifeTime        int
	delayTime       time.Duration
}

func NewNewsFeedWorker(lifeTime int, delayTime time.Duration) *NewsFeedWorker {
	return &NewsFeedWorker{newsFeedService: &service.NewsFeedService{}, lifeTime: lifeTime, delayTime: delayTime}
}

func (this *NewsFeedWorker) Run() {
	go func() {
		ticker := time.NewTicker(time.Second * this.delayTime)
		for _ = range ticker.C {
			planetSeries, err := this.checkRemovedPlanets()
			if err != nil {
				fmt.Println(err)
				continue
			}
			containerSeries, err := this.checkRemovedContainers()
			if err != nil {
				fmt.Println(err)
				continue
			}
			series := make([]*influxdbc.Series, 0)
			series = append(series, planetSeries...)
			series = append(series, containerSeries...)

			_, err = dao.Series.WriteSeries(series)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}()
}

func (this *NewsFeedWorker) checkRemovedPlanets() ([]*influxdbc.Series, error) {
	result := make([]*influxdbc.Series, 0)

	passed, err := dao.Planet.GetPlanetStatusesPassLifeTime(this.lifeTime)
	if err != nil {
		return nil, err
	}
	current, err := dao.Planet.GetPlanetStatusesInLifeTime(this.lifeTime)
	if err != nil {
		return nil, err
	}

	for _, passedCont := range passed {
		exist := false
		for _, currentCont := range current {
			if passedCont.Name == currentCont.Name {
				exist = true
				break
			}
		}

		if exist == false {
			// Removed one
			comps := strings.Split(passedCont.Name, ".")
			token := comps[1]
			planet := comps[2]
			series, err := this.newsFeedService.MakeNewsFeed(token, planet, service.FEED_TYPE_REMOVE_PLANET)
			if err != nil {
				fmt.Println(err)
				continue
			}
			result = append(result, series)

			// make one explicitly
			statusSeries := influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Status", token, planet), "txt_value", "num_value")
			statusSeries.AddPoint("Exited", nil)
			result = append(result, statusSeries)

			statusSeries = influxdbc.NewSeries(fmt.Sprintf("PLANET.%s.%s.Name", token, planet), "txt_value", "num_value")
			statusSeries.AddPoint(planet, nil)
			result = append(result, statusSeries)
		}
	}

	return result, nil
}

func (this *NewsFeedWorker) checkRemovedContainers() ([]*influxdbc.Series, error) {
	return nil, nil
}
