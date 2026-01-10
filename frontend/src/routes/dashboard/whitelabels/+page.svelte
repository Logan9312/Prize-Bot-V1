<script lang="ts">
	import { onMount } from 'svelte';
	import {
		whitelabelAPI,
		type Whitelabel,
		type ValidateTokenResponse
	} from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	let loading = $state(true);
	let whitelabels = $state<Whitelabel[]>([]);
	let isAdmin = $state(false);
	let showModal = $state(false);
	let guideExpanded = $state(false);

	// Modal state
	let tokenInput = $state('');
	let validating = $state(false);
	let saving = $state(false);
	let validationResult = $state<ValidateTokenResponse | null>(null);
	let showToken = $state(false);

	// Delete confirmation
	let deleteConfirmId = $state<string | null>(null);
	let deleteConfirmInput = $state('');
	let deleting = $state(false);

	onMount(async () => {
		await loadWhitelabels();
	});

	async function loadWhitelabels() {
		loading = true;
		try {
			const response = await whitelabelAPI.list();
			whitelabels = response.whitelabels || [];
			isAdmin = response.is_admin;
		} catch (err) {
			toast.error('Failed to load whitelabels');
		} finally {
			loading = false;
		}
	}

	async function validateToken() {
		if (!tokenInput.trim()) return;
		validating = true;
		validationResult = null;
		try {
			validationResult = await whitelabelAPI.validate(tokenInput);
		} catch (err) {
			validationResult = { valid: false, error: 'Validation failed' };
		} finally {
			validating = false;
		}
	}

	async function saveWhitelabel() {
		if (!validationResult?.valid) return;
		saving = true;
		try {
			await whitelabelAPI.create(tokenInput);
			toast.success('Whitelabel added! Restart required to activate.');
			closeModal();
			await loadWhitelabels();
		} catch (err: any) {
			toast.error(err.message || 'Failed to save whitelabel');
		} finally {
			saving = false;
		}
	}

	async function deleteWhitelabel(botId: string) {
		deleting = true;
		try {
			await whitelabelAPI.delete(botId);
			toast.success('Whitelabel deleted');
			closeDeleteModal();
			await loadWhitelabels();
		} catch (err: any) {
			toast.error(err.message || 'Failed to delete whitelabel');
		} finally {
			deleting = false;
		}
	}

	function openDeleteModal(botId: string) {
		deleteConfirmId = botId;
		deleteConfirmInput = '';
	}

	function closeDeleteModal() {
		deleteConfirmId = null;
		deleteConfirmInput = '';
	}

	function openModal() {
		showModal = true;
		tokenInput = '';
		validationResult = null;
		showToken = false;
	}

	function closeModal() {
		showModal = false;
		tokenInput = '';
		validationResult = null;
		showToken = false;
	}

	function getBotInitial(name: string): string {
		return name?.charAt(0)?.toUpperCase() || '?';
	}
</script>

<div class="max-w-3xl mx-auto">
	<!-- Page Header -->
	<div class="mb-8">
		<a
			href="/dashboard"
			class="flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary transition-colors mb-4"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
			Back to servers
		</a>
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-semibold text-text-primary">Whitelabels</h1>
				<p class="text-sm text-text-secondary mt-1">Manage your custom bot instances</p>
			</div>
			<button onclick={openModal} class="btn btn-primary">
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
				</svg>
				Add Whitelabel
			</button>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-16">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else}
		<div class="space-y-6">
			<!-- Setup Guide Card -->
			<div class="card overflow-hidden">
				<button
					onclick={() => (guideExpanded = !guideExpanded)}
					class="w-full flex items-center justify-between text-left"
				>
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-lg bg-accent/20 flex items-center justify-center">
							<svg class="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
						</div>
						<div>
							<h2 class="text-lg font-medium text-text-primary">Setup Guide</h2>
							<p class="text-sm text-text-secondary">Learn how to create and configure your custom bot</p>
						</div>
					</div>
					<svg
						class="w-5 h-5 text-text-muted transition-transform duration-200 {guideExpanded ? 'rotate-180' : ''}"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
					</svg>
				</button>

				{#if guideExpanded}
					<div class="mt-6 pt-6 border-t border-surface-600">
						<div class="space-y-6">
							<!-- Step 1 -->
							<div class="flex gap-4">
								<div class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-bold text-white">1</div>
								<div class="flex-1 pt-1">
									<h3 class="font-medium text-text-primary mb-1">Create a Discord Application</h3>
									<p class="text-sm text-text-secondary mb-2">
										Go to the <a href="https://discord.com/developers/applications" target="_blank" rel="noopener" class="text-accent hover:underline">Discord Developer Portal</a> and click "New Application". Give your bot a name and save.
									</p>
								</div>
							</div>

							<!-- Step 2 -->
							<div class="flex gap-4">
								<div class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-bold text-white">2</div>
								<div class="flex-1 pt-1">
									<h3 class="font-medium text-text-primary mb-1">Get Your Bot Token</h3>
									<p class="text-sm text-text-secondary mb-2">
										Navigate to the "Bot" section in the sidebar. Click "Reset Token" to generate a new token, then copy it. <span class="text-status-warning">Keep this token secret!</span>
									</p>
								</div>
							</div>

							<!-- Step 3 -->
							<div class="flex gap-4">
								<div class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-bold text-white">3</div>
								<div class="flex-1 pt-1">
									<h3 class="font-medium text-text-primary mb-1">Configure Bot Settings</h3>
									<p class="text-sm text-text-secondary mb-2">
										In the Bot section, enable <strong class="text-text-primary">Message Content Intent</strong> under "Privileged Gateway Intents". Optionally disable "Public Bot" if you want only you to add it.
									</p>
								</div>
							</div>

							<!-- Step 4 -->
							<div class="flex gap-4">
								<div class="flex-shrink-0 w-8 h-8 rounded-full bg-accent flex items-center justify-center text-sm font-bold text-white">4</div>
								<div class="flex-1 pt-1">
									<h3 class="font-medium text-text-primary mb-1">Invite Your Bot</h3>
									<p class="text-sm text-text-secondary mb-2">
										Go to OAuth2 &rarr; URL Generator. Select "bot" and "applications.commands" scopes, then select "Administrator" permission. Use the generated URL to invite your bot to your server.
									</p>
								</div>
							</div>

							<!-- Step 5 -->
							<div class="flex gap-4">
								<div class="flex-shrink-0 w-8 h-8 rounded-full bg-status-success flex items-center justify-center text-sm font-bold text-white">5</div>
								<div class="flex-1 pt-1">
									<h3 class="font-medium text-text-primary mb-1">Add Token Here</h3>
									<p class="text-sm text-text-secondary">
										Click the "Add Whitelabel" button above, paste your bot token, and save. Your bot will run with all Prize Bot features!
									</p>
								</div>
							</div>
						</div>

						<!-- Important Notes -->
						<div class="mt-6 p-4 bg-status-warning/10 border border-status-warning/30 rounded-lg">
							<div class="flex gap-3">
								<svg class="w-5 h-5 text-status-warning flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
								</svg>
								<div class="text-sm">
									<p class="font-medium text-status-warning mb-1">Important Notes</p>
									<ul class="text-text-secondary space-y-1">
										<li>Never share your bot token publicly</li>
										<li>Whitelabels require an active Premium subscription</li>
										<li>After adding a token, a server restart is required to activate</li>
									</ul>
								</div>
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- Whitelabels List -->
			{#if whitelabels.length === 0}
				<!-- Empty State -->
				<div class="card text-center py-12">
					<div class="w-16 h-16 rounded-full bg-surface-600 flex items-center justify-center mx-auto mb-4">
						<svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
						</svg>
					</div>
					<h3 class="text-lg font-medium text-text-primary mb-2">No Whitelabels Yet</h3>
					<p class="text-text-secondary mb-6 max-w-md mx-auto">
						Create your first custom bot instance to run Prize Bot with your own branding.
					</p>
					<button onclick={openModal} class="btn btn-primary">
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
						Add Your First Whitelabel
					</button>
				</div>
			{:else}
				<div class="space-y-3">
					{#each whitelabels as bot (bot.bot_id)}
						<div class="card flex items-center gap-4 hover:bg-surface-600/50 transition-colors">
							<!-- Bot Avatar -->
							{#if bot.bot_avatar}
								<img
									src={bot.bot_avatar}
									alt={bot.bot_name}
									class="w-12 h-12 rounded-full bg-surface-600"
								/>
							{:else}
								<div class="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center text-accent font-bold text-lg">
									{getBotInitial(bot.bot_name)}
								</div>
							{/if}

							<!-- Bot Info -->
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-2">
									<h3 class="font-medium text-text-primary truncate">
										{bot.bot_name || 'Unknown Bot'}
									</h3>
									<span class="px-2 py-0.5 rounded text-xs font-medium bg-status-success/20 text-status-success">
										Active
									</span>
								</div>
								<div class="flex items-center gap-4 mt-1 text-sm text-text-secondary">
									<span class="font-mono text-xs">{bot.bot_id}</span>
									{#if isAdmin}
										<span class="text-text-muted">Owner: {bot.user_id}</span>
									{/if}
									{#if bot.created_at}
										<span class="text-text-muted">Added {bot.created_at}</span>
									{/if}
								</div>
							</div>

							<!-- Actions -->
							<div class="flex items-center gap-2">
								<button
									onclick={() => openDeleteModal(bot.bot_id)}
									class="btn btn-ghost py-1.5 px-3 min-h-0 text-sm text-status-danger hover:bg-status-danger/10"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
									Delete
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Add Whitelabel Modal -->
{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<!-- Backdrop -->
		<button
			onclick={closeModal}
			class="absolute inset-0 bg-black/60 backdrop-blur-sm"
			aria-label="Close modal"
		></button>

		<!-- Modal -->
		<div class="relative bg-surface-800 border border-surface-600 rounded-xl w-full max-w-lg shadow-2xl">
			<!-- Header -->
			<div class="flex items-center justify-between p-5 border-b border-surface-600">
				<h2 class="text-lg font-semibold text-text-primary">Add Whitelabel</h2>
				<button onclick={closeModal} class="p-1 rounded-lg hover:bg-surface-600 transition-colors" aria-label="Close modal">
					<svg class="w-5 h-5 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<!-- Content -->
			<div class="p-5 space-y-5">
				<!-- Token Input -->
				<div>
					<label for="token" class="label">Bot Token</label>
					<div class="relative">
						<input
							id="token"
							type={showToken ? 'text' : 'password'}
							bind:value={tokenInput}
							placeholder="MTAwNzQyMjQ3MzYyNjczMDUxNg.GYviWS..."
							class="input pr-20 font-mono text-sm"
						/>
						<button
							type="button"
							onclick={() => (showToken = !showToken)}
							class="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-primary transition-colors"
						>
							{#if showToken}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
								</svg>
							{:else}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
								</svg>
							{/if}
						</button>
					</div>
					<p class="text-xs text-text-muted mt-2">
						Paste your Discord bot token from the Developer Portal
					</p>
				</div>

				<!-- Validate Button -->
				<button
					onclick={validateToken}
					disabled={!tokenInput.trim() || validating}
					class="btn btn-secondary w-full"
				>
					{#if validating}
						<span class="spinner w-4 h-4"></span>
						Validating...
					{:else}
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						Validate Token
					{/if}
				</button>

				<!-- Validation Result -->
				{#if validationResult}
					{#if validationResult.valid}
						<div class="p-4 bg-status-success/10 border border-status-success/30 rounded-lg">
							<div class="flex items-center gap-4">
								{#if validationResult.bot_avatar}
									<img
										src={validationResult.bot_avatar}
										alt={validationResult.bot_name}
										class="w-14 h-14 rounded-full"
									/>
								{:else}
									<div class="w-14 h-14 rounded-full bg-accent/20 flex items-center justify-center text-accent font-bold text-xl">
										{getBotInitial(validationResult.bot_name || '')}
									</div>
								{/if}
								<div>
									<div class="flex items-center gap-2">
										<svg class="w-5 h-5 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
										<span class="font-medium text-status-success">Token Valid</span>
									</div>
									<p class="text-lg font-medium text-text-primary mt-1">{validationResult.bot_name}</p>
									<p class="text-sm text-text-secondary font-mono">{validationResult.bot_id}</p>
								</div>
							</div>
						</div>
					{:else}
						<div class="p-4 bg-status-danger/10 border border-status-danger/30 rounded-lg">
							<div class="flex items-center gap-3">
								<svg class="w-5 h-5 text-status-danger flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
								<div>
									<span class="font-medium text-status-danger">Invalid Token</span>
									<p class="text-sm text-text-secondary mt-0.5">{validationResult.error}</p>
								</div>
							</div>
						</div>
					{/if}
				{/if}
			</div>

			<!-- Footer -->
			<div class="flex justify-end gap-3 p-5 border-t border-surface-600 bg-surface-900/50">
				<button onclick={closeModal} class="btn btn-ghost">
					Cancel
				</button>
				<button
					onclick={saveWhitelabel}
					disabled={!validationResult?.valid || saving}
					class="btn btn-primary"
				>
					{#if saving}
						<span class="spinner w-4 h-4 border-white/30 border-t-white"></span>
						Saving...
					{:else}
						Save Whitelabel
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if deleteConfirmId}
	{@const botToDelete = whitelabels.find(b => b.bot_id === deleteConfirmId)}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<!-- Backdrop -->
		<button
			onclick={closeDeleteModal}
			class="absolute inset-0 bg-black/60 backdrop-blur-sm"
			aria-label="Close modal"
		></button>

		<!-- Modal -->
		<div class="relative bg-surface-800 border border-surface-600 rounded-xl w-full max-w-md shadow-2xl">
			<!-- Header -->
			<div class="flex items-center justify-between p-5 border-b border-surface-600">
				<h2 class="text-lg font-semibold text-status-danger">Delete Whitelabel</h2>
				<button onclick={closeDeleteModal} class="p-1 rounded-lg hover:bg-surface-600 transition-colors" aria-label="Close modal">
					<svg class="w-5 h-5 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<!-- Content -->
			<div class="p-5 space-y-4">
				<div class="flex items-center gap-4 p-4 bg-surface-700 rounded-lg">
					{#if botToDelete?.bot_avatar}
						<img
							src={botToDelete.bot_avatar}
							alt={botToDelete.bot_name}
							class="w-12 h-12 rounded-full bg-surface-600"
						/>
					{:else}
						<div class="w-12 h-12 rounded-full bg-accent/20 flex items-center justify-center text-accent font-bold text-lg">
							{getBotInitial(botToDelete?.bot_name || '')}
						</div>
					{/if}
					<div>
						<p class="font-medium text-text-primary">{botToDelete?.bot_name || 'Unknown Bot'}</p>
						<p class="text-sm text-text-secondary font-mono">{botToDelete?.bot_id}</p>
					</div>
				</div>

				<div class="p-4 bg-status-danger/10 border border-status-danger/30 rounded-lg">
					<div class="flex gap-3">
						<svg class="w-5 h-5 text-status-danger flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
						</svg>
						<div class="text-sm">
							<p class="font-medium text-status-danger mb-1">This action cannot be undone</p>
							<p class="text-text-secondary">This will permanently delete the whitelabel bot. You will need to re-add it if you want to use it again.</p>
						</div>
					</div>
				</div>

				<div>
					<label for="deleteConfirm" class="label">Type <span class="font-mono text-status-danger">delete</span> to confirm</label>
					<input
						id="deleteConfirm"
						type="text"
						bind:value={deleteConfirmInput}
						placeholder="delete"
						class="input font-mono"
						autocomplete="off"
					/>
				</div>
			</div>

			<!-- Footer -->
			<div class="flex justify-end gap-3 p-5 border-t border-surface-600 bg-surface-900/50">
				<button onclick={closeDeleteModal} class="btn btn-ghost">
					Cancel
				</button>
				<button
					onclick={() => deleteWhitelabel(deleteConfirmId!)}
					disabled={deleteConfirmInput.toLowerCase() !== 'delete' || deleting}
					class="btn btn-danger"
				>
					{#if deleting}
						<span class="spinner w-4 h-4 border-white/30 border-t-white"></span>
						Deleting...
					{:else}
						Delete Whitelabel
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}
