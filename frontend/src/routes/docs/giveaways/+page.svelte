<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Discord Giveaway Bot - Host Giveaways in Your Server | Prize Bot';
	const description = 'Run giveaways in your Discord server with reaction-based entries, multiple winners, reroll, and customizable duration.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord giveaway bot, giveaway bot for Discord, Discord giveaway system, host giveaways Discord, Discord raffle bot" />
	<link rel="canonical" href="{siteUrl}/docs/giveaways" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/giveaways" />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content="{siteUrl}/icon.png" />
	<meta property="og:site_name" content="Prize Bot" />

	<!-- Twitter -->
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content="{siteUrl}/icon.png" />

	<!-- JSON-LD -->
	{@html `<script type="application/ld+json">${JSON.stringify({
		"@context": "https://schema.org",
		"@type": "HowTo",
		"name": "How to Host Giveaways in Discord",
		"description": description,
		"step": [
			{
				"@type": "HowToStep",
				"name": "Add Prize Bot",
				"text": "Add Prize Bot to your Discord server"
			},
			{
				"@type": "HowToStep",
				"name": "Create Giveaway",
				"text": "Use /giveaway create to start a new giveaway"
			},
			{
				"@type": "HowToStep",
				"name": "Users Enter",
				"text": "Members react with the gift emoji to enter"
			},
			{
				"@type": "HowToStep",
				"name": "Winners Selected",
				"text": "Winners are randomly selected when the timer ends"
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
			<span class="text-text-secondary">Giveaways</span>
		</div>
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Giveaways
		</h1>
		<p class="text-lg text-text-secondary">
			Run giveaways in your Discord server with reaction-based entries and automatic winner selection.
		</p>
	</div>

	<!-- Discord Embed Preview -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How It Looks</h2>
		<div class="bg-[#2f3136] rounded-lg p-4 border-l-4 border-status-success">
			<div class="flex items-center gap-2 mb-2">
				<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded-full" />
				<span class="font-semibold text-white">Prize Bot</span>
				<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white">BOT</span>
			</div>
			<div class="text-white font-semibold mb-2">Nitro Giveaway</div>
			<div class="text-[#dcddde] text-sm space-y-1 mb-3">
				<p>One lucky winner will receive a month of Discord Nitro!</p>
			</div>
			<div class="grid grid-cols-2 gap-4 text-sm mb-3">
				<div>
					<div class="text-[#72767d] text-xs uppercase">Prize</div>
					<div class="text-white font-medium">1x Discord Nitro</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">Ends In</div>
					<div class="text-white font-medium">23h 45m</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">Winners</div>
					<div class="text-white font-medium">1</div>
				</div>
				<div>
					<div class="text-[#72767d] text-xs uppercase">Entries</div>
					<div class="text-white font-medium">47</div>
				</div>
			</div>
			<div class="flex items-center gap-2 text-sm text-[#dcddde]">
				<span class="text-xl">🎁</span>
				<span>React to enter!</span>
			</div>
			<div class="text-[#72767d] text-xs mt-3">Hosted by @Admin</div>
		</div>
	</div>

	<!-- Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Features</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Reaction Entry</h3>
				<p class="text-sm text-text-secondary">Members react with the 🎁 emoji to enter.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Multiple Winners</h3>
				<p class="text-sm text-text-secondary">Select any number of winners for your giveaway.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Reroll Winners</h3>
				<p class="text-sm text-text-secondary">Winner didn't respond? Easily reroll to select new winners.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Custom Duration</h3>
				<p class="text-sm text-text-secondary">Set giveaways to run for minutes, hours, or days.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Alert Notifications</h3>
				<p class="text-sm text-text-secondary">Ping a role when giveaways start to maximize participation.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Image Support</h3>
				<p class="text-sm text-text-secondary">Attach images to showcase your prizes visually.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/giveaway create</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Start a new giveaway.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">item</span>
						<span class="text-text-secondary">What you're giving away (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">winners</span>
						<span class="text-text-secondary">Number of winners to select (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">duration</span>
						<span class="text-text-secondary">How long the giveaway runs, e.g. 1d, 12h (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">channel</span>
						<span class="text-text-secondary">Channel to post the giveaway in</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">description</span>
						<span class="text-text-secondary">Description of the giveaway item</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">image</span>
						<span class="text-text-secondary">Attach an image to the giveaway</span>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/giveaway delete</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Cancel and delete a giveaway.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">message_id</span>
						<span class="text-text-secondary">Message ID of the giveaway to delete (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 flex-shrink-0">channel</span>
						<span class="text-text-secondary">Channel containing the giveaway</span>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="flex items-center gap-2 mb-2">
					<span class="text-text-primary font-medium">Reroll Button</span>
				</div>
				<p class="text-sm text-text-secondary">After a giveaway ends, use the "Reroll" button on the message to select new winners.</p>
			</div>
		</div>
	</div>

	<!-- Settings -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Configuration</h2>
		<p class="text-text-secondary mb-4">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/settings giveaways</code> to configure:</p>
		<div class="card">
			<ul class="space-y-3 text-sm">
				<li class="flex gap-3">
					<span class="text-accent">Log Channel</span>
					<span class="text-text-secondary">Where giveaway results are posted when they end</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">Alert Role</span>
					<span class="text-text-secondary">Role to ping when new giveaways start</span>
				</li>
			</ul>
		</div>
	</div>

	<!-- How it works -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How Giveaways Work</h2>
		<div class="card">
			<ol class="space-y-4">
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">1</span>
					<div>
						<h3 class="font-medium text-text-primary">Create the giveaway</h3>
						<p class="text-sm text-text-secondary mt-1">Use /giveaway create with your prize details and duration.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">2</span>
					<div>
						<h3 class="font-medium text-text-primary">Members enter</h3>
						<p class="text-sm text-text-secondary mt-1">Users react with the 🎁 emoji to enter the giveaway.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">3</span>
					<div>
						<h3 class="font-medium text-text-primary">Timer expires</h3>
						<p class="text-sm text-text-secondary mt-1">When the duration ends, Prize Bot automatically selects winners.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">4</span>
					<div>
						<h3 class="font-medium text-text-primary">Winners announced</h3>
						<p class="text-sm text-text-secondary mt-1">Winners are mentioned in the channel and logged to your log channel.</p>
					</div>
				</li>
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
