export function formatMessageTime(time: string) {
  const now = new Date();
  const timestamp = new Date(time);

  let day = timestamp.getDate().toString();
  let month = timestamp.getMonth().toString();
  const year = timestamp.getFullYear().toString().slice(2).toString();
  let hour = timestamp.getHours().toString();
  let minutes = timestamp.getMinutes().toString();

  day = Number(day) < 10 ? '0' + day : day;
  month = Number(month) < 10 ? '0' + month : month;
  hour = Number(hour) < 10 ? '0' + hour : hour;
  minutes = Number(minutes) < 10 ? '0' + minutes : minutes;

  if (now.toDateString() === timestamp.toDateString()) {
    return `${hour}:${minutes}`;
  }

  return `${day}/${month}/${year}, ${hour}:${minutes}`;
}

export const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));
