package worker

import (
	"fmt"
	"time"

	"github.com/cosmos-io/cosmos/influxdb"
)

type NewsFeedWorker struct {
	db        *influxdb.InfluxDB
	delayTime time.Duration
}

var (
	lastCheckTime time.Time
)

func NewNewsFeedWorker(db *influxdb.InfluxDB, delayTime time.Duration) *NewsFeedWorker {
	return &NewsFeedWorker{db: db, delayTime: delayTime}
}

func (this *NewsFeedWorker) Run() {
	lastCheckTime = time.Now()

	go func() {
		ticker := time.NewTicker(time.Millisecond * this.delayTime)
		for _ = range ticker.C {
			this.checkNewContainer()
			this.checkContainerUpDown()
			lastCheckTime = time.Now()
		}
	}()
}

func (this *NewsFeedWorker) checkNewContainer() error {
	result, err := this.db.QueryFirstContainerMetrics()
	if err != nil {
		return err
	}
	for _, v := range result {
		time, err := time.Parse(time.RFC3339Nano, v["time"])
		if err == nil {
			if time.Unix() > lastCheckTime.Unix() {
				// New Container
			}
		}
	}

	return nil
}

func (this *NewsFeedWorker) checkContainerUpDown() error {
	start := fmt.Sprintf("%dm", (this.delayTime/1000/60)*2)
	end := fmt.Sprintf("%dm", (this.delayTime / 1000 / 60))

	containersPast, err := this.db.QueryContainersInRange(start, end, end)
	if err != nil {
		return err
	}

	containersCurrent, err := this.db.QueryContainersInRange(end, "0m", end)
	if err != nil {
		return err
	}

	fmt.Println(containersPast)
	fmt.Println(containersCurrent)

	for c, _ := range containersPast {
		if _, ok := containersCurrent[c]; ok == false {
		}
	}

	for c, _ := range containersCurrent {
		if _, ok := containersPast[c]; ok == false {
		}
	}

	return nil
}
