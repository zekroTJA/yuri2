"use strict";

function eventDebug(name, data) {
    console.log(`WS API :: EVENT < ${name} > ::`, data);
}

function addButton(name) {
    var btn = document.createElement('button');
    btn.innerText = name;
    btn.id = 'sound-btn-name';
    btn.className = 'btn btn-primary m-2';
    btn.onclick = (e) => 
        ws.emit('PLAY', {
            'ident': name,
            'source': 0,
        });
    $('#container-sound-btns').append(btn);
};

ws.onEmit((e, raw) => console.log(`WS API :: COMMAND < ${e.name} > ::`, e.data));

// --------------------------
// --- INIT BUTTONS

getLocalSounds().then((sounds) => {
    sounds.forEach((s) => addButton(s.name));
}).catch((r, s) => {
    console.log('REST :: ERROR :: ', r, s);
});

// --------------------------
// --- WS EVENT HANDLERS

ws.on('ERROR', (data) => {
    eventDebug('ERROR', data);
});

ws.on('HELLO', (data) => {
    eventDebug('HELLO', data);
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