package model

import "time"

type Counter struct {
	Routing         string    `json:"routing" bson:"routing"`
	Count           int       `json:"count" bson:"count"`
	MinExecTime     int64     `json:"min_exec_time" bson:"min_exec_time"`
	MaxExecTime     int64     `json:"max_exec_time" bson:"max_exec_time"`
	MedianExecTime  int64     `json:"median_exec_time" bson:"median_exec_time"`
	AverageExecTime int       `json:"average_exec_time" bson:"average_exec_time"`
	AddedTime       time.Time `json:"added_time" bson:"added_time"`
	UpdateTime      time.Time `json:"update_time" bson:"update_time"`
}

func (Counter) Collection() string {
	return "counter_statistics"
}
