package models

import "time"

type Bar struct {
	Symbol   string
	Exchange string
	Time     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}
