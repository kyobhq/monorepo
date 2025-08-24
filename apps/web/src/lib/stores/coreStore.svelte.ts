import type { User } from '$lib/types/types';
import type { setControllableTimeout } from 'utils/timeout';
import { backend } from './backendStore.svelte';
import { userStore } from './userStore.svelte';

const MAX_SEEN_PROFILES = 5;

export class Core {
	firstLoad = $state({ sidebar: false });
	serversLoaded = $state(false);
	richInputLength = $state(0);
	serverDialog = $state(false);
	categoryDialog = $state(false);
	friendsDialog = $state(false);
	channelDialog = $state({ open: false, category_id: '' });
	modDialog = $state({ open: false, action: 'kick', server_id: '', user_id: '' });
	restrictionDialog = $state({ open: false, restriction: '', title: '', reason: '' });
	settingsDialog = $state({ open: false, section: '' });
	serverSettingsDialog = $state({ open: false, server_id: '', section: '' });
	categorySettingsDialog = $state({ open: false, category_id: '', section: '' });
	channelSettingsDialog = $state({ open: false, channel_id: '', section: '' });
	userSettingsDialog = $state({ open: false, section: '' });
	destructiveDialog = $state({
		open: false,
		title: '',
		subtitle: '',
		buttonText: '',
		onclick: () => {}
	});
	pressingShift = $state(false);
	isUsingKeyboard = $state(false);
	profile = $state<{ open: boolean; user: User | null; el: HTMLElement | null; side: string }>({
		open: false,
		user: null,
		el: null,
		side: 'top'
	});
	seenProfiles = $state<User[]>([]);
	callTimeout = $state<ReturnType<typeof setControllableTimeout>>();

	private handleShiftUp = (e: Event) => {
		const keyboardEvent = e as KeyboardEvent;
		if (keyboardEvent.key === 'Shift') {
			this.pressingShift = false;
		}
	};

	private handleKeyboardInput = (e: Event) => {
		const keyboardEvent = e as KeyboardEvent;
		if (keyboardEvent.key === 'Tab') {
			this.setInputMethod('keyboard');
		}
		if (keyboardEvent.key === 'Shift') {
			this.pressingShift = true;
		}
	};

	private handleNonKeyboardInput = () => {
		this.setInputMethod('mouse');
	};

	private setInputMethod(method: 'keyboard' | 'mouse') {
		this.isUsingKeyboard = method === 'keyboard';
		this.updateDataAttributes();
	}

	private updateDataAttributes() {
		if (typeof document !== 'undefined') {
			document.documentElement.setAttribute('data-keyboard-user', String(this.isUsingKeyboard));
		}
	}

	initializeKeyboardDetection(): void {
		if (typeof document === 'undefined') return;

		this.updateDataAttributes();

		const keyboardEvents = ['keydown', 'keyup'];
		const nonKeyboardEvents = [
			'mousedown',
			'mouseup',
			'pointerdown',
			'pointerup',
			'touchstart',
			'touchend'
		];

		keyboardEvents.forEach((event) => {
			const handler = event === 'keyup' ? this.handleShiftUp : this.handleKeyboardInput;
			document.addEventListener(event, handler, true);
		});

		nonKeyboardEvents.forEach((event) => {
			document.addEventListener(event, this.handleNonKeyboardInput, true);
		});
	}

	cleanupKeyboardDetection(): void {
		if (typeof document === 'undefined') return;

		const keyboardEvents = ['keydown', 'keyup'];
		const nonKeyboardEvents = [
			'mousedown',
			'mouseup',
			'pointerdown',
			'pointerup',
			'touchstart',
			'touchend'
		];

		keyboardEvents.forEach((event) => {
			const handler = event === 'keyup' ? this.handleShiftUp : this.handleKeyboardInput;
			document.removeEventListener(event, handler, true);
		});

		nonKeyboardEvents.forEach((event) => {
			document.removeEventListener(event, this.handleNonKeyboardInput, true);
		});
	}

	openMyProfile(el: HTMLElement, side: 'top' | 'bottom' | 'left' | 'right' = 'top'): void {
		if (this.profile.user) {
			this.closeProfile();
		} else {
			this.setProfileAndOpen(userStore.user!, el, side);
		}
	}

	async openProfile(
		userID: string,
		el: HTMLElement,
		side: 'top' | 'bottom' | 'left' | 'right' = 'top'
	): Promise<void> {
		if (this.profile.user) {
			this.closeProfile();
			return;
		}

		const cachedUser = this.seenProfiles.find((profile) => profile.id === userID);
		if (cachedUser) {
			this.setProfileAndOpen(cachedUser, el, side);
			return;
		}

		const res = await backend.getUserProfile(userID);
		res.match(
			(user) => {
				this.setProfileAndOpen(user, el, side);
				this.cacheUserProfile(user);
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	closeProfile(): void {
		this.profile = {
			open: false,
			user: null,
			el: null,
			side: 'top'
		};
	}

	private setProfileAndOpen(user: User, el: HTMLElement, side: string): void {
		this.profile = {
			open: false,
			user,
			el,
			side
		};

		setTimeout(() => {
			if (this.profile.user) {
				this.profile.open = true;
			}
		}, 10);
	}

	private cacheUserProfile(user: User): void {
		if (this.seenProfiles.length === MAX_SEEN_PROFILES) {
			this.seenProfiles.shift();
		}
		this.seenProfiles.push(user);
	}
}

export const coreStore = new Core();
