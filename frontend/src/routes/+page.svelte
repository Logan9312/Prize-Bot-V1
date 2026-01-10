<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	const API_BASE = import.meta.env.VITE_API_URL || '/api';
	let loading = $state(true);

	const siteUrl = 'https://prizebot.dev';
	const title = 'Prize Bot - Discord Auction & Giveaway Bot';
	const description = 'Discord bot for auctions, giveaways, and prize management. Free to use.';
	const imageUrl = `${siteUrl}/icon.png`;

	onMount(() => {
		const unsubscribe = auth.subscribe((state) => {
			loading = state.loading;
			if (!state.loading && state.user) {
				goto('/dashboard');
			}
		});
		return unsubscribe;
	});
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord auction bot, Discord giveaway bot, prize bot, Discord economy bot, auction bot for Discord" />
	<link rel="canonical" href={siteUrl} />

	<!-- Open Graph / Facebook -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content={siteUrl} />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content={imageUrl} />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:url" content={siteUrl} />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={imageUrl} />

	<!-- JSON-LD -->
	{@html `<script type="application/ld+json">${JSON.stringify({
		"@context": "https://schema.org",
		"@type": "SoftwareApplication",
		"name": "Prize Bot",
		"applicationCategory": "SocialNetworkingApplication",
		"operatingSystem": "Discord",
		"description": description,
		"url": siteUrl,
		"offers": {
			"@type": "Offer",
			"price": "0",
			"priceCurrency": "USD"
		}
	})}</script>`}
</svelte:head>

<div class="min-h-screen flex flex-col">
	<!-- Navigation -->
	<nav class="flex items-center justify-between px-4 sm:px-6 lg:px-8 py-4 max-w-7xl mx-auto w-full">
		<div class="flex items-center gap-2">
			<img src="/icon.png" alt="Prize Bot" class="w-8 h-8 rounded-lg" />
			<span class="font-semibold text-text-primary">Prize Bot</span>
		</div>
		<div class="flex items-center gap-6">
			<a href="/docs" class="text-sm text-text-secondary hover:text-text-primary transition-colors">Docs</a>
			<a href="{API_BASE}/auth/discord" class="text-sm text-text-secondary hover:text-text-primary transition-colors">Login</a>
		</div>
	</nav>

	<!-- Hero -->
	<div class="flex-1 flex flex-col items-center justify-center px-4">
		<div class="text-center max-w-xl">
			<!-- Logo -->
			<div class="mb-8">
				<img src="/icon.png" alt="Prize Bot" class="w-20 h-20 rounded-xl mx-auto" />
			</div>

			<!-- Heading -->
			<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
				Prize Bot
			</h1>

			<p class="text-lg text-text-secondary mb-8">
				Auctions, giveaways, and prize claims for your Discord server.
			</p>

			<!-- CTA Buttons -->
			{#if loading}
				<div class="flex justify-center">
					<div class="spinner spinner-lg"></div>
				</div>
			{:else}
				<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
					<a
						href="{API_BASE}/auth/discord"
						class="inline-flex items-center gap-3 px-6 py-3 bg-[#5865F2] text-white font-medium rounded-lg hover:bg-[#4752c4] transition-colors"
					>
						<svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
							<path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028 14.09 14.09 0 0 0 1.226-1.994.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
						</svg>
						Add to Discord
					</a>
					<a href="/docs" class="inline-flex items-center gap-2 px-6 py-3 bg-surface-700 text-text-primary font-medium rounded-lg hover:bg-surface-600 transition-colors border border-surface-600">
						View Documentation
					</a>
				</div>
			{/if}

			<!-- Features -->
			<div class="mt-12 grid grid-cols-3 gap-4 text-center">
				<a href="/docs/auctions" class="group">
					<div class="text-2xl mb-1">
						<svg class="w-6 h-6 mx-auto text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
					<span class="text-sm text-text-muted group-hover:text-text-primary transition-colors">Auctions</span>
				</a>
				<a href="/docs/giveaways" class="group">
					<div class="text-2xl mb-1">
						<svg class="w-6 h-6 mx-auto text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
						</svg>
					</div>
					<span class="text-sm text-text-muted group-hover:text-text-primary transition-colors">Giveaways</span>
				</a>
				<a href="/docs/claims" class="group">
					<div class="text-2xl mb-1">
						<svg class="w-6 h-6 mx-auto text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
					<span class="text-sm text-text-muted group-hover:text-text-primary transition-colors">Claims</span>
				</a>
			</div>
		</div>
	</div>

	<!-- Footer -->
	<footer class="px-4 py-6 text-center text-sm text-text-muted">
		<a href="https://discord.gg/RxP2z5NGtj" target="_blank" rel="noopener" class="hover:text-text-primary transition-colors">Support Server</a>
	</footer>
</div>
