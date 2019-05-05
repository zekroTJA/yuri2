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

    getGuildStats()
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
            spinner.removeClass('d-flex');
            spinner.addClass('d-none');
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
});

// ------------------------------------------------------------

updateStats();
