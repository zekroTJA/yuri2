/** @format */

'use strict';

document.addEventListener('paste', (e) => {
    e.stopPropagation();
    e.preventDefault();
    var pastedData = (e.clipboardData || window.clipboardData).getData('Text');

    var i = pastedData.indexOf('&');
    if (i > -1) pastedData = pastedData.substring(0, i);

    var i = pastedData.indexOf('?');
    if (i > -1 && pastedData.substr(i - 5, 5) !== 'watch')
        pastedData = pastedData.substring(0, i);

    var cut = 0;

    i = pastedData.indexOf('youtube.com/watch?v=');
    if (i > -1) cut = i + 'youtube.com/watch?v='.length;

    i = pastedData.indexOf('youtu.be/');
    if (i > -1) cut = i + 'youtu.be/'.length;

    pastedData = pastedData.substring(cut);

    $.getJSON(
        `https://noembed.com/embed?url=https://www.youtube.com/watch?v=${pastedData}`
    )
        .done((res) => {
            if (res.error) {
                displayError(
                    `This YouTube video can not be accessed. Maybe wrong video ID?`
                );
            } else {
                var ytUrl = `https://www.youtube-nocookie.com/embed/${pastedData}?controls=0`;

                var mod = $('#modalYouTube');
                mod.find('iframe')[0].src = ytUrl;
                mod.find('button[name=accept]').one('click', () => {
                    ws.emit('PLAY', {
                        ident: pastedData,
                        source: 1,
                    });
                });
                mod.modal('show');
            }
        })
        .fail((err) => {
            displayError(
                `<code>REST API ERROR</code> failed getting video inforamtion:<br/><code>${err}</code>`
            );
        });
});
