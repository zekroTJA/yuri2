package wsmgr

import (
	"net/http"
	"sync"

	"github.com/zekroTJA/yuri2/pkg/multierror"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// An OnErrorFunc handles an error event inside
// the web socket message handling.
type OnErrorFunc func(info string, err error)

// An OnEventFunc handles an server side ingoing
// web socket event.
type OnEventFunc func(e *Event)

// wsError contains error
// infotmation.
type wsError struct {
	info string
	err  error
}

// WebSocketManager maintains multiple web socket
// connections.
type WebSocketManager struct {
	mx *sync.Mutex

	onError OnErrorFunc

	conns  map[*WebSocketConn]interface{}
	events map[string]OnEventFunc

	errPipe     chan *wsError
	eventPipe   chan *Event
	closingPipe chan *WebSocketConn
}

func New() *WebSocketManager {
	wsm := &WebSocketManager{
		mx:          new(sync.Mutex),
		conns:       make(map[*WebSocketConn]interface{}),
		events:      make(map[string]OnEventFunc),
		errPipe:     make(chan *wsError),
		eventPipe:   make(chan *Event),
		closingPipe: make(chan *WebSocketConn),
	}

	go wsm.pipeListener()

	return wsm
}

func (wsm *WebSocketManager) OnError(onError OnErrorFunc) {
	wsm.onError = onError
}

func (wsm *WebSocketManager) On(name string, handler OnEventFunc) {
	wsm.events[name] = handler
}

func (wsm *WebSocketManager) NewConn(w http.ResponseWriter, r *http.Request, ident interface{}) (*WebSocketConn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	wsc := NewWebSocketConn(conn, ident, wsm.errPipe, wsm.eventPipe, wsm.closingPipe)

	wsm.mx.Lock()
	wsm.conns[wsc] = nil
	wsm.mx.Unlock()

	return wsc, nil
}

func (wsm *WebSocketManager) Broadcast(e *Event, exclude ...*WebSocketConn) error {
	mErr := multierror.New(nil)
	var err error

	for wsc := range wsm.conns {
		for _, ex := range exclude {
			if ex == wsc {
				continue
			}
			err = wsc.Out(e)
			mErr.Append(err)
		}
	}

	return mErr.Concat()
}

func (wsm *WebSocketManager) pipeListener() {

	for {
		select {

		// on event
		case event, ok := <-wsm.eventPipe:
			if ok {
				wsm.mx.Lock()
				handler, ok := wsm.events[event.Name]
				wsm.mx.Unlock()

				if ok && handler != nil {
					handler(event)
				}
			}

		// on error
		case err, ok := <-wsm.errPipe:
			if ok && wsm.onError != nil {
				wsm.onError(err.info, err.err)
			}

		// on ws connection closing
		case wsc, ok := <-wsm.closingPipe:
			if ok && wsc != nil {
				wsm.mx.Lock()
				delete(wsm.conns, wsc)
				wsm.mx.Unlock()
			}
		}
	}

}
