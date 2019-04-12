"use strict";

function eventDebug(name, data) {
    console.log(`WS API :: EVENT < ${name} > ::`, data);
}

function addButton(name) {
    var btn = document.createElement('button');
    btn.innerText = name;
    btn.id = `soundBtn-${name}`;
    btn.className = 'btn btn-primary btn-sound m-2';
    btn.onclick = (e) => 
        ws.emit('PLAY', {
            'ident': name,
            'source': 0,
        });
    $('#container-sound-btns').append(btn);
};

function fetchSoundsList(sort, cb) {
    $('#container-sound-btns').empty();
    $('#spinnerLoadingSounds').addClass('d-flex');
    $('#spinnerLoadingSounds').removeClass('d-none');
    getLocalSounds(sort ? sort : 'NAME').then((sounds) => {
        sounds.forEach((s) => addButton(s.name));
        $('#spinnerLoadingSounds').removeClass('d-flex');
        $('#spinnerLoadingSounds').addClass('d-none');
        if (cb) cb();
    }).catch((r, s) => {
        console.log('REST :: ERROR :: ', r, s);
        displayError(`<code>REST API ERROR</code> getting sounds failed:<br/><code>${r}</code>`);
        $('#spinnerLoadingSounds').removeClass('d-flex');
        $('#spinnerLoadingSounds').addClass('d-none');
        if (cb) cb();
    });
}

function displayError(desc, time) {
    if (!time) time = 8000;

    var alertBox = $('#errorAlert')[0];
    $('#errorAlertText')[0].innerHTML = desc;
    alertBox.style.display = 'block';
    setTimeout(() => {
        alertBox.style.opacity = '1';
    }, 10);
    setTimeout(() => {
        alertBox.style.opacity = '0';
    }, time);
    setTimeout(() => {
        alertBox.style.display = 'none';
    }, time + 250);
}

ws.onEmit((e, raw) => console.log(`WS API :: COMMAND < ${e.name} > ::`, e.data));

// --------------------------
// --- INIT

var sortBy = getCookieValue('sort_by');
var inChannel = false;

if (getCookieValue('cookies_accepted') !== '1') {
    $('#cookieInformation')[0].style.display = 'block';
}

$('#btnSortBy').on('click', (e) => {
    sortBy = sortBy == 'DATE' ? 'NAME' : 'DATE'; 
    document.cookie = `sort_by=${sortBy}; paht=/`;
    $('#btnSortBy')[0].innerText = 'SORT BY ' + (sortBy == 'DATE' ? 'NAME' : 'DATE');
    fetchSoundsList(sortBy);
});

$('#btCookieAccept').on('click', (e) => {
    document.cookie = 'cookies_accepted=1; path=/';
    $('#cookieInformation')[0].style.display = 'none';
});

$('#btCookieDecline').on('click', (e) => {
    deleteAllCookies();
    window.location = '/static/cookies-declined.html';
});

$('#btnStop').on('click', (e) => {
    ws.emit('STOP');
});

$('#btnJoinLeave').on('click', (e) => {
    if (inChannel)
        ws.emit('LEAVE');
    else
        ws.emit('JOIN');
});

$('#btnLog').on('click', (e) => {
    alert('Gibbet noch ned!');
});

$('#btnStats').on('click', (e) => {
    alert('Gibbet noch ned!');
});

if (sortBy)
    $('#btnSortBy')[0].innerText = 'SORT BY ' + (sortBy == 'DATE' ? 'NAME' : 'DATE');

fetchSoundsList(sortBy);

// --------------------------
// --- WS EVENT HANDLERS

ws.on('ERROR', (data) => {
    eventDebug('ERROR', data);
    displayError(`<code>${data.data.code} - ${data.data.type}</code>&nbsp; ${data.data.message}`);
});

ws.on('HELLO', (data) => {
    eventDebug('HELLO', data);
});

ws.on('PLAYING', (data) => {
    eventDebug('PLAYING', data);
    if (data.data.ident) {
        $(`#soundBtn-${data.data.ident}`).addClass('playing');
    }
    inChannel = true;
    $('#btnJoinLeave')[0].innerText = 'LEFT';
});

ws.on('END', (data) => {
    eventDebug('END', data);
    if (data.data.ident) {
        $(`#soundBtn-${data.data.ident}`).removeClass('playing');
    }
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
    inChannel = true;
    $('#btnJoinLeave')[0].innerText = 'LEFT';
});

ws.on('LEFT', (data) => {
    eventDebug('LEFT', data);
    inChannel = false;
    $('#btnJoinLeave')[0].innerText = 'JOIN';
});