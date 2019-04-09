package wsmgr

import "encoding/json"

const (
	eventInit    = "INIT"
	eventClosing = "CLOSING"
)

type Event struct {
	Sender *WebSocketConn `json:"-"`

	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func NewEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

func (e *Event) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) ParseDataTo(v interface{}) error {
	bData, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bData, v)
}
