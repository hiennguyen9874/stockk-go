package influxdb

import (
	"fmt"

	"github.com/hiennguyen9874/stockk-go/config"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	addr := fmt.Sprintf("http://%v:%v", cfg.InfluxDB.InfluxDBHost, cfg.InfluxDB.InfluxDBPort)

	client := influxdb2.NewClient(addr, cfg.InfluxDB.InfluxDBToken)

	return client, nil
}
