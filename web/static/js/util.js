/** @format */

'use strict';

/**
 * Returns an object of cookies with
 * the name of the cookie as key and
 * the value of the cookies as value.
 * @returns cookies
 */
function getCookies() {
    var c = {};
    document.cookie
        .split(';')
        .map((v) => v.trim().split('='))
        .forEach((v) => (c[v[0]] = v[1]));
    return c;
}

/**
 * Gets the value of a cookie by its
 * name. if the cookie is not set,
 * 'undefined' will be returned.
 * @param {string} name name of the cookie
 * @returns cookie value
 */
function getCookieValue(name) {
    return getCookies()[name];
}

/**
 * Deletes all cookies from this page.
 */
function deleteAllCookies() {
    Object.keys(getCookies()).forEach(
        (k) =>
            (document.cookie = `${k}=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/`)
    );
}

/**
 * Retrurns the passed time parsed as
 * string. If passed date is null, current
 * time will be used.
 * @param {Date?} date date to parse
 * @returns parsed date string
 */
function getTime(date) {
    function btf(inp) {
        if (inp < 10) return '0' + inp;
        return inp;
    }
    if (!date) date = new Date();
    var y = date.getFullYear(),
        m = btf(date.getMonth() + 1),
        d = btf(date.getDate()),
        h = btf(date.getHours()),
        min = btf(date.getMinutes()),
        s = btf(date.getSeconds());
    return `${y}/${m}/${d} - ${h}:${min}:${s}`;
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

function addRowToTable(tb, ...entries) {
    var tr = document.createElement('tr');
    entries.forEach((e) => {
        if (typeof e === 'object') {
            tr.appendChild(e);
        } else {
            var td = document.createElement('td');
            td.innerText = e;
            tr.appendChild(td);
        }
    });
    tb.append(tr);
}

function elemInnerText(elem, innerText) {
    var e = document.createElement(elem);
    e.innerText = innerText;
    return e;
}

function byteCountFormatter(bc) {
    const k = 1024;
    const fix = 2;
    if (bc < k) return `${bc} B`;
    if (bc < k * k) return `${(bc / k).toFixed(fix)} kiB`;
    if (bc < k * k * k) return `${(bc / k / k).toFixed(fix)} MiB`;
    if (bc < k * k * k * k) return `${(bc / k / k / k).toFixed(fix)} GiB`;
    return `${(bc / k / k / k / k).toFixed(fix)} TiB`;
}

function padFront(num, len, char) {
    num = num.toString();
    while (num.length < len) {
        num = char + num;
    }
    return num;
}

function toDDHHMMSS(secs) {
    var dd = Math.floor(secs / 86400);
    var hh = Math.floor((secs % 86400) / 3600);
    var mm = Math.floor(((secs % 86400) % 3600) / 60);
    var ss = Math.floor(((secs % 86400) % 3600) % 60);
    return `${padFront(dd, 2, '0')}:${padFront(hh, 2, '0')}:${padFront(
        mm,
        2,
        '0'
    )}:${padFront(ss, 2, '0')}`;
}
