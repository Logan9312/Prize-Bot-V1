<script lang="ts">
	import { onMount } from 'svelte';
	import { premiumAPI, type UserPremiumStatus } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	let loading = $state(true);
	let portalLoading = $state(false);
	let status: UserPremiumStatus | null = $state(null);

	onMount(async () => {
		try {
			status = await premiumAPI.getUserStatus();
		} catch (err) {
			toast.error('Failed to load subscription status');
		} finally {
			loading = false;
		}
	});

	async function openBillingPortal() {
		portalLoading = true;
		try {
			const response = await premiumAPI.createPortalSession();
			window.location.href = response.url;
		} catch (err) {
			toast.error('Failed to open billing portal. Make sure you have an active subscription.');
		} finally {
			portalLoading = false;
		}
	}

	function formatDate(timestamp: number): string {
		return new Date(timestamp * 1000).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'active':
				return 'text-status-success';
			case 'past_due':
				return 'text-status-warning';
			case 'canceled':
				return 'text-status-danger';
			default:
				return 'text-text-secondary';
		}
	}
</script>

<div class="max-w-2xl mx-auto">
	<!-- Page Header -->
	<div class="mb-6">
		<a
			href="/dashboard"
			class="flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary transition-colors mb-4"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
			Back to servers
		</a>
		<h1 class="text-xl font-semibold text-text-primary">Subscription</h1>
		<p class="text-sm text-text-secondary mt-1">Manage your Prize Bot Premium subscription</p>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else if status}
		<div class="space-y-6">
			<!-- Premium Status Card -->
			<div class="card">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium text-text-primary">Premium Status</h2>
					{#if status.is_premium}
						<span class="px-3 py-1 rounded-full text-xs font-medium bg-accent/20 text-accent">
							Active
						</span>
					{:else}
						<span class="px-3 py-1 rounded-full text-xs font-medium bg-surface-600 text-text-muted">
							Inactive
						</span>
					{/if}
				</div>

				{#if status.is_premium && status.subscriptions && status.subscriptions.length > 0}
					<p class="text-text-secondary">
						You have an active Premium subscription. Thank you for supporting Prize Bot!
					</p>
				{:else if status.is_premium}
					<p class="text-text-secondary">
						You have Premium access. Thank you for supporting Prize Bot!
					</p>
				{:else}
					<p class="text-text-secondary mb-4">
						Upgrade to Premium to unlock additional features like auction scheduling, bulk
						auctions, and more.
					</p>
					<a
						href="https://prizebot.dev/premium"
						target="_blank"
						rel="noopener"
						class="btn btn-primary"
					>
						Get Premium
					</a>
				{/if}
			</div>

			<!-- Subscriptions List -->
			{#if status.subscriptions && status.subscriptions.length > 0}
				<div class="card">
					<h2 class="text-lg font-medium text-text-primary mb-4">Your Subscriptions</h2>
					<div class="space-y-4">
						{#each status.subscriptions as subscription}
							<div class="p-4 bg-surface-800 border border-surface-600 rounded-lg">
								<div class="flex items-center justify-between mb-2">
									<span class="font-medium text-text-primary">{subscription.plan_name}</span>
									<span class="text-sm capitalize {getStatusColor(subscription.status)}">
										{subscription.status.replace('_', ' ')}
									</span>
								</div>
								<div class="text-sm text-text-secondary space-y-1">
									{#if subscription.guild_id}
										<p>
											Linked to server: <span class="text-text-primary"
												>{subscription.guild_id}</span
											>
										</p>
									{:else}
										<p class="text-status-warning">
											Not linked to any server. Use /premium activate in Discord.
										</p>
									{/if}
									<p>Renews: {formatDate(subscription.current_period_end)}</p>
								</div>
							</div>
						{/each}
					</div>
				</div>

				<!-- Billing Portal -->
				<div class="card">
					<h2 class="text-lg font-medium text-text-primary mb-2">Manage Billing</h2>
					<p class="text-text-secondary mb-4">
						Update payment methods, view invoices, or cancel your subscription.
					</p>
					<button onclick={openBillingPortal} disabled={portalLoading} class="btn btn-secondary">
						{#if portalLoading}
							<span class="spinner spinner-sm mr-2"></span>
						{/if}
						Open Billing Portal
					</button>
				</div>
			{:else if status.is_premium}
				<!-- Premium without Stripe subscription (granted/lifetime) -->
				<div class="card">
					<h2 class="text-lg font-medium text-text-primary mb-2">Premium Access</h2>
					<p class="text-text-secondary">
						Your premium access has been granted directly. No billing management is required.
					</p>
				</div>
			{/if}
		</div>
	{/if}
</div>
