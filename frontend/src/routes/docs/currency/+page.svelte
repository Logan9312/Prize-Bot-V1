<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Discord Currency Bot - Server Economy System | Prize Bot';
	const description = 'Economy system for your Discord server with custom currency, user balances, role-based management, and auction integration.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord currency bot, Discord economy bot, server currency Discord, Discord balance system, Discord economy system" />
	<link rel="canonical" href="{siteUrl}/docs/currency" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/currency" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary" />
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
			<span class="text-text-secondary">Currency</span>
		</div>
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Currency
		</h1>
		<p class="text-lg text-text-secondary">
			Server economy with custom currency, user balances, and auction integration.
		</p>
	</div>

	<!-- Discord Embed Preview -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How It Looks</h2>
		<div class="bg-[#2f3136] rounded-lg p-4 border-l-4 border-accent">
			<div class="flex items-center gap-2 mb-3">
				<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded-full" />
				<span class="font-semibold text-white">Prize Bot</span>
				<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white">BOT</span>
			</div>
			<div class="flex items-center gap-4 mb-3">
				<div class="w-16 h-16 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-2xl font-bold">G</div>
				<div>
					<div class="text-white font-semibold text-lg">GamerPro</div>
					<div class="text-[#72767d] text-sm">Member since Jan 2024</div>
				</div>
			</div>
			<div class="bg-[#202225] rounded-lg p-3">
				<div class="text-[#72767d] text-xs uppercase mb-1">Balance</div>
				<div class="text-white text-2xl font-bold">🪙 1,250</div>
			</div>
		</div>
	</div>

	<!-- Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Features</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Custom Currency</h3>
				<p class="text-sm text-text-secondary">Set your own currency symbol (emoji or text) and position (left or right of amount).</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">User Balances</h3>
				<p class="text-sm text-text-secondary">Track individual balances for all server members. View with /profile.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Role Operations</h3>
				<p class="text-sm text-text-secondary">Add, set, or subtract currency for everyone with a specific role at once.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Auction Integration</h3>
				<p class="text-sm text-text-secondary">Currency displays automatically in auction embeds with your custom formatting.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Balance Leaderboard</h3>
				<p class="text-sm text-text-secondary">View all users and their balances sorted by amount.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Bulk Management</h3>
				<p class="text-sm text-text-secondary">Easily manage currency for many users at once with role-based operations.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<code class="text-accent font-mono">/profile [user]</code>
				<p class="text-sm text-text-secondary mt-2">View a user's profile including their currency balance. Shows your own profile if no user specified.</p>
			</div>

			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/currency edit</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Modify currency balances for users or roles.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">action</span>
						<span class="text-text-secondary">Add, set, or subtract from balance</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">amount</span>
						<span class="text-text-secondary">The amount to add/set/subtract</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">user</span>
						<span class="text-text-secondary">Target user (optional if using role)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">role</span>
						<span class="text-text-secondary">Target role for bulk operations (optional if using user)</span>
					</div>
				</div>
			</div>

			<div class="card">
				<code class="text-accent font-mono">/currency list</code>
				<p class="text-sm text-text-secondary mt-2">Display all users' currency balances in the server, sorted by amount.</p>
			</div>
		</div>
	</div>

	<!-- Settings -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Configuration</h2>
		<p class="text-text-secondary mb-4">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/settings currency</code> to configure:</p>
		<div class="card">
			<ul class="space-y-3 text-sm">
				<li class="flex gap-3">
					<span class="text-accent">Symbol</span>
					<span class="text-text-secondary">Your currency symbol (e.g., 🪙, $, coins, credits)</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">Position</span>
					<span class="text-text-secondary">Left (🪙 100) or right (100 🪙) of the amount</span>
				</li>
			</ul>
		</div>
	</div>

	<!-- Examples -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Example Configurations</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Gaming Server</h3>
				<div class="text-2xl font-bold text-accent mb-1">🪙 1,500</div>
				<p class="text-sm text-text-muted">Symbol: 🪙 | Position: Left</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Economy Server</h3>
				<div class="text-2xl font-bold text-accent mb-1">$1,500</div>
				<p class="text-sm text-text-muted">Symbol: $ | Position: Left</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Community Server</h3>
				<div class="text-2xl font-bold text-accent mb-1">1,500 points</div>
				<p class="text-sm text-text-muted">Symbol: points | Position: Right</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Trading Server</h3>
				<div class="text-2xl font-bold text-accent mb-1">1,500 credits</div>
				<p class="text-sm text-text-muted">Symbol: credits | Position: Right</p>
			</div>
		</div>
	</div>

	<!-- Use Cases -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Use Cases</h2>
		<div class="card">
			<ul class="space-y-3 text-sm text-text-secondary">
				<li class="flex gap-3">
					<span class="text-accent">•</span>
					<span><strong class="text-text-primary">Auction bidding</strong> - Members spend their currency to bid on items</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">•</span>
					<span><strong class="text-text-primary">Rewards</strong> - Give currency to active members, event winners, or role holders</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">•</span>
					<span><strong class="text-text-primary">Games</strong> - Use currency as stakes for mini-games and competitions</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">•</span>
					<span><strong class="text-text-primary">Trading</strong> - Let members trade items using your server's economy</span>
				</li>
			</ul>
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
