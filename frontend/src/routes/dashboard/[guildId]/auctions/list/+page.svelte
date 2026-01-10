<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type AuctionListItem } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import FeatureTabs from '$lib/components/FeatureTabs.svelte';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let auctions: AuctionListItem[] = $state([]);

	onMount(async () => {
		await loadAuctions();
	});

	async function loadAuctions() {
		loading = true;
		try {
			const res = await guildsAPI.listAuctions(guildId);
			auctions = res.auctions || [];
		} catch (e) {
			console.error('Failed to load auctions:', e);
			toast.error('Failed to load auctions');
		} finally {
			loading = false;
		}
	}

	function formatTimeRemaining(endTime: string): string {
		const end = new Date(endTime);
		const now = new Date();
		const diff = end.getTime() - now.getTime();

		if (diff <= 0) return 'Ended';

		const days = Math.floor(diff / (1000 * 60 * 60 * 24));
		const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
		const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

		if (days > 0) return `${days}d ${hours}h`;
		if (hours > 0) return `${hours}h ${minutes}m`;
		return `${minutes}m`;
	}

	function formatCurrency(amount: number, currency: string, side: string): string {
		const curr = currency || '$';
		if (side === 'right') {
			return `${amount.toLocaleString()}${curr}`;
		}
		return `${curr}${amount.toLocaleString()}`;
	}

	function getTimeUrgency(endTime: string): 'urgent' | 'warning' | 'normal' {
		const end = new Date(endTime);
		const now = new Date();
		const diff = end.getTime() - now.getTime();
		const hours = diff / (1000 * 60 * 60);

		if (hours <= 1) return 'urgent';
		if (hours <= 6) return 'warning';
		return 'normal';
	}
</script>

<div class="space-y-6">
	<!-- Header with tabs -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<FeatureTabs feature="auctions" activeTab="list" {guildId} listCount={auctions.length} />

		<button
			onclick={loadAuctions}
			disabled={loading}
			class="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-800 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
		>
			<svg class="w-4 h-4 {loading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
			</svg>
			Refresh
		</button>
	</div>

	<!-- Summary -->
	<p class="text-sm text-text-secondary">
		{auctions.length} active auction{auctions.length !== 1 ? 's' : ''}
	</p>

	<!-- Content -->
	{#if loading}
		<div class="flex justify-center py-12">
			<div class="w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin"></div>
		</div>
	{:else if auctions.length === 0}
		<div class="bg-surface-800 border border-surface-600 rounded-lg text-center py-16 px-4">
			<div class="w-16 h-16 mx-auto mb-4 rounded-full bg-surface-700 flex items-center justify-center">
				<svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
				</svg>
			</div>
			<h3 class="text-lg font-medium text-text-primary mb-2">No active auctions</h3>
			<p class="text-text-secondary max-w-sm mx-auto">
				Start an auction with the <code class="px-1.5 py-0.5 bg-surface-700 rounded text-xs">/auction create</code> command in Discord.
			</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each auctions as auction}
				{@const urgency = getTimeUrgency(auction.end_time)}
				<div class="bg-surface-800 border border-surface-600 rounded-lg p-4 hover:border-surface-500 transition-colors">
					<div class="flex justify-between items-start gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<span
									class="w-2 h-2 rounded-full shrink-0 {urgency === 'urgent' ? 'bg-red-400 animate-pulse' : urgency === 'warning' ? 'bg-yellow-400' : 'bg-green-400'}"
									title="{urgency === 'urgent' ? 'Ending soon!' : urgency === 'warning' ? 'Ending in a few hours' : 'Active'}"
								></span>
								<h3 class="font-medium text-text-primary truncate">{auction.item}</h3>
							</div>
							<div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm">
								<span class="text-text-secondary">
									Bid: <span class="text-accent font-semibold">{formatCurrency(auction.bid, auction.currency, auction.currency_side)}</span>
								</span>
								{#if auction.winner}
									<span class="text-text-muted">
										Leader: <span class="text-text-primary font-mono text-xs">{auction.winner}</span>
									</span>
								{:else}
									<span class="text-text-muted italic">No bids yet</span>
								{/if}
							</div>
						</div>
						<div class="text-right shrink-0">
							<p class="text-xs text-text-muted mb-1">Ends in</p>
							<p class="font-mono text-sm font-medium {urgency === 'urgent' ? 'text-red-400' : urgency === 'warning' ? 'text-yellow-400' : 'text-text-primary'}">
								{formatTimeRemaining(auction.end_time)}
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
