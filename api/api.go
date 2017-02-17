// Package api enables to specify the connection parameter and to create corresponding ingest messages and queries.
package api

import (
	"fmt"
	"net/http"

	"github.com/oudsteyn/go-predix-timeseries/ingest"
	"github.com/oudsteyn/go-predix-timeseries/query"
	"github.com/gorilla/websocket"
	"gopkg.in/h2non/gentleman.v0"
)

// Stores information required to send ingest messages and to perform queries
type Api struct {
	ingestUrl    string
	queryUrl     string
	authToken    string
	predixZoneId string
	conn         *websocket.Conn
	client       *gentleman.Client
}

// Creates a new Api object capable to do both ingest and query requests
func New(ingestUrl, queryUrl, authToken, predixZoneId string) *Api {
	return &Api{
		ingestUrl:    ingestUrl,
		queryUrl:     queryUrl,
		authToken:    authToken,
		predixZoneId: predixZoneId,
	}
}

// Creates Api objects capable to send ingest messages only
func Ingest(ingestUrl, authToken, predixZoneId string) *Api {
	return New(ingestUrl, "", authToken, predixZoneId)
}

// Creates Api objects capable to do query requests only
func Query(queryUrl, authToken, predixZoneId string) *Api {
	return New("", queryUrl, authToken, predixZoneId)
}

// Creates an ingest message used to push data to Time Series
func (a *Api) IngestMessage() *ingest.Message {
	if a.conn == nil {
		if a.ingestUrl == "" {
			return nil
		}
		h := http.Header{}
		h.Add("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
		h.Add("Predix-Zone-Id", a.predixZoneId)
		h.Add("Origin", "http://localhost")
		conn, _, err := websocket.DefaultDialer.Dial(a.ingestUrl, h)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil
		}
		a.conn = conn
	}
	return ingest.NewMessage(a.conn)
}

// Creates a query to do a read request to Time Series
func (a *Api) Query() query.Query {
	if a.client == nil {
		if a.queryUrl == "" {
			return nil
		}
		a.client = gentleman.New()
		a.client.BaseURL(a.queryUrl)
		a.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
		a.client.SetHeader("Predix-Zone-Id", a.predixZoneId)
	}
	return query.New(a.client)
}
