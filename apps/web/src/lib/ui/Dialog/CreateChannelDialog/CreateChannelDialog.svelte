<script lang="ts">
	import { page } from '$app/state';
	import { CreateChannelSchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import DefaultDialog from '../DefaultDialog/DefaultDialog.svelte';
	import DialogFooter from '../DialogFooter/DialogFooter.svelte';
	import { ChannelTypes } from '$lib/types/types';
	import HashChat from 'ui/icons/HashChat.svelte';
	import VolumeIcon from 'ui/icons/VolumeIcon.svelte';
	import GalleryIcon from 'ui/icons/GalleryIcon.svelte';
	import CarouselIcon from 'ui/icons/CarouselIcon.svelte';

	const CHANNEL_ICONS = {
		textual: HashChat,
		voice: VolumeIcon,
		gallery: GalleryIcon,
		kanban: CarouselIcon
	};

	const CHANNEL_DESCRIPTIONS = {
		textual: 'Textual channels are used for text-based conversations.',
		voice: 'Voice channels are used for voice-based conversations.',
		gallery: 'Gallery channels are used for image and video sharing.',
		kanban: 'Kanban channels are used for Kanban boards.'
	};

	const { form, errors, enhance } = superForm(defaults(valibot(CreateChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateChannelSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.server_id = page.params.server_id || '';
				form.data.category_id = coreStore.channelDialog.category_id;
				form.data.position = channelStore.getChannelsLastPositionInCategory(
					page.params.server_id || '',
					coreStore.channelDialog.category_id
				);
				coreStore.channelDialog.open = false;

				const res = await backend.createChannel(form.data);
				res.match(
					() => {},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		}
	});

	$effect(() => {
		if (!$form.type || $form.type === '') {
			$form.type = 'textual';
		}
	});
</script>

<DefaultDialog
	bind:state={coreStore.channelDialog.open}
	title="Create a channel"
	subtitle="Plenty of cool ways to interact with people!"
>
	<form method="post" use:enhance>
		<FormInput
			title="Name"
			id="channel-name"
			type="text"
			bind:error={$errors.name}
			bind:inputValue={$form.name}
			placeholder="General"
			class="mt-4 px-8"
		/>

		<FormInput
			title="Description"
			id="description-name"
			type="text"
			bind:error={$errors.description}
			bind:inputValue={$form.description}
			placeholder="Here you can talk about everything!"
			class="mt-4 px-8"
		/>

		<ul
			class="flex flex-col gap-y-1 mx-8 border-[0.5px] p-1 border-main-800 bg-main-900 rounded-md mt-4 overflow-auto h-[10rem]"
		>
			{#each Object.values(ChannelTypes) as type}
				{@const Icon = CHANNEL_ICONS[type as keyof typeof CHANNEL_ICONS]}
				<li>
					<label
						for={`type-${type}`}
						class={[
							'flex items-center gap-x-4 border-[0.5px] rounded-sm px-4 py-2 cursor-pointer hover:bg-main-800 transition-colors',
							$form.type === type
								? 'bg-main-800 border-main-700'
								: 'opacity-70 bg-main-850 border-main-800'
						]}
					>
						<input
							id={`type-${type}`}
							type="radio"
							name="type"
							value={type}
							class="opacity-0 absolute h-0 w-0"
							bind:group={$form.type}
						/>
						<Icon class={$form.type === type ? 'text-main-200' : 'text-main-300'} />
						<div class="flex flex-col">
							<p>{type.charAt(0).toUpperCase() + type.slice(1)}</p>
							<p
								class={[
									'text-sm leading-tight',
									$form.type === type ? 'text-main-300' : 'text-main-500'
								]}
							>
								{CHANNEL_DESCRIPTIONS[type]}
							</p>
						</div>
					</label>
				</li>
			{/each}
		</ul>

		<!-- <div class="px-8 mt-10"> -->
		<!-- 	<Switch -->
		<!-- 		active={$form.e2ee} -->
		<!-- 		action={() => ($form.e2ee = true)} -->
		<!-- 		reverse={() => ($form.e2ee = false)} -->
		<!-- 		label="End to end encryption*" -->
		<!-- 	/> -->
		<!-- </div> -->
		<!---->
		<!-- <p class="text-main-600 text-xs px-8 mt-4"> -->
		<!-- 	*End to end encryption make all your messages in this channel unreadable for everyone outside -->
		<!-- 	this channel, even us. However be careful, since encryption happens on your side it makes some -->
		<!-- 	features unusable like: Search, Automod, Usage of bots, etc. We advise to use those channels -->
		<!-- 	for conversations that should stay fully private even if it means less features. -->
		<!-- </p> -->

		<DialogFooter buttonText="Create channel" />
	</form>
</DefaultDialog>
