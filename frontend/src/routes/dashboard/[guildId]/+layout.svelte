<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, type Channel, type Role } from '$lib/api/client';
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import MobileNav from '$lib/components/MobileNav.svelte';

	$: guildId = $page.params.guildId!;

	const channels = writable<Channel[]>([]);
	const roles = writable<Role[]>([]);

	setContext('channels', channels);
	setContext('roles', roles);

	onMount(async () => {
		try {
			const [channelsRes, rolesRes] = await Promise.all([
				guildsAPI.getChannels(guildId),
				guildsAPI.getRoles(guildId)
			]);
			channels.set(channelsRes.channels || []);
			roles.set(rolesRes.roles || []);
		} catch {
			// Will be empty, which is fine
		}
	});

	const navItems = [
		{ href: `/dashboard/${guildId}`, label: 'Overview', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
		{ href: `/dashboard/${guildId}/auctions`, label: 'Auctions', icon: 'M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z' },
		{ href: `/dashboard/${guildId}/claims`, label: 'Claims', icon: 'M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z' },
		{ href: `/dashboard/${guildId}/giveaways`, label: 'Giveaways', icon: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7' },
		{ href: `/dashboard/${guildId}/currency`, label: 'Currency', icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
		{ href: `/dashboard/${guildId}/shop`, label: 'Shop', icon: 'M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z' }
	];

	$: currentPath = $page.url.pathname;
</script>

<div class="flex flex-col lg:flex-row gap-4 lg:gap-6">
	<!-- Mobile: Hamburger + Drawer -->
	<div class="lg:hidden mb-4">
		<MobileNav {navItems} {currentPath} {guildId} />
	</div>

	<!-- Desktop: Fixed Sidebar -->
	<nav class="hidden lg:block w-44 flex-shrink-0">
		<a
			href="/dashboard"
			class="flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary transition-colors mb-5"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
			Back to servers
		</a>

		<div class="space-y-1">
			{#each navItems as item}
				<a
					href={item.href}
					class="nav-item {currentPath === item.href ? 'nav-item-active' : ''}"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
					</svg>
					{item.label}
				</a>
			{/each}
		</div>
	</nav>

	<!-- Content -->
	<div class="flex-1">
		<slot />
	</div>
</div>
