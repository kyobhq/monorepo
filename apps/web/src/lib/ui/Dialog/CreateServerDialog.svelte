<script lang="ts">
	import { CreateServerSchema } from '$lib/types/schemas';
	import type { Server } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import Cropper from 'svelte-easy-crop';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import DefaultDialog from './DefaultDialog.svelte';
	import Switch from 'ui/Switch/Switch.svelte';
	import DialogFooter from './DialogFooter.svelte';
	import { goto } from '$app/navigation';

	let avatar = $state<string | undefined>();
	let crop = $state({ x: 0, y: 0 });
	let zoom = $state(1);
	let minZoom = $state(3);
	let maxZoom = $state(5);

	const { form, errors, enhance } = superForm(defaults(valibot(CreateServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateServerSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.position = serverStore.getLastPosition();

				const res = await backend.createServer(form.data);
				res.match(
					(s) => {
						const server: Server = {
							...s,
							categories: {},
							main_color: '12,14,14'
						};

						serverStore.addServer(server);
						coreStore.serverDialog = false;
						goto(`/servers/${server.id}`);
					},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		}
	});

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		const image = target.files?.[0];

		if (image) {
			const dataUrl = URL.createObjectURL(image);
			const img = new Image();

			img.onload = () => {
				const aspectAvatar = 1;
				const aspectImage = img.naturalWidth / img.naturalHeight;

				minZoom = aspectImage > 1 ? aspectImage / aspectAvatar : aspectAvatar / aspectImage;
				zoom = minZoom;

				URL.revokeObjectURL(dataUrl);
			};
			img.src = dataUrl;
			avatar = dataUrl;
			$form.avatar = image;
		}
	}
</script>

<DefaultDialog
	bind:state={coreStore.serverDialog}
	title="Create a server"
	subtitle="It can be either public or private."
>
	<form method="post" use:enhance enctype="multipart/form-data">
		<div class="flex items-center justify-between px-8">
			<div
				class={[
					'group relative h-[85px] w-[85px] overflow-hidden text-transparent transition-colors hover:cursor-pointer border-[0.5px]',
					$errors.avatar
						? 'hocus:bg-red-400/25 border-red-400 hocus:inner-red-400/40 bg-red-400/15'
						: 'border-main-800 bg-main-950 hocus:bg-main-900'
				]}
			>
				<input
					type="file"
					id="avatar"
					name="avatar"
					onchange={onFile}
					aria-label="Server avatar"
					class="absolute h-full w-full text-transparent hover:cursor-pointer"
				/>
				{#if $form.avatar}
					<Cropper
						image={avatar}
						cropSize={{ height: 85, width: 85 }}
						cropShape="rect"
						showGrid={false}
						bind:crop
						bind:zoom
						{minZoom}
						{maxZoom}
						oncropcomplete={(e) => {
							$form.crop = e.pixels;
						}}
					/>
				{/if}
			</div>
		</div>

		<FormInput
			title="Server name"
			id="server-name"
			type="text"
			bind:error={$errors.name}
			bind:inputValue={$form.name}
			placeholder="My cool community"
			class="mt-4 px-8"
		/>

		<FormInput
			title="Server description"
			id="server-description"
			type="rich"
			bind:error={$errors.description}
			bind:inputValue={$form.description}
			placeholder="Here we do..."
			class="mt-4 px-8"
			inputClass="w-full"
		/>

		<div class="px-8 mt-4">
			<Switch
				active={$form.public}
				action={() => ($form.public = true)}
				reverse={() => ($form.public = false)}
				label="Make this server public (people will be able to see it)"
			/>
		</div>

		<DialogFooter buttonText="Create server" />
	</form>
</DefaultDialog>
