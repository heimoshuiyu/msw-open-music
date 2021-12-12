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

// convert unix timestamp to %Y-%m-%d %H:%M:%S
export function convertIntToDateTime(timestamp) {
  var date = new Date(timestamp * 1000);
  var year = date.getFullYear();
  var month = date.getMonth() + 1;
  var day = date.getDate();
  var hour = date.getHours();
  var minute = date.getMinutes();
  var second = date.getSeconds();
  var time =
    year +
    "-" +
    (month < 10 ? "0" + month : month) +
    "-" +
    (day < 10 ? "0" + day : day) +
    " " +
    (hour < 10 ? "0" + hour : hour) +
    ":" +
    (minute < 10 ? "0" + minute : minute) +
    ":" +
    (second < 10 ? "0" + second : second);
  return time;
}

export function SayHello() {
  return "Hello";
}
