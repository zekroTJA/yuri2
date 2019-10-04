/** @format */

import { padFront } from './util.format';

export function toDDHHMMSS(secs: number) {
  const dd = padFront(Math.floor(secs / 86400), 2, 0);
  const hh = padFront(Math.floor((secs % 86400) / 3600), 2, 0);
  const mm = padFront(Math.floor(((secs % 86400) % 3600) / 60), 2, 0);
  const ss = padFront(Math.floor(((secs % 86400) % 3600) % 60), 2, 0);
  return `${dd}:${hh}:${mm}:${ss}`;
}
