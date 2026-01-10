<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type GiveawayListItem } from '$lib/api/client';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let giveaways: GiveawayListItem[] = $state([]);

	onMount(async () => {
		try {
			const res = await guildsAPI.listGiveaways(guildId);
			giveaways = res.giveaways || [];
		} catch (e) {
			console.error('Failed to load giveaways:', e);
		} finally {
			loading = false;
		}
	});

	function formatTimeRemaining(endTime: string): string {
		const end = new Date(endTime);
		const now = new Date();
		const diff = end.getTime() - now.getTime();

		if (diff <= 0) return 'Ended';

		const hours = Math.floor(diff / (1000 * 60 * 60));
		const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

		if (hours > 24) {
			const days = Math.floor(hours / 24);
			return `${days}d ${hours % 24}h`;
		}
		return `${hours}h ${minutes}m`;
	}
</script>

<div>
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-xl font-semibold text-text-primary">Running Giveaways</h1>
			<p class="text-sm text-text-secondary mt-1">
				{giveaways.length} running giveaway{giveaways.length !== 1 ? 's' : ''}
			</p>
		</div>
		<a href="/dashboard/{guildId}" class="text-sm text-text-secondary hover:text-text-primary transition-colors">
			Back to Overview
		</a>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin"></div>
		</div>
	{:else if giveaways.length === 0}
		<div class="bg-surface-800 border border-surface-600 rounded-lg text-center py-12">
			<svg class="w-12 h-12 mx-auto text-text-muted mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
			</svg>
			<p class="text-text-secondary">No running giveaways</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each giveaways as giveaway}
				<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
					<div class="flex justify-between items-start">
						<div class="flex-1 min-w-0">
							<h3 class="font-medium text-text-primary truncate">{giveaway.item}</h3>
							<p class="text-sm text-text-secondary mt-1">
								{giveaway.winners} winner{giveaway.winners !== 1 ? 's' : ''}
							</p>
						</div>
						<div class="text-right ml-4">
							<p class="text-xs text-text-muted">Ends in</p>
							<p class="font-mono text-text-primary">
								{formatTimeRemaining(giveaway.end_time)}
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
