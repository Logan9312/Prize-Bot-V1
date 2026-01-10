<script lang="ts">
	const siteUrl = 'https://prizebot.dev';
	const title = 'Discord Prize Claim System - Ticket-Based Prize Distribution | Prize Bot';
	const description = 'Prize distribution in your Discord server with automatic ticket creation, claim tracking, inventory management, and logging.';
</script>

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta name="keywords" content="Discord prize claim, Discord ticket system, prize distribution Discord, Discord claim bot, prize management Discord" />
	<link rel="canonical" href="{siteUrl}/docs/claims" />

	<!-- Open Graph -->
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/docs/claims" />
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
			<span class="text-text-secondary">Claims</span>
		</div>
		<h1 class="text-3xl sm:text-4xl font-bold text-text-primary mb-4">
			Claims
		</h1>
		<p class="text-lg text-text-secondary">
			Prize distribution with automatic ticket creation, claim tracking, and logging.
		</p>
	</div>

	<!-- Discord Embed Preview -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">How It Looks</h2>
		<div class="bg-[#2f3136] rounded-lg p-4 border-l-4 border-status-warning">
			<div class="flex items-center gap-2 mb-2">
				<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded-full" />
				<span class="font-semibold text-white">Prize Bot</span>
				<span class="text-xs bg-accent px-1.5 py-0.5 rounded text-white">BOT</span>
			</div>
			<div class="text-white font-semibold mb-2">You Won a Prize!</div>
			<div class="text-[#dcddde] text-sm space-y-1 mb-3">
				<p>Congratulations! You've won <strong>Discord Nitro (1 Month)</strong> from the recent giveaway.</p>
				<p class="mt-2">Click the button below to claim your prize and open a ticket.</p>
			</div>
			<div class="flex gap-2 mb-3">
				<button class="px-4 py-2 bg-status-success text-white text-sm rounded font-medium">Claim Prize</button>
			</div>
			<div class="text-[#72767d] text-xs">Prize ID: #1234</div>
		</div>
	</div>

	<!-- Features -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Features</h2>
		<div class="grid sm:grid-cols-2 gap-4">
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Automatic Tickets</h3>
				<p class="text-sm text-text-secondary">When a user claims a prize, a private ticket channel is created for prize delivery.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Prize Inventory</h3>
				<p class="text-sm text-text-secondary">Track unclaimed prizes per user. Users can view their pending prizes anytime.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Claim Logging</h3>
				<p class="text-sm text-text-secondary">All claims are logged to your designated channel with completion status.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Custom Instructions</h3>
				<p class="text-sm text-text-secondary">Set custom instructions that appear in claim tickets for prize delivery.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Bulk Claims</h3>
				<span class="text-xs bg-accent/20 text-accent px-2 py-0.5 rounded">Premium</span>
				<p class="text-sm text-text-secondary mt-2">Give prizes to everyone with a specific role at once.</p>
			</div>
			<div class="card">
				<h3 class="font-medium text-text-primary mb-2">Completion Workflow</h3>
				<p class="text-sm text-text-secondary">Mark claims as complete when prizes are delivered. Reopen tickets if needed.</p>
			</div>
		</div>
	</div>

	<!-- Commands -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Commands</h2>
		<div class="space-y-4">
			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/claim create</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Create a claimable prize.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">target</span>
						<span class="text-text-secondary">User or role to receive the prize (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">item</span>
						<span class="text-text-secondary">The prize you're giving (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">log_channel</span>
						<span class="text-text-secondary">Output channel for the prize (required)</span>
					</div>
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">description</span>
						<span class="text-text-secondary">Description of the prize</span>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/claim inventory</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">View a user's unclaimed prizes.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">user</span>
						<span class="text-text-secondary">User to view prizes for (required)</span>
					</div>
				</div>
			</div>

			<div class="card">
				<div class="flex items-start justify-between gap-4 mb-2">
					<code class="text-accent font-mono">/claim refresh</code>
				</div>
				<p class="text-sm text-text-secondary mb-3">Resend all claim messages. Does not delete old ones.</p>
				<div class="text-sm space-y-2">
					<div class="flex gap-2">
						<span class="text-text-muted w-32 shrink-0">channel</span>
						<span class="text-text-secondary">Output channel for the claims (required)</span>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Claim Workflow -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Claim Workflow</h2>
		<div class="card">
			<ol class="space-y-4">
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">1</span>
					<div>
						<h3 class="font-medium text-text-primary">Prize Created</h3>
						<p class="text-sm text-text-secondary mt-1">Use /claim create or win an auction/giveaway. User receives a DM with a claim button.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">2</span>
					<div>
						<h3 class="font-medium text-text-primary">User Claims</h3>
						<p class="text-sm text-text-secondary mt-1">User clicks the "Claim Prize" button. A private ticket channel is created.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">3</span>
					<div>
						<h3 class="font-medium text-text-primary">Prize Delivery</h3>
						<p class="text-sm text-text-secondary mt-1">Staff can communicate with the winner in the ticket to deliver the prize.</p>
					</div>
				</li>
				<li class="flex gap-4">
					<span class="shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-medium text-white">4</span>
					<div>
						<h3 class="font-medium text-text-primary">Complete Claim</h3>
						<p class="text-sm text-text-secondary mt-1">Click "Complete" to close the ticket. The claim is logged to your log channel.</p>
					</div>
				</li>
			</ol>
		</div>
	</div>

	<!-- Settings -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Configuration</h2>
		<p class="text-text-secondary mb-4">Use <code class="px-1.5 py-0.5 bg-surface-800 rounded text-accent">/settings claiming</code> to configure:</p>
		<div class="card">
			<ul class="space-y-3 text-sm">
				<li class="flex gap-3">
					<span class="text-accent">Category</span>
					<span class="text-text-secondary">Where claim ticket channels are created</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">Log Channel</span>
					<span class="text-text-secondary">Where completed claims are logged</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">Instructions</span>
					<span class="text-text-secondary">Custom message shown in claim tickets</span>
				</li>
				<li class="flex gap-3">
					<span class="text-accent">Disable Claiming</span>
					<span class="text-text-secondary">Record prizes without creating claim tickets (record-only mode)</span>
				</li>
			</ul>
		</div>
	</div>

	<!-- Ticket Actions -->
	<div class="mb-10">
		<h2 class="text-lg font-semibold text-text-primary mb-4">Ticket Actions</h2>
		<div class="card">
			<div class="space-y-4">
				<div>
					<h3 class="font-medium text-text-primary mb-1">Complete</h3>
					<p class="text-sm text-text-secondary">Mark the claim as complete and close the ticket. The claim is logged.</p>
				</div>
				<div>
					<h3 class="font-medium text-text-primary mb-1">Cancel</h3>
					<p class="text-sm text-text-secondary">Cancel the claim without completing it. The prize remains unclaimed.</p>
				</div>
				<div>
					<h3 class="font-medium text-text-primary mb-1">Reopen</h3>
					<p class="text-sm text-text-secondary">Reopen a closed ticket if more communication is needed.</p>
				</div>
			</div>
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
