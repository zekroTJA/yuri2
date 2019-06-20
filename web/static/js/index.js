/** @format */

'use strict';

var chain = null;
var favorites = [];
var test;

// ------------------------------------------------------------

function eventDebug(name, data) {
    console.log(`WS API :: EVENT < ${name} > ::`, data);
}

function addButton(name) {
    var btn = document.createElement('button');
    btn.innerText = name;
    btn.id = `soundBtn-${name}`;
    btn.className = 'btn btn-primary btn-sound m-2';
    btn.onclick = (e) => {
        if (e.target.className.includes('btn-favorite')) return;

        if (chain !== null) {
            var t = $(e.target);
            if (t.hasClass('sound-enqueued')) {
                t.removeClass('sound-enqueued');
                delete chain[name];
                return;
            }
            chain[name] = t;
            t.addClass('sound-enqueued');
            return;
        }

        ws.emit('PLAY', {
            ident: name,
            source: 0,
        });
    };

    var fav = document.createElement('a');
    fav.className = 'btn-favorite';
    if (favorites.includes(name)) fav.className += ' d-block';
    fav.onclick = (e) => {
        if (favorites.includes(name)) {
            deleteFavorites(name)
                .then(() => {
                    test = $(e.target).find('.btn-favorite');
                    favorites.splice(favorites.indexOf(name), 1);
                    $(e.target)
                        .find('.btn-favorite')
                        .prevObject.removeClass('d-block');
                })
                .catch((err) => {
                    displayError(
                        `<code>REST API ERROR</code> unsetting favorite failed:<br/><code>${err}</code>`
                    );
                });
        } else {
            postFavorites(name)
                .then(() => {
                    test = $(e.target).find('.btn-favorite');
                    $(e.target)
                        .find('.btn-favorite')
                        .prevObject.addClass('d-block');
                    favorites.push(name);
                })
                .catch((err) => {
                    displayError(
                        `<code>REST API ERROR</code> setting favorite failed:<br/><code>${err}</code>`
                    );
                });
        }
    };

    btn.appendChild(fav);
    $('#container-sound-btns').append(btn);
}

function addRandomButton() {
    var btn = document.createElement('button');
    btn.innerText = '🎲';
    btn.id = `soundBtn-RANDOM`;
    btn.className = 'btn btn-primary m-2';
    btn.onclick = (e) => ws.emit('RANDOM');
    $('#container-sound-btns').append(btn);
}

function fetchSoundsList(sort, cb) {
    var spinner = $('#spinnerLoadingSounds');
    spinner.addClass('d-flex');
    spinner.removeClass('d-none');
    $('#container-sound-btns').empty();
    getLocalSounds(sort ? sort : 'NAME')
        .then((sounds) => {
            updateSoundList(sounds);
            updateFastTriggerSelector(sounds);
            spinner.removeClass('d-flex');
            spinner.addClass('d-none');
            if (cb) cb(sounds);
        })
        .catch((r, s) => {
            console.log('REST :: ERROR :: ', r, s);
            displayError(
                `<code>REST API ERROR</code> getting sounds failed:<br/><code>${r}</code>`
            );
            spinner.removeClass('d-flex');
            spinner.addClass('d-none');
            if (cb) cb();
        });
}

function updateSoundList(soundList) {
    $('#container-sound-btns').empty();
    addRandomButton();
    soundList
        .filter((s) => favorites.includes(s.name))
        .forEach((s) => {
            addButton(s.name);
        });
    soundList
        .filter((s) => !favorites.includes(s.name))
        .forEach((s) => {
            addButton(s.name);
        });
}

function updateFastTriggerSelector(soundList) {
    var ddff = $('#ddFastTrigger');
    ddff.children().empty();

    var opt = document.createElement('option');
    opt.innerText = 'RANDOM';
    ddff.append(opt);

    soundList.forEach((s) => {
        var opt = document.createElement('option');
        opt.innerText = s.name;
        ddff.append(opt);
    });

    getFastTrigger()
        .then((res) => {
            if (res.random) {
                ddff.val('RANDOM');
            } else {
                ddff.val(res.ident);
            }
        })
        .catch((r, s) => {
            console.log('REST :: ERROR :: ', r, s);
            displayError(
                `<code>REST API ERROR</code> getting fast trigger:<br/><code>${r}</code>`
            );
        });
}

function filterSoundsList(query, sl) {
    if (!query) updateSoundList(sl);
    query = query.toLowerCase();
    updateSoundList(
        sl.filter((s) => {
            if (query.startsWith('*')) return s.name.endsWith(query.substr(1));
            if (query.endsWith('*'))
                return s.name.startsWith(query.substr(0, query.length - 1));
            return s.name.includes(query);
        })
    );
}

function setVolume(v) {
    var container = $('#containerVol');
    if (container.hasClass('invis')) container.removeClass('invis');
    $('#sliderVol').val(v);
    $('#labelVol')[0].innerText = v + '%';
}

function resetChaining() {
    var t = $('#btnChaining');
    t.removeClass('btn-chaning-active');
    t.text('CHAINING');
    Object.keys(chain).forEach((s) => {
        chain[s].removeClass('sound-enqueued');
    });
    chain = null;
}

function playFromQueue() {
    if (Object.keys(chain).length === 0) {
        resetChaining();
        return false;
    }

    var sound = Object.keys(chain)[0];

    ws.emit('PLAY', {
        ident: sound,
        source: 0,
    });

    chain[sound].removeClass('sound-enqueued');
    delete chain[sound];

    return true;
}

ws.onEmit((e, raw) =>
    console.log(`WS API :: COMMAND < ${e.name} > ::`, e.data)
);

// --------------------------
// --- INIT

var sounds = [];
var sortBy = getCookieValue('sort_by');
var inChannel = false;
var guildID = null;

if (getCookieValue('cookies_accepted') !== '1') {
    $('#cookieInformation')[0].style.display = 'block';
}

getFavorites()
    .then((res) => {
        favorites = res.results;
        fetchSoundsList(sortBy, (s) => {
            if (s) sounds = s;
        });
    })
    .catch((err) => {
        console.log('REST :: ERROR :: ', err);
        displayError(
            `<code>REST API ERROR</code> getting favorites:<br/><code>${r}</code>`
        );
    });

// BUTTON EVENT HOOKS

$('#btnSortBy').on('click', (e) => {
    sortBy = sortBy == 'DATE' ? 'NAME' : 'DATE';
    document.cookie = `sort_by=${sortBy}; Max-Age=2147483647; Paht=/`;
    $('#btnSortBy')[0].innerText =
        'SORT BY ' + (sortBy == 'DATE' ? 'NAME' : 'DATE');
    fetchSoundsList(sortBy, (s) => {
        if (s) sounds = s;
    });
});

$('#btCookieAccept').on('click', (e) => {
    document.cookie = 'cookies_accepted=1; Max-Age=2147483647; Path=/';
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
    if (inChannel) ws.emit('LEAVE');
    else ws.emit('JOIN');
});

$('#btnChaining').on('click', (e) => {
    var t = $(e.target);

    if (chain === null) {
        chain = {};
        t.addClass('btn-chaning-active');
        t.text('RECORDING');
        return;
    }

    if (chain.length === 0) {
        t.removeClass('btn-chaning-active');
        t.text('CHAINING');
        chain = null;
        return;
    }

    t.text('PLAYING');
    playFromQueue();
});

$('#btnLog').on('click', (e) => {
    getGuildLog(guildID)
        .then((res) => {
            $('#modalLog').modal('show');
            var tab = $('#modalLog div.modal-body > table');

            tab.children('tr').each((_, tr) => tr.remove());
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
                        tdSource.innerHTML =
                            '<span class="badge badge-primary">L</span>';
                        break;
                    case 'youtube':
                        tdSource.innerHTML =
                            '<span class="badge badge-warning">Y</span>';
                        break;
                    case 'http':
                        tdSource.innerHTML =
                            '<span class="badge badge-info">H</span>';
                        break;
                    default:
                        tdSource.innerHTML =
                            '<span class="badge badge-dark">?</span>';
                }
                tr.appendChild(tdSource);

                tab.append(tr);
            });
        })
        .catch((err) => {
            displayError('You need to be in a voice channel to see the log.');
        });
});

$('#btnStats').on('click', (e) => {
    getGuildStats(guildID)
        .then((res) => {
            $('#modalStats').modal('show');
            var tab = $('#modalStats div.modal-body > table');

            tab.children('tr').each((_, tr) => tr.remove());
            res.forEach((r, i) => {
                var tr = document.createElement('tr');

                var tdNumber = document.createElement('td');
                tdNumber.innerHTML = `<span class="badge badge-primary">${i +
                    1}</span>`;
                tr.appendChild(tdNumber);

                var tdSound = document.createElement('td');
                tdSound.innerText = r.sound;
                tr.appendChild(tdSound);

                var tdCount = document.createElement('td');
                tdCount.innerText = r.count;
                tr.appendChild(tdCount);

                tab.append(tr);
            });
        })
        .catch((err) => {
            displayError('You need to be in a voice channel to see stats.');
        });
});

$('#ddFastTrigger').on('change', (e) => {
    var ident = $('#ddFastTrigger').val();
    postFastTrigger(ident === 'RANDOM', ident).catch((err) => {
        console.log('REST :: ERROR :: ', err);
        displayError(
            `<code>REST API ERROR</code> setting fast trigger:<br/><code>${r}</code>`
        );
    });
});

$('#sliderVol').on('input', (e) => {
    var val = $('#sliderVol').val();
    $('#labelVol')[0].innerText = val + '%';
});

$('#sliderVol').on('change', (e) => {
    var val = $('#sliderVol').val();
    ws.emit('VOLUME', parseInt(val));
});

$('#searchBox').on('input', (e) => {
    var tb = e.currentTarget;
    var val = tb.value;
    setTimeout(() => {
        if (val == tb.value) filterSoundsList(val, sounds);
    }, 250);
});

if (sortBy)
    $('#btnSortBy')[0].innerText =
        'SORT BY ' + (sortBy == 'DATE' ? 'NAME' : 'DATE');

// --------------------------
// --- WS EVENT HANDLERS

ws.on('ERROR', (data) => {
    eventDebug('ERROR', data);
    displayError(
        `<code>${data.data.code} - ${data.data.type}</code>&nbsp; ${
            data.data.message
        }`
    );
});

ws.on('HELLO', (data) => {
    eventDebug('HELLO', data);
    if (!data.data) return;

    if (data.data.admin) {
        $('#btnAdmin').removeClass('d-none');
    }

    if (data.data.connected) {
        $('#btnJoinLeave')[0].innerText = 'LEAVE';
        inChannel = true;
        guildID = data.data.voice_state.guild_id;
        setVolume(data.data.vol);
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
    setVolume(data.data.vol);
});

ws.on('END', (data) => {
    eventDebug('END', data);
    if (data.data.ident) {
        $(`#soundBtn-${data.data.ident}`).removeClass('playing');
    }

    if (chain !== null) {
        playFromQueue();
    }
});

ws.on('PLAY_ERROR', (data) => {
    eventDebug('PLAY_ERROR', data);
    if (chain !== null) {
        resetChaining();
    }
});

ws.on('STUCK', (data) => {
    eventDebug('STUCK', data);
});

ws.on('VOLUME_CHANGED', (data) => {
    eventDebug('VOLUME_CHANGED', data);
    setVolume(data.data.vol);
});

ws.on('JOINED', (data) => {
    eventDebug('JOINED', data);
    inChannel = true;
    $('#btnJoinLeave')[0].innerText = 'LEAVE';
    guildID = data.data.guild_id;
    setVolume(data.data.vol);
});

ws.on('LEFT', (data) => {
    eventDebug('LEFT', data);
    inChannel = false;
    $('#btnJoinLeave')[0].innerText = 'JOIN';
});
