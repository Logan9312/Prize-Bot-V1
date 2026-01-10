<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Discord Auction Bot - Host Auctions in Your Server | Prize Bot';
	const description = 'Create and manage auctions in your Discord server. Features include anti-snipe protection, auction scheduling, bid history, and custom currency.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord auction bot, auction bot for Discord, Discord bidding bot, host auctions Discord, Discord auction system" />
	<link rel="canonical" href="{siteUrl}/docs/auctions" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/auctions" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content="{siteUrl}/icon.png" />

	<!-- JSON-LD -->
	{@html `<script type="application/ld+json">${JSON.stringify({
		"@context": "https://schema.org",
		"@type": "HowTo",
		"name": "How to Host Auctions in Discord",
		"description": description,
		"step": [
			{
				"@type": "HowToStep",
				"name": "Add Prize Bot",
				"text": "Add Prize Bot to your Discord server"
			},
			{
				"@type": "HowToStep",
				"name": "Configure Settings",
				"text": "Use /settings auctions to configure auction channels and options"
			},
			{
				"@type": "HowToStep",
				"name": "Create Auction",
				"text": "Use /auction create to start your first auction"
			}
		]
	})}</script>`}
</svelte:head>

<div class="max-w-4xl">
	<!-- Hero -->
	<div class="mb-10">
		<div class="flex items-center gap-2 text-sm text-text-muted mb-4">
			<a href="/docs" class="hover:text-text-primary transition-colors">Docs</a>
			<span>/</span>
			<span class="text-text-secondary">Auctions</span>
		</div>
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Auctions
		</h1>
		<p class="text-lg text-text-secondary">
			Host auctions in your Discord server with bidding, anti-snipe protection, and scheduling.
		</p>
	</div>

	<!-- Discord Embed Preview -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How It Looks</h2>
		<div class="bg-[#2f3136] rounded-lg p-4 border-l-4 border-accent">
			<div class="flex items-center gap-2 mb-2">
				<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded-full" />
				<span class="font-semibold text-white">Prize Bot</span>
				<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white">BOT</span>
			</div>
			<div class="text-white font-semibold mb-2">Rare Gaming Mouse</div>
			<div class="text-[#dcddde] text-sm space-y-1 mb-3">
				<p>A limited edition gaming mouse with RGB lighting and 16000 DPI sensor.</p>
			</div>
			<div class="grid grid-cols-2 gap-4 text-sm mb-3">
				<div>
					<div class="text-[#72767d] text-xs uppercase">Current Bid</div>
					<div class="text-white font-medium">$150</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">Time Remaining</div>
					<div class="text-white font-medium">2h 34m</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">High Bidder</div>
					<div class="text-accent">@GamerPro</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">Total Bids</div>
					<div class="text-white font-medium">12</div>
				</div>
			</div>
			<div class="flex gap-2">
				<button class="px-4 py-2 bg-accent text-white text-sm rounded font-medium">Place Bid</button>
				<button class="px-4 py-2 bg-[#4f545c] text-white text-sm rounded font-medium">Bid History</button>
			</div>
			<div class="text-[#72767d] text-xs mt-3">Anti-snipe enabled: Bids in last 5 minutes extend by 2 minutes</div>
		</div>
	</div>

	<!-- Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Features</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Anti-Snipe Protection</h3>
				<p class="text-sm text-text-secondary">Automatically extends auction time when bids are placed near the end.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Custom Currency</h3>
				<p class="text-sm text-text-secondary">Use your server's currency system or create custom currency symbols for your auctions.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Buyout Price</h3>
				<p class="text-sm text-text-secondary">Set a "buy it now" price to let users instantly win the auction.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Bid History</h3>
				<p class="text-sm text-text-secondary">Track all bids with timestamps and bidder information.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Auction Queue</h3>
				<p class="text-sm text-text-secondary">Schedule auctions in advance to run automatically. <span class="text-accent">(Premium)</span></p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Dedicated Channels</h3>
				<p class="text-sm text-text-secondary">Each auction gets its own channel for organized bidding and discussion.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/auction create</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Create a new auction.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">item</span>
						<span class="text-text-secondary">Name of the item being auctioned (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">bid</span>
						<span class="text-text-secondary">Starting bid amount (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">duration</span>
						<span class="text-text-secondary">How long the auction runs, e.g. 1h, 30m, 1d (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">description</span>
						<span class="text-text-secondary">Description of the item</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">currency</span>
						<span class="text-text-secondary">One-time currency symbol for this auction</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">use_currency</span>
						<span class="text-text-secondary">Winner pays with their currency balance</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">integer_only</span>
						<span class="text-text-secondary">Only allow whole number bids</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">increment_min</span>
						<span class="text-text-secondary">Minimum bid increment</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">increment_max</span>
						<span class="text-text-secondary">Maximum bid increment</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">buyout</span>
						<span class="text-text-secondary">Price to instantly win the auction</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">target_price</span>
						<span class="text-text-secondary">Hidden minimum price (no winner if not reached)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">category</span>
						<span class="text-text-secondary">Category to create the auction channel in</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">image</span>
						<span class="text-text-secondary">Attach an image to the auction</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">schedule</span>
						<span class="text-text-secondary">Delay before auction starts, e.g. 24h <span class="text-accent">(Premium)</span></span>
					</div>
				</div>
			</div>

			<div class="card">
				<code class="text-accent font-mono">/bid [amount]</code>
				<p class="text-sm text-text-secondary mt-2">Place a bid on the active auction in the current channel.</p>
			</div>

			<div class="card">
				<code class="text-accent font-mono">/auction queue</code>
				<p class="text-sm text-text-secondary mt-2">View scheduled auctions.</p>
			</div>

			<div class="card">
				<code class="text-accent font-mono">/auction end</code>
				<p class="text-sm text-text-secondary mt-2">End the auction in the current channel early.</p>
			</div>

			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/auction edit</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Modify an active auction.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">item</span>
						<span class="text-text-secondary">Change the item name</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">host</span>
						<span class="text-text-secondary">Change the auction host</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">bid</span>
						<span class="text-text-secondary">Change the current bid</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">winner</span>
						<span class="text-text-secondary">Set the current winner</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">extend</span>
						<span class="text-text-secondary">Extend or reduce time (use negative for reduce)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">description</span>
						<span class="text-text-secondary">Change the description</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">currency</span>
						<span class="text-text-secondary">Change the currency symbol</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">currency_side</span>
						<span class="text-text-secondary">Left or right of amount</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">buyout</span>
						<span class="text-text-secondary">Change the buyout price</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">target_price</span>
						<span class="text-text-secondary">Change the hidden target price</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">image</span>
						<span class="text-text-secondary">Change the image</span>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Settings -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Configuration</h2>
		<p class="text-text-secondary mb-4">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/settings auctions</code> to configure server-wide defaults:</p>
		<div class="card">
			<div class="text-sm space-y-2">
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">log_channel</span>
					<span class="text-text-secondary">Channel where auction results are posted when they end</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">category</span>
					<span class="text-text-secondary">Category to create auction channels in</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">alert_role</span>
					<span class="text-text-secondary">Role to ping when auctions start</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">integer_only</span>
					<span class="text-text-secondary">Only allow whole number bids</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">use_currency</span>
					<span class="text-text-secondary">Winner pays with their currency balance</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">snipe_range</span>
					<span class="text-text-secondary">Remaining time to trigger anti-snipe (e.g. 5m)</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">snipe_extension</span>
					<span class="text-text-secondary">Time to add when anti-snipe triggers (e.g. 2m)</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">snipe_limit</span>
					<span class="text-text-secondary">Maximum total time that can be added via anti-snipe</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">snipe_cap</span>
					<span class="text-text-secondary">Maximum auction duration after extensions</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">channel_lock</span>
					<span class="text-text-secondary">Run auction in the same channel as the command</span>
				</div>
				<div class="flex gap-2">
					<span class="text-accent w-32 flex-shrink-0">channel_prefix</span>
					<span class="text-text-secondary">Prefix for auction channel names (default: 💸│)</span>
				</div>
			</div>
		</div>
	</div>

	<!-- Anti-Snipe Details -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Anti-Snipe System</h2>
		<p class="text-text-secondary mb-4">
			Anti-snipe prevents last-second bidding by extending the auction when bids are placed near the end.
		</p>
		<div class="card">
			<h3 class="font-medium text-text-primary mb-3">How it works:</h3>
			<ol class="space-y-2 text-sm text-text-secondary list-decimal list-inside">
				<li>Set a <strong class="text-text-primary">snipe range</strong> (e.g., 5 minutes)</li>
				<li>When a bid is placed within the snipe range, the auction extends</li>
				<li>Set an <strong class="text-text-primary">extension time</strong> (e.g., 2 minutes)</li>
				<li>Optionally set a <strong class="text-text-primary">limit</strong> on total extensions</li>
				<li>Optionally set a <strong class="text-text-primary">cap</strong> on maximum auction duration</li>
			</ol>
		</div>
	</div>

	<!-- CTA -->
	<div class="card bg-accent/10 border-accent/20">
		<div class="text-center">
			<h3 class="font-semibold text-text-primary mb-2">Ready to get started?</h3>
			<p class="text-sm text-text-secondary mb-4">Add Prize Bot to your server.</p>
			<a href="/" class="btn btn-primary">
				Add to Discord
			</a>
		</div>
	</div>
</div>
