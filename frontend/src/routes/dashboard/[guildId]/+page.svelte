<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type GuildStats } from '$lib/api/client';

	$: guildId = $page.params.guildId!;

	let stats: GuildStats | null = null;
	let loading = true;

	onMount(async () => {
		await loadStats();
	});

	$: if (guildId) {
		loadStats();
	}

	async function loadStats() {
		loading = true;
		try {
			stats = await guildsAPI.getStats(guildId);
		} catch (e) {
			console.error('Failed to load stats:', e);
			stats = null;
		} finally {
			loading = false;
		}
	}

	const features = [
		{
			id: 'auctions',
			title: 'Auctions',
			description: 'Configure auction channels, anti-snipe timers, and bid settings',
			href: (id: string) => `/dashboard/${id}/auctions`,
			icon: 'M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z'
		},
		{
			id: 'giveaways',
			title: 'Giveaways',
			description: 'Set up giveaway alerts and winner announcements',
			href: (id: string) => `/dashboard/${id}/giveaways`,
			icon: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7'
		},
		{
			id: 'claims',
			title: 'Claims',
			description: 'Configure ticket categories and staff roles',
			href: (id: string) => `/dashboard/${id}/claims`,
			icon: 'M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z'
		},
		{
			id: 'currency',
			title: 'Currency',
			description: 'Set your server\'s currency symbol and formatting',
			href: (id: string) => `/dashboard/${id}/currency`,
			icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
		},
		{
			id: 'shop',
			title: 'Shop',
			description: 'Configure shop settings and item limits',
			href: (id: string) => `/dashboard/${id}/shop`,
			icon: 'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z'
		}
	];
</script>

<div>
	<!-- Page Header -->
	<div class="mb-6">
		<h1 class="text-xl font-semibold text-text-primary">Server Settings</h1>
		<p class="text-sm text-text-secondary mt-1">Configure Prize Bot features for your server</p>
	</div>

	<!-- Feature Cards -->
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
		{#each features as feature}
			<a href={feature.href(guildId)} class="feature-card group">
				<div class="flex items-start gap-4">
					<div class="w-10 h-10 rounded-lg bg-surface-600 flex items-center justify-center flex-shrink-0 group-hover:bg-accent transition-colors">
						<svg class="w-5 h-5 text-text-secondary group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={feature.icon} />
						</svg>
					</div>
					<div class="flex-1 min-w-0">
						<h3 class="font-medium text-text-primary group-hover:text-accent transition-colors">
							{feature.title}
						</h3>
						<p class="text-sm text-text-secondary mt-1">
							{feature.description}
						</p>
					</div>
				</div>
			</a>
		{/each}
	</div>

	<!-- Stats -->
	<div class="mt-8 pt-6 border-t border-surface-600">
		<h3 class="text-xs font-medium text-text-secondary uppercase tracking-wide mb-4">Quick Stats</h3>
		<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
			<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
				<p class="text-2xl font-semibold text-text-primary">
					{#if loading}
						<span class="animate-pulse">-</span>
					{:else}
						{stats?.active_auctions ?? 0}
					{/if}
				</p>
				<p class="text-xs text-text-muted mt-1">Active Auctions</p>
			</div>
			<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
				<p class="text-2xl font-semibold text-text-primary">
					{#if loading}
						<span class="animate-pulse">-</span>
					{:else}
						{stats?.running_giveaways ?? 0}
					{/if}
				</p>
				<p class="text-xs text-text-muted mt-1">Running Giveaways</p>
			</div>
			<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
				<p class="text-2xl font-semibold text-text-primary">
					{#if loading}
						<span class="animate-pulse">-</span>
					{:else}
						{stats?.open_claims ?? 0}
					{/if}
				</p>
				<p class="text-xs text-text-muted mt-1">Open Claims</p>
			</div>
			<div class="bg-surface-800 border border-surface-600 rounded-lg p-4">
				<p class="text-2xl font-semibold text-text-primary">
					{#if loading}
						<span class="animate-pulse">-</span>
					{:else}
						{stats?.shop_items ?? 0}
					{/if}
				</p>
				<p class="text-xs text-text-muted mt-1">Shop Items</p>
			</div>
		</div>
	</div>
</div>
