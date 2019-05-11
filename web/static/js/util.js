/** @format */

'use strict';

var btnStats = {};

// ------------------------------------------------------------

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

/**
 * Fades in an error box and sets its
 * text to the passed error description.
 * After the specified time in milliseconds,
 * the alert box will be faded out.
 * @param {string} desc error description
 * @param {number} time time until fade out (in ms)
 */
function displayError(desc, time) {
    if (!time) time = 8000;

    var alertBox = $('#errorAlert')[0];
    $('#errorAlertText')[0].innerHTML = desc;
    window.scrollTo(0, 0);

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

/**
 * Fades in an info box and sets its
 * text to the passed info description.
 * After the specified time in milliseconds,
 * the alert box will be faded out.
 * @param {string} desc info description
 * @param {number} time time until fade out (in ms)
 */
function displayInfo(desc, time) {
    if (!time) time = 8000;

    var alertBox = $('#infoAlert');
    $('#infoAlert')[0].innerHTML = desc;

    // fade in
    alertBox.removeClass('d-none');
    setTimeout(() => {
        alertBox[0].style.opacity = '1';
        alertBox[0].style.transform = 'translateY(0px)';
    }, 10);
    // fade out
    setTimeout(() => {
        alertBox[0].style.opacity = '0';
        alertBox[0].style.transform = 'translateY(-20px)';
    }, time);
    setTimeout(() => {
        alertBox.addClass('d-none');
    }, time + 250);
}

/**
 * Adds a row element to the passed table
 * containing the
 * @param {Object} tb
 * @param  {...any} entries
 */
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

/**
 * Creates a new element of the passed type
 * and sets its inner text to the specified
 * text.
 * @param {string} elem type of object to be created
 * @param {string} innerText inner text
 * @returns created DOM element
 */
function elemInnerText(elem, innerText) {
    var e = document.createElement(elem);
    e.innerText = innerText;
    return e;
}

/**
 * Formats a passed number of bytes to
 * a useful number with the adequate
 * unit prefix.
 * @param {number} bc byte count
 */
function byteCountFormatter(bc) {
    const k = 1024;
    const fix = 2;
    if (bc < k) return `${bc} B`;
    if (bc < k * k) return `${(bc / k).toFixed(fix)} kiB`;
    if (bc < k * k * k) return `${(bc / k / k).toFixed(fix)} MiB`;
    if (bc < k * k * k * k) return `${(bc / k / k / k).toFixed(fix)} GiB`;
    return `${(bc / k / k / k / k).toFixed(fix)} TiB`;
}

/**
 * Pads char to the passed num until
 * nums length is equal to len.
 * num and char will we converted to
 * string using.
 * @param {*} num object to be padded
 * @param {number} len minimum result string length
 * @param {*} char the string which will be padded
 * @returns result string
 */
function padFront(num, len, char) {
    num = num.toString();
    char = char.toString();
    while (num.length < len) {
        num = char + num;
    }
    return num;
}

/**
 * Parses a duration in seconds to
 * DD:HH:MM:SS format.
 * @param {number} secs seconds
 * @returns formatted string
 */
function toDDHHMMSS(secs) {
    var dd = padFront(Math.floor(secs / 86400), 2, 0);
    var hh = padFront(Math.floor((secs % 86400) / 3600), 2, 0);
    var mm = padFront(Math.floor(((secs % 86400) % 3600) / 60), 2, 0);
    var ss = padFront(Math.floor(((secs % 86400) % 3600) % 60), 2, 0);
    return `${dd}:${hh}:${mm}:${ss}`;
}

/**
 * Taskes a button onclick event. The fisr time, the
 * button was pressed, the color of the background
 * of the button will change to red and the text will
 * be changed to 'Sure?'. If the button was pressed
 * again in this state, the passed callback winn be
 * executed. The callback is getting passed the button
 * DOM element object. After that, the button will
 * reset to its initial state.
 * If a second click does not occur after 3 seconds,
 * the button will return to its initial state.
 * @param {Object} e button onclick event
 * @param {function} cb callback on acception
 */
function btnAccept(e, cb) {
    var t = $(e.target);
    var reset = () => {
        t.removeClass('need-accept');
        t.addClass('bg-clr-orange');
        t[0].innerText = btnStats[t[0].id].innerText;
        delete btnStats[t[0].id];
    };

    if (t.hasClass('need-accept')) {
        cb.bind(null, t).call();
        reset();
    } else {
        t.removeClass('bg-clr-orange');
        t.addClass('need-accept');
        t.oldText = t[0].innerText;
        btnStats[t[0].id] = {
            innerText: t[0].innerText,
        };
        var width = t.width();
        t[0].innerText = 'Sure?';
        t.width(width);
        setTimeout(() => {
            if (t.hasClass('need-accept')) {
                reset();
            }
        }, 3000);
    }
}
