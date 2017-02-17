package query

import (
	"github.com/oudsteyn/go-predix-timeseries/datapoint"
)

type Group map[string]interface{}
type Filter map[string]interface{}
type Attribute map[string]interface{}

type Result struct {
	Tags []struct {
		Name    string `json:"name"`
		Results []struct {
			Groups     []Group               `json:"groups"`
			Filters    Filter                `json:"filters"`
			Values     []datapoint.Datapoint `json:"values"`
			Attributes Attribute             `json:"attributes"`
		} `json:"results"`
		Stats struct {
			RawCount int `json:"rawCount"`
		} `json:"stats"`
	} `json:"tags"`
}
