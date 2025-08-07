<script lang="ts">
	import Globe from 'ui/icons/Globe.svelte';
	import type { User } from '$lib/types/types';
	import HeartIcon from 'ui/icons/HeartIcon.svelte';
	import type { Component } from 'svelte';
	import StarsIcon from 'ui/icons/StarsIcon.svelte';
	import SmileIcon from 'ui/icons/SmileIcon.svelte';
	import ArchiveIcon from 'ui/icons/ArchiveIcon.svelte';
	import GraduationCapIcon from 'ui/icons/GraduationCapIcon.svelte';
	import BriefCaseIcon from 'ui/icons/BriefCaseIcon.svelte';
	import ConfettiIcon from 'ui/icons/ConfettiIcon.svelte';
	import SpeakIcon from 'ui/icons/SpeakIcon.svelte';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();

	const factIconMap: Record<string, Component> = {
		'I like': HeartIcon,
		'I love': HeartIcon,
		'My birthday is on the': ConfettiIcon,
		'I am from': Globe,
		'My favorite': StarsIcon,
		'I enjoy': SmileIcon,
		'I collect': ArchiveIcon,
		'I speak': SpeakIcon,
		'I studied': GraduationCapIcon,
		'I work as': BriefCaseIcon
	};

	const getIconForFact = (label: string): Component => {
		return factIconMap[label] || Globe;
	};
</script>

{#if user.facts.length! > 0}
	<p class="text-white/45 mix-blend-plus-lighter text-sm mt-4 font-medium relative z-[4]">Facts</p>
	{#each user.facts as fact, idx}
		{@const Icon = getIconForFact(fact.label)}

		<div class={['flex items-center gap-x-1.5', idx === 0 ? 'mt-1.5' : 'mt-2']}>
			<Icon height={20} width={20} class="text-white/45 relative z-[4]" />
			<p class="text-white/45 mix-blend-plus-lighter">
				{fact.label}
			</p>
			<p class="text-white/85 mix-blend-plus-lighter font-medium relative z-[4]">{fact.value}</p>
		</div>
	{/each}
{/if}
