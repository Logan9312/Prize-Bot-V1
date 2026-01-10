<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type ClaimListItem } from '$lib/api/client';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let claims: ClaimListItem[] = $state([]);

	onMount(async () => {
		try {
			const res = await guildsAPI.listClaims(guildId);
			claims = res.claims || [];
		} catch (e) {
			console.error('Failed to load claims:', e);
		} finally {
			loading = false;
		}
	});

	function getStatusColor(status: string): string {
		switch (status) {
			case 'pending':
				return 'text-yellow-400';
			case 'claimed':
				return 'text-green-400';
			case 'cancelled':
				return 'text-red-400';
			default:
				return 'text-text-secondary';
		}
	}
</script>

<div>
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-xl font-semibold text-text-primary">Open Claims</h1>
			<p class="text-sm text-text-secondary mt-1">
				{claims.length} open claim{claims.length !== 1 ? 's' : ''}
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
	{:else if claims.length === 0}
		<div class="bg-surface-800 border border-surface-600 rounded-lg text-center py-12">
			<svg class="w-12 h-12 mx-auto text-text-muted mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
			</svg>
			<p class="text-text-secondary">No open claims</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each claims as claim}
				<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
					<div class="flex justify-between items-start">
						<div class="flex-1 min-w-0">
							<h3 class="font-medium text-text-primary truncate">{claim.item}</h3>
							<p class="text-sm text-text-secondary mt-1">
								Winner: <span class="text-text-primary">{claim.winner}</span>
							</p>
							{#if claim.cost > 0}
								<p class="text-xs text-text-muted mt-1">
									Cost: ${claim.cost}
								</p>
							{/if}
						</div>
						<div class="text-right ml-4">
							<p class="text-xs text-text-muted">Status</p>
							<p class="font-medium capitalize {getStatusColor(claim.status)}">
								{claim.status}
							</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
