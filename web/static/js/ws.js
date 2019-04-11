function Handler(id, cb) {
    this.cb = cb;
    this.id = id;
    
    this.call = (...args) => 
        cb.apply(this, args);
}

function WsClient(url) {
    this.ws = new WebSocket(url);

    this._rollingID = 0;
    this.eventListener = {};

    this.onEmit = (cb) => this._onEmit = cb;

    this.on = (event, cb) => {
        if (!this.eventListener[event])
            this.eventListener[event] = [];

        var id = this._rollingID++;
        this.eventListener[event].push(
            new Handler(id, cb));

        return () => {
            if (this.eventListener[event]) {
                var i = this.eventListener[event]
                    .findIndex((h) => h.id == id);
                this.eventListener[event].splice(i, 1);
            }
        };
    }

    this.emit = (name, data) => {
        let event = {
            name: name, 
            data: data,
        }
        let rawData = JSON.stringify(event);

        if (this._onEmit) this._onEmit(event, rawData);

        this.ws.send(rawData);
    }

    this.ws.onmessage = (response) => {
        try {
            let data = JSON.parse(response.data);
            if (data) {
                let cbs = this.eventListener[data.name]
                if (cbs)
                    cbs.forEach((h) => h.call(data));
            }
        } catch (e) {
            console.log(e)
        }
    }

    this.onOpen = (handler) => this.ws.onopen = handler;
    this.onClose = (handler) => this.ws.onclose = handler;
}

var ws = new WsClient(
    window.location.href.replace(/((http)|(https)):\/\//gm, 
        window.location.href.startsWith('http://') ? 'ws:/' : 'wss://') + 'ws'
);

// --------------------------------------------------------------------------------------

ws.onOpen(() => {
    ws.emit('INIT', {
        user_id: getCookieValue('userid'),
        token: getCookieValue('token'),
    });
});

// --------------------------------------------------------------------------------------

function getCookieValue(name) {
    var c = document.cookie
        .split(";")
        .map((c) => c.trim().split('='))
        .find((c) => c[0] == name);
    if (c)
        return c[1];
}