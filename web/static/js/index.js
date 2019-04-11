
function eventDebug(name, data) {
    console.log(`WSAPI :: EVENT < ${name} > ::`, data);
}

ws.onEmit((e, raw) => console.log(`WSAPI :: COMMAND < ${e.name} > ::`, e.data));


ws.on('ERROR', (data) => {
    eventDebug('ERROR', data);
});

ws.on('PLAYING', (data) => {
    eventDebug('PLAYING', data);
});

ws.on('END', (data) => {
    eventDebug('END', data);
});

ws.on('PLAY_ERROR', (data) => {
    eventDebug('PLAY_ERROR', data);
});

ws.on('STUCK', (data) => {
    eventDebug('STUCK', data);
});

ws.on('VOLUME_CHANGED', (data) => {
    eventDebug('VOLUME_CHANGED', data);
});

ws.on('JOINED', (data) => {
    eventDebug('JOINED', data);
});

ws.on('LEFT', (data) => {
    eventDebug('LEFT', data);
});