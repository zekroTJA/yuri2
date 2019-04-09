function WsClient(url) {
    this.ws = new WebSocket(url);

    this.eventListener = {};

    this.on = (event, cb) => this.eventListener[event] = cb;

    this.emit = (name, data) => {
        let event = {
            name: name, 
            data: data,
        }
        let rawData = JSON.stringify(event);
        console.log("sending: ", rawData);
        this.ws.send(rawData);
    }

    this.ws.onmessage = (response) => {
        try {
            let data = JSON.parse(response.data);
            if (data) {
                let cb = this.eventListener[data.name]
                if (cb)
                    cb(data.data)
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

ws.on('ERROR', (data) => {
    console.log("ERROR :: ", data)
});

ws.onOpen(() => {
    ws.emit('INIT', {
        user_id: getCookieValue('userid'),
        token: getCookieValue('token'),
    });

    setTimeout(() => {
        ws.emit('PLAY', 'test123')
    }, 1000);
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