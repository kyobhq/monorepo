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

export function joinedAt(time: string) {
	const now = new Date();
	const timestamp = new Date(time);

	const diff = now.getTime() - timestamp.getTime();

	const years = Math.floor(diff / (1000 * 60 * 60 * 24 * 365));
	const months = Math.floor(diff / (1000 * 60 * 60 * 24 * 30));
	const days = Math.floor(diff / (1000 * 60 * 60 * 24));
	const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
	const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

	if (years > 0) return `${years > 1 ? `${years} years` : 'a year'} ago`;
	if (months > 0) return `${months > 1 ? `${months} months` : 'a month'} ago`;
	if (days > 0) return `${days > 1 ? `${days} days` : 'a day'} ago`;
	if (hours > 0) return `${hours > 1 ? `${hours} hours` : 'an hour'} ago`;
	if (minutes > 0) return `${minutes > 1 ? `${minutes} minutes` : 'a minute'} ago`;
	return 'just now';
}

export function expiresAt(time: string) {
	const now = new Date();
	const timestamp = new Date(time);

	const diff = timestamp.getTime() - now.getTime();
	const days = Math.floor(diff / (1000 * 60 * 60 * 24));
	const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
	const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

	if (days > 0) return `${days} days`;
	if (hours > 0) return `${hours} hours`;
	if (minutes > 0) return `${minutes} minutes`;
	return 'expired';
}

export const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));
