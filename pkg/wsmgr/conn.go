package wsmgr

import (
	"github.com/gorilla/websocket"
)

// A WebSocketConn maintians the communication
// and data transfer of a WS connection.
type WebSocketConn struct {
	conn  *websocket.Conn
	ident interface{}

	out chan []byte
	in  chan []byte

	errPipe     chan *wsError
	eventPipe   chan *Event
	closingPipe chan *WebSocketConn

	close chan bool
}

// NewWebSocketConn creates a new WebSocketConnection
// by upgrading an incomming request. You may pass an
// ident to identify the sender of an event in event
// handlers.
func NewWebSocketConn(conn *websocket.Conn, ident interface{},
	errPipe chan *wsError, eventPipe chan *Event, closingPipe chan *WebSocketConn) *WebSocketConn {

	wsc := &WebSocketConn{
		conn:        conn,
		ident:       ident,
		out:         make(chan []byte),
		in:          make(chan []byte),
		errPipe:     errPipe,
		eventPipe:   eventPipe,
		closingPipe: closingPipe,
		close:       make(chan bool, 1),
	}

	go wsc.reader()
	go wsc.writer()

	return wsc
}

// Out sends an event to the connected client.
// This function blocks until the message was
// send to the client.
func (wsc *WebSocketConn) Out(e *Event) error {
	data, err := e.ToJSON()
	if err != nil {
		return err
	}

	wsc.out <- data

	return nil
}

// SetIdent sets the ident of the connection.
func (wsc *WebSocketConn) SetIdent(ident interface{}) {
	wsc.ident = ident
}

// GetIdent returns the ident value of
// the connection.
func (wsc *WebSocketConn) GetIdent() interface{} {
	return wsc.ident
}

// reader starts a blocking loop waiting for
// incomming messages, which will be tried to be
// JSON-parsed into Event objects, which will then
// be send to the eventPipe. If something fails,
// this will be send into the errPipe as wsError
// object.
func (wsc *WebSocketConn) reader() {
	defer func() {
		wsc.conn.Close()
		wsc.closingPipe <- wsc
	}()

	for {
		event := new(Event)
		if err := wsc.conn.ReadJSON(event); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				break
			}

			wsc.errPipe <- &wsError{
				info: "INCOMMING :: ReadJSON",
				err:  err,
			}
		}

		if event != nil {
			event.Sender = wsc
			wsc.eventPipe <- event
		}
	}
}

// writer starts a blocking loop waiting for
// outgoing messages in the out channel or
// for closing the connection by signal
// into close channel.
func (wsc *WebSocketConn) writer() {
	for {
		select {

		case msg, ok := <-wsc.out:
			if !ok {
				wsc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			err := wsc.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				wsc.errPipe <- &wsError{
					info: "OUTGOING :: NextWriter",
					err:  err,
				}
			}

		case <-wsc.close:
			wsc.conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}
	}
}

// Close sends a close signal into
// the close channel so the connection
// will be closed.
func (wsc *WebSocketConn) Close() {
	wsc.close <- true
}
