package influxdb

import (
	"fmt"

	"github.com/hiennguyen9874/stockk-go/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	addr := fmt.Sprintf("http://%v:%v", cfg.InfluxDB.Host, cfg.InfluxDB.Port)

	client := influxdb2.NewClient(addr, cfg.InfluxDB.Token)

	return client, nil
}
