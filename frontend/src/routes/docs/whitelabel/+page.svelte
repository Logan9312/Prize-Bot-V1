<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Whitelabel Discord Bot - Custom Branded Bot | Prize Bot';
	const description = 'Run Prize Bot under your own brand with a custom bot name and avatar.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="whitelabel Discord bot, custom Discord bot, branded Discord bot, private Discord bot, custom auction bot" />
	<link rel="canonical" href="{siteUrl}/docs/whitelabel" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/whitelabel" />
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
			<span class="text-text-secondary">Whitelabel</span>
		</div>
		<div class="flex items-center gap-2 mb-4">
			<h1 class="text-3xl sm:text-4xl font-bold text-text-primary">
				Whitelabel Bots
			</h1>
			<span class="text-xs bg-accent/20 text-accent px-2 py-1 rounded font-medium">Premium</span>
		</div>
		<p class="text-lg text-text-secondary">
			Run Prize Bot under your own brand with a custom bot name and avatar.
		</p>
	</div>

	<!-- Preview -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Your Bot, Your Brand</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<div class="text-xs text-text-muted uppercase mb-3">Before (Prize Bot)</div>
				<div class="bg-[#2f3136] rounded-lg p-3">
					<div class="flex items-center gap-2">
						<img src="/icon.png" alt="Prize Bot" class="w-8 h-8 rounded-full" />
						<div>
							<span class="font-semibold text-white">Prize Bot</span>
							<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white ml-1">BOT</span>
						</div>
					</div>
				</div>
			</div>
			<div class="card bg-accent/5 border-accent/20">
				<div class="text-xs text-accent uppercase mb-3">After (Your Bot)</div>
				<div class="bg-[#2f3136] rounded-lg p-3">
					<div class="flex items-center gap-2">
						<div class="w-8 h-8 rounded-full bg-linear-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white font-bold text-sm">Y</div>
						<div>
							<span class="font-semibold text-white">Your Server Bot</span>
							<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white ml-1">BOT</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- How It Works -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How It Works</h2>
		<div class="card">
			<ol class="space-y-4">
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">1</span>
					<div>
						<h3 class="font-medium text-text-primary">Create a Discord Application</h3>
						<p class="text-sm text-text-secondary mt-1">Go to the <a href="https://discord.com/developers/applications" target="_blank" rel="noopener" class="text-accent hover:underline">Discord Developer Portal</a> and create a new application with your desired name and avatar.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">2</span>
					<div>
						<h3 class="font-medium text-text-primary">Get your bot token</h3>
						<p class="text-sm text-text-secondary mt-1">In the Bot section, reset and copy your bot token. Keep this secret!</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">3</span>
					<div>
						<h3 class="font-medium text-text-primary">Add your token to Prize Bot</h3>
						<p class="text-sm text-text-secondary mt-1">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/whitelabel add_token</code> to securely add your bot token.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">4</span>
					<div>
						<h3 class="font-medium text-text-primary">Invite your bot</h3>
						<p class="text-sm text-text-secondary mt-1">Use the OAuth2 URL generator in the Developer Portal to create an invite link for your custom bot.</p>
					</div>
				</li>
			</ol>
		</div>
	</div>

	<!-- Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">What You Get</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Custom Identity</h3>
				<p class="text-sm text-text-secondary">Your bot uses its own name, avatar, and status - fully customizable in Discord's Developer Portal.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">All Features Included</h3>
				<p class="text-sm text-text-secondary">Your whitelabel bot has access to all Prize Bot features including premium features.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Secure Token Storage</h3>
				<p class="text-sm text-text-secondary">Your bot token is encrypted and stored securely. Only you can access it.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Automatic Updates</h3>
				<p class="text-sm text-text-secondary">Your whitelabel bot automatically receives all Prize Bot updates and new features.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<code class="text-accent font-mono">/whitelabel add_token [token]</code>
				<p class="text-sm text-text-secondary mt-2">Add or update your custom bot token. The token is encrypted and stored securely.</p>
			</div>
		</div>
	</div>

	<!-- Setup Guide -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Detailed Setup Guide</h2>

		<div class="space-y-6">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-3">Step 1: Create Discord Application</h3>
				<ol class="space-y-2 text-sm text-text-secondary list-decimal list-inside">
					<li>Go to <a href="https://discord.com/developers/applications" target="_blank" rel="noopener" class="text-accent hover:underline">discord.com/developers/applications</a></li>
					<li>Click "New Application"</li>
					<li>Enter your desired bot name (e.g., "Your Server Bot")</li>
					<li>Upload a custom icon in the General Information tab</li>
				</ol>
			</div>

			<div class="card">
				<h3 class="font-medium text-text-primary mb-3">Step 2: Configure Bot Settings</h3>
				<ol class="space-y-2 text-sm text-text-secondary list-decimal list-inside">
					<li>Go to the "Bot" tab in the sidebar</li>
					<li>Click "Reset Token" and copy the new token</li>
					<li>Enable these Privileged Gateway Intents:
						<ul class="list-disc list-inside ml-4 mt-1">
							<li>Server Members Intent</li>
							<li>Message Content Intent</li>
						</ul>
					</li>
					<li>Disable "Public Bot" if you want only you to be able to add it</li>
				</ol>
			</div>

			<div class="card">
				<h3 class="font-medium text-text-primary mb-3">Step 3: Add Token to Prize Bot</h3>
				<ol class="space-y-2 text-sm text-text-secondary list-decimal list-inside">
					<li>In Discord, use the Prize Bot (not your new bot yet)</li>
					<li>Run <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/whitelabel add_token</code></li>
					<li>Paste your bot token when prompted</li>
					<li>Prize Bot will validate and securely store your token</li>
				</ol>
			</div>

			<div class="card">
				<h3 class="font-medium text-text-primary mb-3">Step 4: Generate Invite Link</h3>
				<ol class="space-y-2 text-sm text-text-secondary list-decimal list-inside">
					<li>Go to the "OAuth2" → "URL Generator" tab</li>
					<li>Select the "bot" and "applications.commands" scopes</li>
					<li>Select the permissions your bot needs (Administrator recommended)</li>
					<li>Copy the generated URL and use it to invite your bot</li>
				</ol>
			</div>
		</div>
	</div>

	<!-- FAQ -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">FAQ</h2>
		<div class="space-y-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Is my token secure?</h3>
				<p class="text-sm text-text-secondary">Yes, your bot token is encrypted using industry-standard encryption before being stored. Only the Prize Bot system can decrypt it to run your bot.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Can I use the same bot on multiple servers?</h3>
				<p class="text-sm text-text-secondary">Yes! Once your whitelabel bot is set up, you can invite it to any servers where you have the "Manage Server" permission.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">What happens if I cancel premium?</h3>
				<p class="text-sm text-text-secondary">Your whitelabel bot will stop working when your premium subscription ends. You can continue using the main Prize Bot with free features.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Can I change my bot's name and avatar later?</h3>
				<p class="text-sm text-text-secondary">Yes, you can update your bot's name, avatar, and other settings anytime in the Discord Developer Portal.</p>
			</div>
		</div>
	</div>

	<!-- CTA -->
	<div class="card bg-accent/10 border-accent/20">
		<div class="text-center">
			<h3 class="font-semibold text-text-primary mb-2">Interested?</h3>
			<p class="text-sm text-text-secondary mb-4">Whitelabeling is included with Premium ($5/month).</p>
			<a href="/pricing" class="btn btn-primary">
				Get Premium
			</a>
		</div>
	</div>
</div>
