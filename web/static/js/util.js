/**
 * Returns an object of cookies with
 * the name of the cookie as key and
 * the value of the cookies as value.
 * @returns cookies
 */
function getCookies() {
    var c = {};
    document.cookie
        .split(";")
        .map((v) => v.trim().split('='))
        .forEach((v) => c[v[0]] = v[1]);
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
    Object.keys(getCookies()).forEach((k) => 
        document.cookie = `${k}=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/`
    );
}