// Package ingest provides primitives to push data to time series.
package ingest

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oudsteyn/go-predix-timeseries/tag"
	"github.com/gorilla/websocket"
)

type Message struct {
	id   int64
	tags map[string]*tag.Tag
	conn *websocket.Conn
}

type acknowledgementMessage struct {
	messageId  int64 `json:"messageId"`
	statusCode int   `json:"statusCode"`
}

func NewMessage(conn *websocket.Conn) *Message {
	return &Message{time.Now().UnixNano() / int64(time.Millisecond), make(map[string]*tag.Tag), conn}
}

func (m *Message) MarshalJSON() ([]byte, error) {
	body := []tag.Tag{}
	for _, tag := range m.tags {
		body = append(body, *tag)
	}
	return json.Marshal(struct {
		Id   int64     `json:"messageId"`
		Body []tag.Tag `json:"body"`
	}{
		Id:   m.id,
		Body: body,
	})

}

// Adds a tag to the message
func (m *Message) AddTag(name string) (*tag.Tag, error) {
	if tag.CorrectNameRE.MatchString(name) {
		if _, ok := m.tags[name]; !ok {
			m.tags[name] = tag.New(name)
		}

		return m.tags[name], nil
	} else {
		return nil, tag.IncorrectName
	}
}

// Gets a tag from the message
func (m *Message) GetTag(name string) (*tag.Tag, bool) {
	t, ok := m.tags[name]
	return t, ok
}

// Deletes a tag from the message
func (m *Message) DeleteTag(name string) {
	delete(m.tags, name)
}

// Sends the message to Time Series, i.e. push datapoints
func (m *Message) Send() error {
	e := m.conn.WriteJSON(m)
	if e != nil {
		return e
	}
	_, p, e := m.conn.ReadMessage()
	if e != nil {
		return e
	}
	var resp acknowledgementMessage
	e = json.Unmarshal(p, &resp)
	if e == nil {
		switch resp.statusCode {
		case 200:
			e = nil
		case 400:
			e = errors.New("Bad request")
		case 401:
			e = errors.New("Unauthorized")
		case 413:
			e = errors.New("Request entity too large")
		case 503:
			e = errors.New("Failed to ingest data")
		}
	}
	return e
}
