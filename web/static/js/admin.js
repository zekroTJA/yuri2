/** @format */

'use strict';

var guilds = [];
var VCs = [];

// SYSTEM STATS
function updateStats() {
    var spinnerInfo = $('#spinnerInfo');
    var spinnerGuilds = $('#spinnerGuilds');
    var spinnerVCs = $('#spinnerVCs');
    var tbStats = $('#tbStats > tbody');
    var tbGuilds = $('#tbGuilds > tbody');
    var tbVCs = $('#tbVCs > tbody');

    tbStats.empty();
    tbGuilds.empty();
    tbVCs.empty();

    spinnerInfo.removeClass('d-none');
    spinnerInfo.addClass('d-flex');
    spinnerGuilds.removeClass('d-none');
    spinnerGuilds.addClass('d-flex');
    spinnerVCs.removeClass('d-none');
    spinnerVCs.addClass('d-flex');

    getAdminStats()
        .then((res) => {
            guilds = res.guilds;
            VCs = res.voice_connections;

            addRowToTable(
                tbStats,
                'Guilds',
                elemInnerText('code', res.guilds.length)
            );
            addRowToTable(
                tbStats,
                'Active VCs',
                elemInnerText('code', res.voice_connections.length)
            );
            addRowToTable(
                tbStats,
                'Uptime',
                elemInnerText('code', toDDHHMMSS(res.system.uptime_seconds))
            );
            addRowToTable(tbStats, 'OS', elemInnerText('code', res.system.os));
            addRowToTable(
                tbStats,
                'Arch',
                elemInnerText('code', res.system.arch)
            );
            addRowToTable(
                tbStats,
                'Go Version',
                elemInnerText('code', res.system.go_version)
            );
            addRowToTable(
                tbStats,
                'Used CPU Threads',
                elemInnerText('code', res.system.cpu_used_cores)
            );
            addRowToTable(
                tbStats,
                'Running Go Routines',
                elemInnerText('code', res.system.go_routines)
            );
            addRowToTable(
                tbStats,
                'Used Heap',
                elemInnerText('code', byteCountFormatter(res.system.heap_use_b))
            );
            addRowToTable(
                tbStats,
                'Used Stack',
                elemInnerText(
                    'code',
                    byteCountFormatter(res.system.stack_use_b)
                )
            );

            spinnerInfo.removeClass('d-flex');
            spinnerInfo.addClass('d-none');

            res.guilds.slice(0, 10).forEach((g) => {
                addRowToTable(tbGuilds, g.name, elemInnerText('code', g.id));
            });

            spinnerGuilds.removeClass('d-flex');
            spinnerGuilds.addClass('d-none');

            res.voice_connections.slice(0, 10).forEach((v) => {
                addRowToTable(
                    tbVCs,
                    v.guild.name,
                    elemInnerText('code', v.vc_id)
                );
            });

            spinnerVCs.removeClass('d-flex');
            spinnerVCs.addClass('d-none');
        })
        .catch((err) => {
            displayError(
                `<code>REST API ERROR</code> getting system stats failed<br/><code>${err}</code>`
            );

            spinnerGuilds.removeClass('d-flex');
            spinnerGuilds.addClass('d-none');
            spinnerInfo.removeClass('d-flex');
            spinnerInfo.addClass('d-none');
            spinnerVCs.removeClass('d-flex');
            spinnerVCs.addClass('d-none');
        });
}

// SOUND STATS
function updateSoundStats() {
    var spinnerSoundInfo = $('#spinnerSoundInfo');
    var tbSoundStats = $('#tbSoundStats > tbody');

    tbSoundStats.empty();

    spinnerSoundInfo.removeClass('d-none');
    spinnerSoundInfo.addClass('d-flex');
    getAdminSoundStats()
        .then((res) => {
            addRowToTable(
                tbSoundStats,
                'Number of Sound Files',
                elemInnerText('code', res.sounds_len)
            );
            addRowToTable(
                tbSoundStats,
                'Size of Sound Files',
                elemInnerText('code', byteCountFormatter(res.size_b))
            );
            addRowToTable(
                tbSoundStats,
                'Log Records',
                elemInnerText('code', res.log_len)
            );

            spinnerSoundInfo.removeClass('d-flex');
            spinnerSoundInfo.addClass('d-none');
        })
        .catch((err) => {
            displayError(
                `<code>REST API ERROR</code> getting sound stats failed<br/><code>${err}</code>`
            );

            spinnerSoundInfo.removeClass('d-flex');
            spinnerSoundInfo.addClass('d-none');
        });
}

// ------------------------------------------------------------

$('#btnMoreGuilds').on('click', () => {
    var tab = $('#modalGuilds div.modal-body > table');
    tab.empty();

    guilds.forEach((g) => {
        addRowToTable(tab, g.name, g.id);
    });

    $('#modalGuilds').modal('show');
});

$('#btnMoreVCs').on('click', () => {
    var tab = $('#modalVCs div.modal-body > table');
    tab.empty();

    VCs.forEach((v) => {
        addRowToTable(tab, v.guild.name, v.guild.id, v.vc_id);
    });

    $('#modalVCs').modal('show');
});

$('#btRefresh').on('click', () => {
    updateStats();
    updateSoundStats();
});

$('#btRestart').on('click', (e) => {
    btnAccept(e, (t) => {
        postRestart()
            .then(() => {
                var alert = $('#warnAlert');
                alert.removeClass('d-none');
                setTimeout(() => {
                    alert[0].style.opacity = '1';
                    alert[0].style.transform = 'translateY(0px)';
                }, 10);
                var secs = 10;
                alert.text(
                    `Restarting... This page atomatically reloads after ${secs} seconds...`
                );
                setInterval(() => {
                    alert.text(
                        `Restarting... This page atomatically reloads after ${--secs} seconds...`
                    );
                    if (secs === 0) window.location = '/admin';
                }, 1000);
            })
            .catch((err) => {
                displayError(
                    `<code>REST API ERROR</code> restart failed<br/><code>${err}</code>`
                );
            });
    });
});

$('#btRefetch').on('click', (e) => {
    btnAccept(e, (t) => {
        postRefetch()
            .then(() => {
                displayInfo('Refetched sounds.', 5000);
            })
            .catch((err) => {
                displayError(
                    `<code>REST API ERROR</code> refetch failed<br/><code>${err}</code>`
                );
            });
    });
});

// ------------------------------------------------------------

updateStats();
updateSoundStats();
