package wsmgr

import "encoding/json"

// An Event contains the Name and the
// Data of a web socket event and may
// contain the sender WebSocketConn
// pointer.
type Event struct {
	Sender *WebSocketConn `json:"-"`

	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// NewEvent creates a new instance of Event
// with the passed name and data.
func NewEvent(name string, data interface{}) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
}

// ToJSON tries to parse the Event object
// to a JSON formatted byte array.
func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ParseDataTo tries to parse the data
// interface to the passed interface
// object by marshaling the data object
// to JSON and then unmarshaling the raw
// JSON data to the passed object.
func (e *Event) ParseDataTo(v interface{}) error {
	bData, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bData, v)
}
