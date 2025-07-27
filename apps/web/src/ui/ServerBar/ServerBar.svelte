<script lang="ts">
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import UserLine from 'ui/UserLine/UserLine.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import Input from 'ui/Input/Input.svelte';
</script>

<div class="bg-main-975 border-l-[0.5px] border-l-main-700 w-[19.5rem] overflow-hidden">
	<section class="p-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search" />
	</section>
	<BarSeparator title="Members - 10" />
	<div class="flex flex-col gap-y-2 p-2.5">
		<CollapsibleBox header="Administrator - 2" canCollapse={false} color="#AC393A">
			<UserLine
				avatar="https://i.pinimg.com/736x/14/b6/aa/14b6aaeff07a05e98cf19e50f19ea0d4.jpg"
				name="John Doe"
				color="#AC393A"
				hoverable
			/>
			<UserLine
				avatar="https://i.pinimg.com/736x/14/b6/aa/14b6aaeff07a05e98cf19e50f19ea0d4.jpg"
				name="John Doe"
				color="#AC393A"
				hoverable
			/>
		</CollapsibleBox>
		<CollapsibleBox header="Moderator - 1" canCollapse={false} color="#3958AC">
			<UserLine
				avatar="https://i.pinimg.com/736x/14/b6/aa/14b6aaeff07a05e98cf19e50f19ea0d4.jpg"
				name="John Doe"
				hoverable
				color="#3958AC"
			/>
		</CollapsibleBox>
	</div>
</div>
