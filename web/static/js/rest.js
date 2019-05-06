/** @format */

'use strict';

function restDebugRequest(method, endpoint) {
    console.log(`REST API :: REQUEST :: ${method} ${endpoint}`);
}

function restDebugRespone(data, status) {
    console.log(`REST API :: RESPONSE :: [${status}]`, data);
}

// ------------------------------
// --- REQUESTS

// GET /api/localsounds
function getLocalSounds(sortBy) {
    var url = '/api/localsounds';
    if (sortBy) url += '?sort=' + sortBy.toUpperCase();

    restDebugRequest('GET', url);

    return new Promise((resolve, rejects) => {
        $.getJSON(url, (res, s) => {
            restDebugRespone(res, s);
            if (s === 'success') {
                resolve(res.results);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });
    });
}

// GET /api/logs/:GUILDID
function getGuildLog(guildID) {
    var url = `/api/logs/${guildID}?limit=50`;
    restDebugRequest('GET', url);

    return new Promise((resolve, rejects) => {
        $.getJSON(url, (res, s) => {
            restDebugRespone(res, s);
            if (s === 'success') {
                resolve(res.results);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });
    });
}

// GET /api/stats/:GUILDID
function getGuildStats(guildID) {
    var url = `/api/stats/${guildID}?limit=50`;
    restDebugRequest('GET', url);

    return new Promise((resolve, rejects) => {
        $.getJSON(url, (res, s) => {
            restDebugRespone(res, s);
            if (s === 'success') {
                resolve(res.results);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });
    });
}

// GET /api/admin/restart
function postRestart() {
    var url = `/api/admin/restart`;
    restDebugRequest('POST', url);

    return new Promise((resolve, rejects) => {
        $.post(url, (res, s) => {
            restDebugRespone(res, s);
            if (s === 'success') {
                resolve(res);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });
    });
}

// GET /api/admin/refetch
function postRefetch() {
    var url = `/api/admin/refetch`;
    restDebugRequest('POST', url);

    return new Promise((resolve, rejects) => {
        $.post(url, (res, s) => {
            restDebugRespone(res, s);
            if (s === 'success') {
                resolve(res);
            } else {
                rejects(res, s);
            }
        }).fail((e) => {
            rejects(e);
        });
    });
}
