<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Prize Bot Documentation - Discord Auction & Giveaway Bot';
	const description = 'Documentation for Prize Bot - a Discord bot for auctions, giveaways, prize claims, and server currency.';

	const features = [
		{
			href: '/docs/auctions',
			title: 'Auctions',
			description: 'Host auctions in your Discord server with bidding, anti-snipe protection, and scheduling.',
			icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
		},
		{
			href: '/docs/giveaways',
			title: 'Giveaways',
			description: 'Run giveaways with reaction-based entries, multiple winners, and reroll.',
			icon: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7'
		},
		{
			href: '/docs/claims',
			title: 'Claims & Tickets',
			description: 'Distribute prizes with automatic ticket creation, claim tracking, and logging.',
			icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
		},
		{
			href: '/docs/currency',
			title: 'Currency System',
			description: 'Server economy with custom currency, user balances, and role-based management.',
			icon: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z'
		},
		{
			href: '/docs/premium',
			title: 'Premium Features',
			description: 'Auction scheduling, bulk claims, and custom branding.',
			icon: 'M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z'
		},
		{
			href: '/docs/whitelabel',
			title: 'Whitelabel Bots',
			description: 'Run Prize Bot under your own brand with a custom bot name and avatar.',
			icon: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01'
		}
	];
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord bot, auction bot, giveaway bot, Discord auctions, Discord giveaways, prize bot, Discord economy bot" />
	<link rel="canonical" href="{siteUrl}/docs" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content="{siteUrl}/icon.png" />

	<!-- JSON-LD Structured Data -->
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
		},
		"aggregateRating": {
			"@type": "AggregateRating",
			"ratingValue": "5",
			"ratingCount": "100"
		}
	})}</script>`}
</svelte:head>

<div class="max-w-4xl">
	<!-- Hero -->
	<div class="mb-12">
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Prize Bot Documentation
		</h1>
		<p class="text-lg text-text-secondary mb-6">
			A Discord bot for auctions, giveaways, and prize management.
		</p>
		<div class="flex flex-wrap gap-3">
			<a href="/" class="btn btn-primary">
				<svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
					<path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028 14.09 14.09 0 0 0 1.226-1.994.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
				</svg>
				Add to Discord
			</a>
			<a href="/pricing" class="btn btn-secondary">
				View Pricing
			</a>
		</div>
	</div>

	<!-- Features Grid -->
	<h2 class="text-xl font-semibold text-text-primary mb-6">Features</h2>
	<div class="grid sm:grid-cols-2 gap-4 mb-12">
		{#each features as feature}
			<a href={feature.href} class="feature-card group">
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center group-hover:bg-accent/20 transition-colors">
						<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={feature.icon} />
						</svg>
					</div>
					<div>
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

	<!-- Getting Started -->
	<h2 class="text-xl font-semibold text-text-primary mb-6">Getting Started</h2>
	<div class="card mb-8">
		<ol class="space-y-4">
			<li class="flex gap-4">
				<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">1</span>
				<div>
					<h3 class="font-medium text-text-primary">Add Prize Bot to your server</h3>
					<p class="text-sm text-text-secondary mt-1">Click "Add to Discord" and select your server. Grant the required permissions.</p>
				</div>
			</li>
			<li class="flex gap-4">
				<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">2</span>
				<div>
					<h3 class="font-medium text-text-primary">Configure your settings</h3>
					<p class="text-sm text-text-secondary mt-1">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/settings</code> to set up channels, roles, and preferences.</p>
				</div>
			</li>
			<li class="flex gap-4">
				<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">3</span>
				<div>
					<h3 class="font-medium text-text-primary">Start hosting events</h3>
					<p class="text-sm text-text-secondary mt-1">Create your first auction with <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/auction create</code> or giveaway with <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/giveaway create</code>.</p>
				</div>
			</li>
		</ol>
	</div>

	<!-- Support CTA -->
	<div class="card bg-accent/10 border-accent/20">
		<div class="flex flex-col sm:flex-row items-center gap-4">
			<div class="flex-1">
				<h3 class="font-medium text-text-primary">Need help?</h3>
				<p class="text-sm text-text-secondary mt-1">Join our Discord server for support, feature requests, and community discussion.</p>
			</div>
			<a href="https://discord.gg/RxP2z5NGtj" target="_blank" rel="noopener" class="btn btn-secondary whitespace-nowrap">
				Join Support Server
			</a>
		</div>
	</div>
</div>
