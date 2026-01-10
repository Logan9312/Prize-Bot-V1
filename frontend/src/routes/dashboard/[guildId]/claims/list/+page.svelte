<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { guildsAPI, claimsAPI, type ClaimListItem } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import FeatureTabs from '$lib/components/FeatureTabs.svelte';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let claims: ClaimListItem[] = $state([]);

	// Edit modal state
	let editingClaim: ClaimListItem | null = $state(null);
	let editForm = $state({ item: '', winner: '', cost: 0 });
	let saving = $state(false);

	// Confirm modal state
	let confirmAction: { type: 'resend' | 'cancel'; claim: ClaimListItem } | null = $state(null);
	let actionLoading = $state(false);

	onMount(async () => {
		await loadClaims();
	});

	async function loadClaims() {
		loading = true;
		try {
			const res = await guildsAPI.listClaims(guildId);
			claims = res.claims || [];
		} catch (e) {
			console.error('Failed to load claims:', e);
			toast.error('Failed to load claims');
		} finally {
			loading = false;
		}
	}

	function openEditModal(claim: ClaimListItem) {
		editingClaim = claim;
		editForm = {
			item: claim.item,
			winner: claim.winner,
			cost: claim.cost
		};
	}

	function closeEditModal() {
		editingClaim = null;
		editForm = { item: '', winner: '', cost: 0 };
	}

	async function saveEdit() {
		if (!editingClaim) return;
		saving = true;
		try {
			await claimsAPI.update(guildId, editingClaim.message_id, editForm);
			toast.success('Claim updated successfully');
			closeEditModal();
			await loadClaims();
		} catch (e: any) {
			toast.error(e.message || 'Failed to update claim');
		} finally {
			saving = false;
		}
	}

	function openConfirmModal(type: 'resend' | 'cancel', claim: ClaimListItem) {
		confirmAction = { type, claim };
	}

	function closeConfirmModal() {
		confirmAction = null;
	}

	async function executeAction() {
		if (!confirmAction) return;
		actionLoading = true;
		try {
			if (confirmAction.type === 'resend') {
				await claimsAPI.resend(guildId, confirmAction.claim.message_id);
				toast.success('Claim resent successfully');
			} else {
				await claimsAPI.cancel(guildId, confirmAction.claim.message_id);
				toast.success('Claim cancelled successfully');
			}
			closeConfirmModal();
			await loadClaims();
		} catch (e: any) {
			toast.error(e.message || `Failed to ${confirmAction.type} claim`);
		} finally {
			actionLoading = false;
		}
	}

	function handleModalKeydown(e: KeyboardEvent, closeFunc: () => void) {
		if (e.key === 'Escape') closeFunc();
	}
</script>

<div class="space-y-6">
	<!-- Header with tabs -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<FeatureTabs feature="claims" activeTab="list" {guildId} listCount={claims.length} />

		<button
			onclick={loadClaims}
			disabled={loading}
			class="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-800 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
		>
			<svg class="w-4 h-4 {loading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
			</svg>
			Refresh
		</button>
	</div>

	<!-- Summary -->
	<p class="text-sm text-text-secondary">
		{claims.length} open claim{claims.length !== 1 ? 's' : ''}
	</p>

	<!-- Content -->
	{#if loading}
		<div class="flex justify-center py-12">
			<div class="w-8 h-8 border-2 border-accent border-t-transparent rounded-full animate-spin"></div>
		</div>
	{:else if claims.length === 0}
		<div class="bg-surface-800 border border-surface-600 rounded-lg text-center py-16 px-4">
			<div class="w-16 h-16 mx-auto mb-4 rounded-full bg-surface-700 flex items-center justify-center">
				<svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
				</svg>
			</div>
			<h3 class="text-lg font-medium text-text-primary mb-2">No open claims</h3>
			<p class="text-text-secondary max-w-sm mx-auto">
				Claims are created automatically when auctions end or giveaways are drawn. Check back later!
			</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each claims as claim}
				<div class="bg-surface-800 border border-surface-600 rounded-lg p-4 hover:border-surface-500 transition-colors">
					<div class="flex justify-between items-start gap-4">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2">
								<span class="w-2 h-2 rounded-full bg-yellow-400 flex-shrink-0" title="Pending"></span>
								<h3 class="font-medium text-text-primary truncate">{claim.item}</h3>
							</div>
							<div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm">
								<span class="text-text-secondary">
									Winner: <span class="text-text-primary font-mono text-xs">{claim.winner}</span>
								</span>
								{#if claim.cost > 0}
									<span class="text-text-muted">
										Cost: <span class="text-text-primary">${claim.cost}</span>
									</span>
								{/if}
							</div>
						</div>
						<div class="flex items-center gap-1">
							<button
								onclick={() => openEditModal(claim)}
								class="p-2 text-text-muted hover:text-text-primary hover:bg-surface-700 rounded-lg transition-colors cursor-pointer"
								title="Edit claim"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
								</svg>
							</button>
							<button
								onclick={() => openConfirmModal('resend', claim)}
								class="p-2 text-text-muted hover:text-accent hover:bg-surface-700 rounded-lg transition-colors cursor-pointer"
								title="Resend claim message"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
								</svg>
							</button>
							<button
								onclick={() => openConfirmModal('cancel', claim)}
								class="p-2 text-text-muted hover:text-red-400 hover:bg-surface-700 rounded-lg transition-colors cursor-pointer"
								title="Cancel claim"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Edit Modal -->
{#if editingClaim}
	<div
		class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
		role="dialog"
		aria-modal="true"
		aria-labelledby="edit-modal-title"
		onclick={closeEditModal}
		onkeydown={(e) => handleModalKeydown(e, closeEditModal)}
		tabindex="-1"
	>
		<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
		<div class="bg-surface-800 border border-surface-600 rounded-lg p-6 w-full max-w-md mx-4" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="document">
			<h2 id="edit-modal-title" class="text-lg font-semibold text-text-primary mb-4">Edit Claim</h2>

			<div class="space-y-4">
				<div>
					<label class="block text-sm text-text-secondary mb-1" for="edit-item">Item</label>
					<input
						id="edit-item"
						type="text"
						bind:value={editForm.item}
						class="w-full bg-surface-700 border border-surface-600 rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-accent"
					/>
				</div>

				<div>
					<label class="block text-sm text-text-secondary mb-1" for="edit-winner">Winner (User ID)</label>
					<input
						id="edit-winner"
						type="text"
						bind:value={editForm.winner}
						class="w-full bg-surface-700 border border-surface-600 rounded-lg px-3 py-2 text-text-primary font-mono text-sm focus:outline-none focus:border-accent"
					/>
				</div>

				<div>
					<label class="block text-sm text-text-secondary mb-1" for="edit-cost">Cost</label>
					<input
						id="edit-cost"
						type="number"
						bind:value={editForm.cost}
						min="0"
						step="0.01"
						class="w-full bg-surface-700 border border-surface-600 rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-accent"
					/>
				</div>
			</div>

			<div class="flex justify-end gap-3 mt-6">
				<button
					onclick={closeEditModal}
					class="px-4 py-2 text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
				>
					Cancel
				</button>
				<button
					onclick={saveEdit}
					disabled={saving}
					class="px-4 py-2 bg-accent hover:bg-accent/80 text-white rounded-lg transition-colors disabled:opacity-50 cursor-pointer"
				>
					{saving ? 'Saving...' : 'Save'}
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Confirm Modal -->
{#if confirmAction}
	<div
		class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
		role="dialog"
		aria-modal="true"
		aria-labelledby="confirm-modal-title"
		onclick={closeConfirmModal}
		onkeydown={(e) => handleModalKeydown(e, closeConfirmModal)}
		tabindex="-1"
	>
		<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
		<div class="bg-surface-800 border border-surface-600 rounded-lg p-6 w-full max-w-md mx-4" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="document">
			<h2 id="confirm-modal-title" class="text-lg font-semibold text-text-primary mb-2">
				{confirmAction.type === 'resend' ? 'Resend Claim' : 'Cancel Claim'}
			</h2>
			<p class="text-text-secondary mb-4">
				{#if confirmAction.type === 'resend'}
					Are you sure you want to resend the claim for <strong class="text-text-primary">{confirmAction.claim.item}</strong>? This will post a new claim message in the log channel.
				{:else}
					Are you sure you want to cancel the claim for <strong class="text-text-primary">{confirmAction.claim.item}</strong>? This action cannot be undone.
				{/if}
			</p>

			<div class="flex justify-end gap-3">
				<button
					onclick={closeConfirmModal}
					class="px-4 py-2 text-text-secondary hover:text-text-primary transition-colors cursor-pointer"
				>
					Cancel
				</button>
				<button
					onclick={executeAction}
					disabled={actionLoading}
					class="px-4 py-2 rounded-lg transition-colors disabled:opacity-50 cursor-pointer {confirmAction.type === 'cancel' ? 'bg-red-600 hover:bg-red-500 text-white' : 'bg-accent hover:bg-accent/80 text-white'}"
				>
					{#if actionLoading}
						Processing...
					{:else}
						{confirmAction.type === 'resend' ? 'Resend' : 'Cancel Claim'}
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}
