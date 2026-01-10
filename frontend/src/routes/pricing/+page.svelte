<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Prize Bot Pricing - Free & Premium Plans for Discord';
	const description = 'Prize Bot pricing. Free tier with auctions and giveaways. Premium for $5/month adds scheduling, bulk claims, and whitelabeling.';

	const freeFeatures = [
		'Unlimited auctions',
		'Unlimited giveaways',
		'Prize claim system',
		'Currency system',
		'Anti-snipe protection',
		'Bid history tracking',
		'Multiple winners for giveaways',
		'Custom currency symbols',
		'Claim logging',
		'24/7 uptime'
	];

	const premiumFeatures = [
		'Everything in Free',
		'Auction scheduling (queue)',
		'Bulk claims (role-based)',
		'Whitelabel bot (custom branding)',
		'Coming soon: Auto-enter giveaways',
		'Coming soon: Custom thumbnails',
		'Coming soon: Custom embed colors',
		'Priority support'
	];

	const comparisonFeatures = [
		{ name: 'Auctions', free: 'Unlimited', premium: 'Unlimited' },
		{ name: 'Giveaways', free: 'Unlimited', premium: 'Unlimited' },
		{ name: 'Claims & Tickets', free: true, premium: true },
		{ name: 'Currency System', free: true, premium: true },
		{ name: 'Anti-Snipe Protection', free: true, premium: true },
		{ name: 'Bid History', free: true, premium: true },
		{ name: 'Custom Currency Symbol', free: true, premium: true },
		{ name: 'Auction Scheduling', free: false, premium: true },
		{ name: 'Bulk Claims', free: false, premium: true },
		{ name: 'Whitelabel Bot', free: false, premium: true },
		{ name: 'Custom Embed Colors', free: false, premium: 'Coming Soon' },
		{ name: 'Custom Thumbnails', free: false, premium: 'Coming Soon' }
	];
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Prize Bot pricing, Discord bot pricing, auction bot price, giveaway bot cost, Discord bot premium" />
	<link rel="canonical" href="{siteUrl}/pricing" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/pricing" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content="{siteUrl}/icon.png" />

	<!-- JSON-LD Pricing -->
	{@html `<script type="application/ld+json">${JSON.stringify({
		"@context": "https://schema.org",
		"@type": "Product",
		"name": "Prize Bot Premium",
		"description": "Premium subscription for Prize Bot Discord bot",
		"brand": {
			"@type": "Brand",
			"name": "Prize Bot"
		},
		"offers": [
			{
				"@type": "Offer",
				"name": "Free",
				"price": "0",
				"priceCurrency": "USD",
				"description": "Free tier with full auction and giveaway features"
			},
			{
				"@type": "Offer",
				"name": "Premium",
				"price": "5",
				"priceCurrency": "USD",
				"billingIncrement": "P1M",
				"description": "Premium tier with auction scheduling, bulk claims, and whitelabeling"
			}
		]
	})}</script>`}
</svelte:head>

<div class="min-h-screen">
	<!-- Header -->
	<header class="sticky top-0 z-40 bg-surface-900/95 backdrop-blur border-b border-surface-600">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<a href="/" class="flex items-center gap-3">
					<img src="/icon.png" alt="Prize Bot" class="w-8 h-8 rounded-lg" />
					<span class="font-semibold text-text-primary">Prize Bot</span>
				</a>
				<nav class="flex items-center gap-6">
					<a href="/docs" class="text-sm text-text-secondary hover:text-text-primary transition-colors">
						Docs
					</a>
					<a href="/" class="btn btn-primary text-sm !py-2 !px-4 !min-h-0">
						Get Started
					</a>
				</nav>
			</div>
		</div>
	</header>

	<main class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
		<!-- Hero -->
		<div class="text-center mb-12">
			<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
				Pricing
			</h1>
			<p class="text-lg text-text-secondary max-w-2xl mx-auto">
				Free tier has most features. Premium adds scheduling, bulk claims, and whitelabeling.
			</p>
		</div>

		<!-- Pricing Cards -->
		<div class="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto mb-16">
			<!-- Free Plan -->
			<div class="card">
				<div class="mb-6">
					<h2 class="text-xl font-semibold text-text-primary mb-1">Free</h2>
					<p class="text-sm text-text-secondary">Most features included</p>
				</div>
				<div class="mb-6">
					<span class="text-4xl font-bold text-text-primary">$0</span>
					<span class="text-text-muted">/month</span>
				</div>
				<a href="/" class="btn btn-secondary w-full mb-6">
					Get Started
				</a>
				<ul class="space-y-3">
					{#each freeFeatures as feature}
						<li class="flex items-start gap-3 text-sm">
							<svg class="w-5 h-5 text-status-success flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
							</svg>
							<span class="text-text-secondary">{feature}</span>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Premium Plan -->
			<div class="card bg-accent/5 border-accent/30 relative">
				<div class="absolute -top-3 left-1/2 -translate-x-1/2">
					<span class="bg-accent text-white text-xs font-medium px-3 py-1 rounded-full">Most Popular</span>
				</div>
				<div class="mb-6">
					<h2 class="text-xl font-semibold text-text-primary mb-1">Premium</h2>
					<p class="text-sm text-text-secondary">Scheduling, bulk claims, whitelabeling</p>
				</div>
				<div class="mb-6">
					<span class="text-4xl font-bold text-text-primary">$5</span>
					<span class="text-text-muted">/month</span>
				</div>
				<a href="/" class="btn btn-primary w-full mb-6">
					Get Premium
				</a>
				<ul class="space-y-3">
					{#each premiumFeatures as feature}
						<li class="flex items-start gap-3 text-sm">
							<svg class="w-5 h-5 text-accent flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
							</svg>
							<span class="text-text-secondary">{feature}</span>
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<!-- Feature Comparison -->
		<div class="mb-16">
			<h2 class="text-2xl font-bold text-text-primary text-center mb-8">Feature Comparison</h2>
			<div class="overflow-x-auto">
				<table class="w-full max-w-3xl mx-auto">
					<thead>
						<tr class="border-b border-surface-600">
							<th class="text-left py-3 px-4 text-text-primary font-medium">Feature</th>
							<th class="text-center py-3 px-4 text-text-primary font-medium">Free</th>
							<th class="text-center py-3 px-4 text-text-primary font-medium">Premium</th>
						</tr>
					</thead>
					<tbody>
						{#each comparisonFeatures as feature}
							<tr class="border-b border-surface-700">
								<td class="py-3 px-4 text-text-secondary">{feature.name}</td>
								<td class="text-center py-3 px-4">
									{#if feature.free === true}
										<svg class="w-5 h-5 text-status-success mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									{:else if feature.free === false}
										<svg class="w-5 h-5 text-text-muted mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
										</svg>
									{:else}
										<span class="text-text-secondary text-sm">{feature.free}</span>
									{/if}
								</td>
								<td class="text-center py-3 px-4">
									{#if feature.premium === true}
										<svg class="w-5 h-5 text-accent mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									{:else if feature.premium === false}
										<svg class="w-5 h-5 text-text-muted mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
										</svg>
									{:else}
										<span class="text-accent text-sm">{feature.premium}</span>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- FAQ -->
		<div class="max-w-3xl mx-auto mb-16">
			<h2 class="text-2xl font-bold text-text-primary text-center mb-8">Frequently Asked Questions</h2>
			<div class="space-y-4">
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">What's included in free?</h3>
					<p class="text-sm text-text-secondary">Unlimited auctions, giveaways, claims, and currency. Premium adds scheduling, bulk claims, and whitelabeling.</p>
				</div>
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">Can I cancel anytime?</h3>
					<p class="text-sm text-text-secondary">Yes. Cancel from the billing portal. You keep access until the end of your billing period.</p>
				</div>
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">What payment methods do you accept?</h3>
					<p class="text-sm text-text-secondary">Credit/debit cards through Stripe. Apple Pay and Google Pay where available.</p>
				</div>
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">How many servers can I use Premium on?</h3>
					<p class="text-sm text-text-secondary">One server per subscription. You can switch servers from the dashboard.</p>
				</div>
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">Do you offer refunds?</h3>
					<p class="text-sm text-text-secondary">Contact us within 7 days for a full refund.</p>
				</div>
				<div class="card">
					<h3 class="font-medium text-text-primary mb-2">What happens to my data if I downgrade?</h3>
					<p class="text-sm text-text-secondary">All data is preserved. You just lose access to premium features.</p>
				</div>
			</div>
		</div>

		<!-- CTA -->
		<div class="text-center">
			<h2 class="text-2xl font-bold text-text-primary mb-4">Get started</h2>
			<p class="text-text-secondary mb-6 max-w-xl mx-auto">
				Add Prize Bot to your server.
			</p>
			<div class="flex justify-center gap-4">
				<a href="/" class="btn btn-primary">
					Add to Discord
				</a>
				<a href="/docs" class="btn btn-secondary">
					View Documentation
				</a>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer class="border-t border-surface-600 mt-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="flex flex-col sm:flex-row items-center justify-between gap-4">
				<div class="flex items-center gap-3">
					<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded" />
					<span class="text-sm text-text-muted">Prize Bot</span>
				</div>
				<div class="flex items-center gap-6 text-sm text-text-muted">
					<a href="/docs" class="hover:text-text-primary transition-colors">Docs</a>
					<a href="https://discord.gg/RxP2z5NGtj" target="_blank" rel="noopener" class="hover:text-text-primary transition-colors">Support</a>
				</div>
			</div>
		</div>
	</footer>
</div>
