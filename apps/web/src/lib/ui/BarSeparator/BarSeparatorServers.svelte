<script lang="ts">
	import { onMount } from 'svelte';
	import { gsap } from 'gsap';

	interface Props {
		title: string;
		tab: 'channels' | 'members';
	}

	let { title, tab = $bindable() }: Props = $props();

	let btnEl: HTMLButtonElement;
	let titleEl: HTMLSpanElement;
	let membersEl: HTMLSpanElement;
	let tl: gsap.core.Timeline;

	const EASE = 'cubic-bezier(0.625, 0.05, 0, 1)';

	// Function to update animation based on tab state
	const updateAnimationState = () => {
		if (!tl) return;

		if (tab === 'members') {
			tl.progress(1); // Show "Members list"
		} else {
			tl.progress(0); // Show title
		}
	};

	onMount(() => {
		gsap.set(titleEl, { opacity: 1, y: 0, filter: 'blur(0px)' });
		gsap.set(membersEl, { opacity: 0, y: 4, filter: 'blur(4px)' });

		tl = gsap
			.timeline({ paused: true, defaults: { duration: 0.1, ease: EASE } })
			.to(titleEl, { opacity: 0, y: -4, filter: 'blur(4px)' }, 0)
			.to(membersEl, { opacity: 1, y: 0, filter: 'blur(0px)' }, '>');

		updateAnimationState();

		const onEnter = () => {
			if (tab === 'channels') {
				tl.play();
			} else {
				tl.reverse();
			}
		};
		const onLeave = () => {
			if (tab === 'channels') {
				tl.reverse();
			} else {
				tl.play();
			}
		};

		btnEl.addEventListener('mouseenter', onEnter);
		btnEl.addEventListener('mouseleave', onLeave);

		return () => {
			btnEl.removeEventListener('mouseenter', onEnter);
			btnEl.removeEventListener('mouseleave', onLeave);
			tl.kill();
		};
	});

	$effect(() => {
		updateAnimationState();
	});
</script>

<button
	bind:this={btnEl}
	class="group relative overflow-hidden text-left bg-main-975 border-t-[0.5px] border-b-[0.5px] border-main-800 h-14 select-none hover:bg-main-950 transition-colors shrink-0"
	onclick={() => (tab = tab === 'channels' ? 'members' : 'channels')}
>
	<span bind:this={titleEl} class="absolute inset-0 flex items-center pointer-events-none px-4">
		{title}
	</span>
	<span bind:this={membersEl} class="absolute inset-0 flex items-center pointer-events-none px-4">
		Members list
	</span>
</button>
