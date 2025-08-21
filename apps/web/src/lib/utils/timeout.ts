import { coreStore } from 'stores/coreStore.svelte';

export function setControllableTimeout(callback: () => void, delay: number) {
	const timeoutId = setTimeout(callback, delay);

	return {
		executeNow() {
			clearTimeout(timeoutId);
			callback();
		},
		clear() {
			clearTimeout(timeoutId);
			coreStore.callTimeout = undefined;
		}
	};
}
