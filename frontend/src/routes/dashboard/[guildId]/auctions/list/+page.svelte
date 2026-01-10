<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type AuctionListItem } from '$lib/api/client';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let auctions: AuctionListItem[] = $state([]);

	onMount(async () => {
		try {
			const res = await guildsAPI.listAuctions(guildId);
			auctions = res.auctions || [];
		} catch (e) {
			console.error('Failed to load auctions:', e);
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

	function formatCurrency(amount: number, currency: string, side: string): string {
		if (side === 'right') {
			return `${amount}${currency}`;
		}
		return `${currency}${amount}`;
	}
</script>

<div>
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-xl font-semibold text-text-primary">Active Auctions</h1>
			<p class="text-sm text-text-secondary mt-1">
				{auctions.length} active auction{auctions.length !== 1 ? 's' : ''}
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
	{:else if auctions.length === 0}
		<div class="bg-surface-800 border border-surface-600 rounded-lg text-center py-12">
			<svg class="w-12 h-12 mx-auto text-text-muted mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
			</svg>
			<p class="text-text-secondary">No active auctions</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each auctions as auction}
				<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
					<div class="flex justify-between items-start">
						<div class="flex-1 min-w-0">
							<h3 class="font-medium text-text-primary truncate">{auction.item}</h3>
							<p class="text-sm text-text-secondary mt-1">
								Current bid: <span class="text-accent font-medium">{formatCurrency(auction.bid, auction.currency || '$', auction.currency_side || 'left')}</span>
							</p>
							{#if auction.winner}
								<p class="text-xs text-text-muted mt-1">
									Leading: <span class="text-text-secondary">{auction.winner}</span>
								</p>
							{/if}
						</div>
						<div class="text-right ml-4">
							<p class="text-xs text-text-muted">Ends in</p>
							<p class="font-mono text-text-primary">
								{formatTimeRemaining(auction.end_time)}
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
