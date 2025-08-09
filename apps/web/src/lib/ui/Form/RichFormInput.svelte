<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { CharacterCount, Placeholder } from '@tiptap/extensions';
	import { Plugin } from '@tiptap/pm/state';

	let element: Element;
	let editor: Editor;
	let contentSet = $state(false);

	let { placeholder, content = $bindable(), class: classes } = $props();

	const limitConsecutiveBreaksPlugin = new Plugin({
		filterTransaction(transaction) {
			if (!transaction.docChanged) return true;

			const maxDeepness = 6;
			let deepness = 0;

			transaction.doc.descendants((node) => {
				if (node.type.name === 'paragraph') {
					deepness++;
				}
			});

			return deepness <= maxDeepness;
		}
	});

	onMount(() => {
		editor = new Editor({
			element: element,
			content: content,
			extensions: [
				CharacterCount.configure({
					limit: 150
				}),
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				}),
				Placeholder.configure({
					placeholder: placeholder
				})
			],
			onTransaction: () => {
				editor = editor;
			},
			onUpdate: ({ editor }) => {
				content = editor.getJSON();
			},
			onCreate: ({ editor }) => {
				editor.registerPlugin(limitConsecutiveBreaksPlugin);
			},
			editorProps: {
				attributes: {
					class: 'about-me-input'
				}
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});

	$effect(() => {
		if (content && !contentSet) {
			editor.commands.setContent(content);
			contentSet = true;
		}
	});
</script>

<div
	class={[
		'pointer-events-auto bg-main-900 border-[0.5px] border-main-800 placeholder:text-main-400 hocus:bg-main-800 hocus:border-main-600 mt-1.5 transition-colors duration-100 focus:ring-0 [&>.tiptap]:px-3 [&>.tiptap]:py-2 [&>.tiptap]:min-h-[5rem] rounded-sm',
		classes
	]}
	bind:this={element}
></div>
