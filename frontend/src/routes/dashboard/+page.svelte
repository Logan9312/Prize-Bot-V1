<script lang="ts">
	import { onMount } from 'svelte';
	import { guildsAPI, type Guild } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	let guilds: Guild[] = $state([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const response = await guildsAPI.list();
			guilds = response.guilds || [];
		} catch (err) {
			toast.error('Failed to load servers');
		} finally {
			loading = false;
		}
	});

	function getIconUrl(guild: Guild): string {
		if (guild.icon_url) return guild.icon_url;
		return '';
	}

	const activeGuilds = $derived(guilds.filter(g => g.bot_in));
	const inactiveGuilds = $derived(guilds.filter(g => !g.bot_in));

	function handleAddClick(event: MouseEvent) {
		event.stopPropagation();
	}
</script>

<div>
	<!-- Page Header -->
	<div class="mb-6">
		<h1 class="text-xl font-semibold text-text-primary">Select a Server</h1>
		<p class="text-sm text-text-secondary mt-1">Choose a server to configure Prize Bot</p>
	</div>

	{#if loading}
		<div class="flex justify-center py-16">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else if guilds.length === 0}
		<!-- Empty State -->
		<div class="card text-center py-12 max-w-md mx-auto">
			<svg class="w-12 h-12 mx-auto mb-4 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
			</svg>
			<h3 class="font-medium text-text-primary mb-2">No servers found</h3>
			<p class="text-sm text-text-secondary mb-4">
				Add Prize Bot to a server where you have admin permissions.
			</p>
			<a
				href="https://discord.com/oauth2/authorize?client_id=YOUR_BOT_ID&permissions=8&scope=bot%20applications.commands"
				target="_blank"
				rel="noopener"
				class="btn btn-primary"
			>
				Add Prize Bot
			</a>
		</div>
	{:else}
		<!-- Active Servers -->
		{#if activeGuilds.length > 0}
			<div class="mb-8">
				<div class="flex items-center gap-2 mb-3">
					<span class="text-xs font-medium text-text-secondary uppercase tracking-wide">Active</span>
					<span class="text-xs text-text-muted">({activeGuilds.length})</span>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
					{#each activeGuilds as guild}
						<a href="/dashboard/{guild.id}" class="server-card">
							{#if getIconUrl(guild)}
								<img
									src={getIconUrl(guild)}
									alt={guild.name}
									class="w-10 h-10 rounded-lg"
								/>
							{:else}
								<div class="w-10 h-10 rounded-lg bg-accent flex items-center justify-center text-sm font-medium text-white">
									{guild.name[0].toUpperCase()}
								</div>
							{/if}

							<div class="flex-1 min-w-0">
								<p class="font-medium text-text-primary truncate">{guild.name}</p>
								<div class="flex items-center gap-1.5 mt-0.5">
									<span class="w-1.5 h-1.5 rounded-full bg-status-success"></span>
									<span class="text-xs text-status-success">Active</span>
								</div>
							</div>

							<svg class="w-4 h-4 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
							</svg>
						</a>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Inactive Servers -->
		{#if inactiveGuilds.length > 0}
			<div>
				<div class="flex items-center gap-2 mb-3">
					<span class="text-xs font-medium text-text-secondary uppercase tracking-wide">Needs Setup</span>
					<span class="text-xs text-text-muted">({inactiveGuilds.length})</span>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
					{#each inactiveGuilds as guild}
						<div class="server-card opacity-60">
							{#if getIconUrl(guild)}
								<img
									src={getIconUrl(guild)}
									alt={guild.name}
									class="w-10 h-10 rounded-lg grayscale"
								/>
							{:else}
								<div class="w-10 h-10 rounded-lg bg-surface-600 flex items-center justify-center text-sm font-medium text-text-muted">
									{guild.name[0].toUpperCase()}
								</div>
							{/if}

							<div class="flex-1 min-w-0">
								<p class="font-medium text-text-secondary truncate">{guild.name}</p>
								<div class="flex items-center gap-1.5 mt-0.5">
									<span class="w-1.5 h-1.5 rounded-full bg-status-danger"></span>
									<span class="text-xs text-status-danger">Not installed</span>
								</div>
							</div>

							<a
								href="https://discord.com/oauth2/authorize?client_id=YOUR_BOT_ID&permissions=8&scope=bot%20applications.commands&guild_id={guild.id}"
								target="_blank"
								rel="noopener"
								class="btn btn-secondary text-xs py-1.5 px-3"
								onclick={handleAddClick}
							>
								Add
							</a>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
