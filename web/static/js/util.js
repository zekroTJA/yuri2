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
