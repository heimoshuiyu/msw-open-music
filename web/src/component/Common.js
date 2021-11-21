export function CalcReadableFilesize(filesize) {
  if (filesize < 1024) {
    return filesize;
  }
  if (filesize < 1024 * 1024) {
    return Math.round(filesize / 1024) + "K";
  }
  if (filesize < 1024 * 1024 * 1024) {
    return Math.round(filesize / 1024 / 1024) + "M";
  }
  if (filesize < 1024 * 1024 * 1024 * 1024) {
    return Math.round(filesize / 1024 / 1024 / 1024) + "G";
  }
}

export function CalcReadableFilesizeDetail(filesize) {
  if (filesize < 1024 * 1024) {
    return filesize;
  }
  if (filesize < 1024 * 1024 * 1024) {
    return numberWithCommas(Math.round(filesize / 1024)) + "K";
  }
  if (filesize < 1024 * 1024 * 1024 * 1024) {
    return numberWithCommas(Math.round(filesize / 1024 / 1024)) + "M";
  }
  if (filesize < 1024 * 1024 * 1024 * 1024 * 1024) {
    return numberWithCommas(Math.round(filesize / 1024 / 1024 / 1024)) + "G";
  }
}

function numberWithCommas(x) {
  x = x.toString();
  var pattern = /(-?\d+)(\d{3})/;
  while (pattern.test(x)) x = x.replace(pattern, "$1,$2");
  return x;
}

export function SayHello() {
  return "Hello";
}
