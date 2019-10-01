/** @format */

export function byteCountFormatter(bc: number) {
  const k = 1024;
  const fix = 2;
  if (bc < k) {
    return `${bc} B`;
  }
  if (bc < k * k) {
    return `${(bc / k).toFixed(fix)} kiB`;
  }
  if (bc < k * k * k) {
    return `${(bc / k / k).toFixed(fix)} MiB`;
  }
  if (bc < k * k * k * k) {
    return `${(bc / k / k / k).toFixed(fix)} GiB`;
  }
  return `${(bc / k / k / k / k).toFixed(fix)} TiB`;
}

export function padFront(num: any, len: number, char: any) {
  num = num.toString();
  char = char.toString();
  while (num.length < len) {
    num = char + num;
  }
  return num;
}
