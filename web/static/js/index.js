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

function addRandomButton() {
    var btn = document.createElement('button');
    btn.innerText = '🎲';
    btn.id = `soundBtn-RANDOM`;
    btn.className = 'btn btn-primary m-2';
    btn.onclick = (e) => 
        ws.emit('RANDOM');
    $('#container-sound-btns').append(btn);
}

function fetchSoundsList(sort, cb) {
    var spinner = $('#spinnerLoadingSounds');
    $('#container-sound-btns').empty();
    spinner.addClass('d-flex');
    spinner.removeClass('d-none');
    getLocalSounds(sort ? sort : 'NAME').then((sounds) => {
        addRandomButton();
        sounds.forEach((s) => addButton(s.name));
        spinner.removeClass('d-flex');
        spinner.addClass('d-none');
        if (cb) cb();
    }).catch((r, s) => {
        console.log('REST :: ERROR :: ', r, s);
        displayError(`<code>REST API ERROR</code> getting sounds failed:<br/><code>${r}</code>`);
        spinner.removeClass('d-flex');
        spinner.addClass('d-none');
        if (cb) cb();
    });
}

function displayError(desc, time) {
    if (!time) time = 8000;

    var alertBox = $('#errorAlert')[0];
    $('#errorAlertText')[0].innerHTML = desc;

    // fade in
    alertBox.style.display = 'block';
    setTimeout(() => {
        alertBox.style.opacity = '1';
        alertBox.style.transform = 'translateY(0px)';
    }, 10);
    // fade out
    setTimeout(() => {
        alertBox.style.opacity = '0';
        alertBox.style.transform = 'translateY(-20px)';
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
var guildID = null;

if (getCookieValue('cookies_accepted') !== '1') {
    $('#cookieInformation')[0].style.display = 'block';
}

// BUTTON EVENT HOOKS

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
    getGuildLog(guildID).then((res) => {
        $('#modalLog').modal('show');
        var tab = $('#modalLog div.modal-body > table');

        Array.forEach(tab.children('tr'), (tr) => tr.remove());
        res.forEach((r) => {
            var tr = document.createElement('tr');

            var tdTime = document.createElement('td');
            tdTime.innerText = getTime(new Date(r.time));
            tr.appendChild(tdTime);

            var tdCaller = document.createElement('td');
            tdCaller.innerText = r.user_tag;
            tr.appendChild(tdCaller);

            var tdSound = document.createElement('td');
            tdSound.innerText = r.sound;
            tr.appendChild(tdSound);

            var tdSource = document.createElement('td');
            switch (r.source) {
                case 'local':
                    tdSource.innerHTML = '<span class="badge badge-primary">L</span>';
                    break;
                case 'youtube':
                    tdSource.innerHTML = '<span class="badge badge-warning">Y</span>';
                    break;
                case 'http':
                    tdSource.innerHTML = '<span class="badge badge-info">H</span>';
                    break;
                default:
                    tdSource.innerHTML = '<span class="badge badge-dark">?</span>';
            }
            tr.appendChild(tdSource);

            tab.append(tr);
        });

    }).catch((r, s) => {
        console.log('REST :: ERROR :: ', r, s);
        displayError(`<code>REST API ERROR</code> getting log failed: You need to be in a voice channel to open the guilds log.`);
    });
});

$('#btnStats').on('click', (e) => {
    getGuildStats(guildID).then((res) => {
        $('#modalStats').modal('show');
        var tab = $('#modalStats div.modal-body > table');

        Array.forEach(tab.children('tr'), (tr) => tr.remove());
        res.forEach((r, i) => {
            var tr = document.createElement('tr');

            var tdNumber = document.createElement('td');
            tdNumber.innerHTML = `<span class="badge badge-primary">${i + 1}</span>`;
            tr.appendChild(tdNumber);

            var tdSound = document.createElement('td');
            tdSound.innerText = r.sound;
            tr.appendChild(tdSound);

            var tdCount = document.createElement('td');
            tdCount.innerText = r.count;
            tr.appendChild(tdCount);

            tab.append(tr);
        });

    }).catch((r, s) => {
        console.log('REST :: ERROR :: ', r, s);
        displayError(`<code>REST API ERROR</code> getting log failed: You need to be in a voice channel to open the guilds stats.`);
    });
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
    if (data.data && data.data.connected) {
        $('#btnJoinLeave')[0].innerText = 'LEAVE';
        inChannel = true;
        guildID = data.data.voice_state.guild_id;
    }
});

ws.on('PLAYING', (data) => {
    eventDebug('PLAYING', data);
    if (data.data.ident) {
        $(`#soundBtn-${data.data.ident}`).addClass('playing');
    }
    inChannel = true;
    $('#btnJoinLeave')[0].innerText = 'LEAVE';
    guildID = data.data.guild_id;
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
    $('#btnJoinLeave')[0].innerText = 'LEAVE';
    guildID = data.data.guild_id;
});

ws.on('LEFT', (data) => {
    eventDebug('LEFT', data);
    inChannel = false;
    $('#btnJoinLeave')[0].innerText = 'JOIN';
});