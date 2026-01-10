<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { Channel } from '$lib/api/client';

	let { value = $bindable(), label = '', type = 'all' }: {
		value?: string;
		label?: string;
		type?: 'text' | 'category' | 'all';
	} = $props();

	// Normalize undefined to empty string for the select
	const selectValue = $derived(value ?? '');

	const id = `channel-select-${Math.random().toString(36).substring(2, 11)}`;
	const channels = getContext<Writable<Channel[]>>('channels');

	const filteredChannels = $derived($channels.filter((ch) => {
		if (type === 'text') return ch.type === 0;
		if (type === 'category') return ch.type === 4;
		return true;
	}));

	const categories = $derived($channels.filter((ch) => ch.type === 4));

	function getChannelsByCategory(parentId: string | null) {
		return filteredChannels.filter((ch) => {
			if (ch.type === 4) return false;
			return parentId ? ch.parent_id === parentId : !ch.parent_id;
		});
	}
</script>

<div>
	{#if label}
		<label for={id} class="label">{label}</label>
	{/if}
	<select {id} value={selectValue} onchange={(e) => value = e.currentTarget.value} class="select">
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
