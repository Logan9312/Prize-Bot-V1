<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { Channel } from '$lib/api/client';

	export let value: string = '';
	export let label: string = '';
	export let type: 'text' | 'category' | 'all' = 'all';

	const channels = getContext<Writable<Channel[]>>('channels');

	$: filteredChannels = $channels.filter((ch) => {
		if (type === 'text') return ch.type === 0;
		if (type === 'category') return ch.type === 4;
		return true;
	});

	$: categories = $channels.filter((ch) => ch.type === 4);

	function getChannelsByCategory(parentId: string | null) {
		return filteredChannels.filter((ch) => {
			if (ch.type === 4) return false;
			return parentId ? ch.parent_id === parentId : !ch.parent_id;
		});
	}
</script>

<div>
	{#if label}
		<label class="label">{label}</label>
	{/if}
	<select bind:value class="select">
		<option value="">None</option>
		{#if type === 'category'}
			{#each filteredChannels as channel}
				<option value={channel.id}>📁 {channel.name}</option>
			{/each}
		{:else}
			{#each getChannelsByCategory(null) as channel}
				<option value={channel.id}># {channel.name}</option>
			{/each}
			{#each categories as category}
				<optgroup label={category.name}>
					{#each getChannelsByCategory(category.id) as channel}
						<option value={channel.id}># {channel.name}</option>
					{/each}
				</optgroup>
			{/each}
		{/if}
	</select>
</div>
