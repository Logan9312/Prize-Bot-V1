<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Prize Bot Premium - Advanced Features for Discord Servers';
	const description = 'Premium features including auction scheduling, bulk claims, and whitelabel bots.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Prize Bot premium, Discord bot premium, auction scheduling, whitelabel Discord bot, bulk claims Discord" />
	<link rel="canonical" href="{siteUrl}/docs/premium" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/premium" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content="{siteUrl}/icon.png" />
</svelte:head>

<div class="max-w-4xl">
	<!-- Hero -->
	<div class="mb-10">
		<div class="flex items-center gap-2 text-sm text-text-muted mb-4">
			<a href="/docs" class="hover:text-text-primary transition-colors">Docs</a>
			<span>/</span>
			<span class="text-text-secondary">Premium</span>
		</div>
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Premium
		</h1>
		<p class="text-lg text-text-secondary">
			Additional features for $5/month.
		</p>
	</div>

	<!-- Pricing Banner -->
	<div class="card bg-gradient-to-r from-accent/20 to-accent/5 border-accent/30 mb-10">
		<div class="flex flex-col sm:flex-row items-center justify-between gap-4">
			<div>
				<div class="flex items-center gap-2 mb-1">
					<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
					</svg>
					<span class="font-semibold text-text-primary">Prize Bot Premium</span>
				</div>
				<p class="text-sm text-text-secondary">All premium features for one simple price</p>
			</div>
			<div class="text-center sm:text-right">
				<div class="text-3xl font-bold text-text-primary">$5<span class="text-lg font-normal text-text-muted">/month</span></div>
				<a href="/pricing" class="text-sm text-accent hover:underline">View pricing details</a>
			</div>
		</div>
	</div>

	<!-- Premium Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">What's Included</h2>
		<div class="space-y-4">
			<!-- Auction Queue -->
			<div class="card">
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
						<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					</div>
					<div>
						<h3 class="font-medium text-text-primary mb-1">Auction Scheduling</h3>
						<p class="text-sm text-text-secondary mb-2">Schedule auctions in advance to run automatically at specific times.</p>
						<code class="text-xs bg-surface-800 px-2 py-1 rounded text-accent">/auction queue</code>
					</div>
				</div>
			</div>

			<!-- Bulk Claims -->
			<div class="card">
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
						<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
						</svg>
					</div>
					<div>
						<h3 class="font-medium text-text-primary mb-1">Bulk Claims</h3>
						<p class="text-sm text-text-secondary mb-2">Give prizes to everyone with a specific role at once.</p>
					</div>
				</div>
			</div>

			<!-- Whitelabel -->
			<div class="card">
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center">
						<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
						</svg>
					</div>
					<div>
						<h3 class="font-medium text-text-primary mb-1">Whitelabel Bot</h3>
						<p class="text-sm text-text-secondary mb-2">Run Prize Bot under your own brand with a custom bot name and avatar.</p>
						<a href="/docs/whitelabel" class="text-sm text-accent hover:underline">Learn more about whitelabeling →</a>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Coming Soon -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Coming Soon</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card border-dashed">
				<div class="flex items-center gap-2 mb-2">
					<span class="text-text-primary font-medium">Auto-Enter Giveaways</span>
					<span class="text-xs bg-surface-600 text-text-muted px-2 py-0.5 rounded">Coming Soon</span>
				</div>
				<p class="text-sm text-text-secondary">Automatically enter users by role into giveaways without requiring manual reactions.</p>
			</div>
			<div class="card border-dashed">
				<div class="flex items-center gap-2 mb-2">
					<span class="text-text-primary font-medium">Custom Thumbnails</span>
					<span class="text-xs bg-surface-600 text-text-muted px-2 py-0.5 rounded">Coming Soon</span>
				</div>
				<p class="text-sm text-text-secondary">Add custom images to your auction embeds for a more visual experience.</p>
			</div>
			<div class="card border-dashed">
				<div class="flex items-center gap-2 mb-2">
					<span class="text-text-primary font-medium">Custom Embed Colors</span>
					<span class="text-xs bg-surface-600 text-text-muted px-2 py-0.5 rounded">Coming Soon</span>
				</div>
				<p class="text-sm text-text-secondary">Customize the color of auction and giveaway embeds to match your server's theme.</p>
			</div>
			<div class="card border-dashed">
				<div class="flex items-center gap-2 mb-2">
					<span class="text-text-primary font-medium">Priority Support</span>
					<span class="text-xs bg-surface-600 text-text-muted px-2 py-0.5 rounded">Coming Soon</span>
				</div>
				<p class="text-sm text-text-secondary">Get faster responses and dedicated support for premium subscribers.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Premium Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<code class="text-accent font-mono">/premium info</code>
				<p class="text-sm text-text-secondary mt-2">Check your premium subscription status and linked servers.</p>
			</div>
			<div class="card">
				<code class="text-accent font-mono">/premium activate</code>
				<p class="text-sm text-text-secondary mt-2">Link your premium subscription to a server to enable premium features.</p>
			</div>
		</div>
	</div>

	<!-- How to Subscribe -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How to Subscribe</h2>
		<div class="card">
			<ol class="space-y-4">
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">1</span>
					<div>
						<h3 class="font-medium text-text-primary">Visit the pricing page</h3>
						<p class="text-sm text-text-secondary mt-1">Go to <a href="/pricing" class="text-accent hover:underline">prizebot.dev/pricing</a> and click "Get Premium".</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">2</span>
					<div>
						<h3 class="font-medium text-text-primary">Complete payment</h3>
						<p class="text-sm text-text-secondary mt-1">Pay securely via Stripe. We accept all major credit cards.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">3</span>
					<div>
						<h3 class="font-medium text-text-primary">Activate on your server</h3>
						<p class="text-sm text-text-secondary mt-1">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/premium activate</code> in Discord to link your subscription.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">4</span>
					<div>
						<h3 class="font-medium text-text-primary">Enjoy premium features</h3>
						<p class="text-sm text-text-secondary mt-1">All premium features are now unlocked on your linked server!</p>
					</div>
				</li>
			</ol>
		</div>
	</div>

	<!-- FAQ -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">FAQ</h2>
		<div class="space-y-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Can I use premium on multiple servers?</h3>
				<p class="text-sm text-text-secondary">Each subscription can be linked to one server at a time. You can switch servers anytime from the dashboard or with /premium activate.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">How do I cancel?</h3>
				<p class="text-sm text-text-secondary">Go to the billing portal in your dashboard to cancel anytime. You'll keep premium access until the end of your billing period.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">What payment methods do you accept?</h3>
				<p class="text-sm text-text-secondary">We accept all major credit and debit cards through Stripe. Apple Pay and Google Pay are also supported.</p>
			</div>
		</div>
	</div>

	<!-- CTA -->
	<div class="card bg-accent/10 border-accent/20">
		<div class="text-center">
			<h3 class="font-semibold text-text-primary mb-2">Interested?</h3>
			<p class="text-sm text-text-secondary mb-4">Premium is $5/month.</p>
			<a href="/pricing" class="btn btn-primary">
				View Pricing
			</a>
		</div>
	</div>
</div>
