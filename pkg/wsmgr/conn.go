package wsmgr

import (
	"github.com/gorilla/websocket"
)

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

func (wsc *WebSocketConn) Out(e *Event) error {
	data, err := e.ToJson()
	if err != nil {
		return err
	}

	// go func() {
	wsc.out <- data
	// }()

	return nil
}

func (wsc *WebSocketConn) SetIdent(ident interface{}) {
	wsc.ident = ident
}

func (wsc *WebSocketConn) GetIdent() interface{} {
	return wsc.ident
}

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

func (wsc *WebSocketConn) Close() {
	wsc.close <- true
}
