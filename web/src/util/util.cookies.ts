/** @format */

export function getCookies() {
  var c = {};
  document.cookie
    .split(';')
    .map((v) => v.trim().split('='))
    .forEach((v) => (c[v[0]] = v[1]));
  return c;
}

export function getCookieValue(name: string) {
  return getCookies()[name];
}
