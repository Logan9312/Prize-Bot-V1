<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type CurrencySettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import MobileActionBar from '$lib/components/MobileActionBar.svelte';

	$: guildId = $page.params.guildId!;

	let loading = true;
	let saving = false;
	let settings: CurrencySettings = { guild_id: guildId };

	onMount(async () => {
		try {
			settings = await settingsAPI.getCurrency(guildId);
		} catch {
			// Empty settings is fine
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		try {
			await settingsAPI.updateCurrency(guildId, settings);
			toast.success('Currency settings saved');
		} catch (err) {
			toast.error('Failed to save settings');
		} finally {
			saving = false;
		}
	}

	async function reset() {
		if (!confirm('Are you sure you want to reset all currency settings?')) return;
		try {
			await settingsAPI.deleteCurrency(guildId);
			settings = { guild_id: guildId };
			toast.success('Settings reset');
		} catch {
			toast.error('Failed to reset settings');
		}
	}
</script>

<div>
	<div class="mb-4 lg:mb-6">
		<MobileActionBar onSave={save} onReset={reset} bind:saving>
			<h1 slot="header" class="text-fluid-xl font-semibold text-text-primary">Currency Settings</h1>
		</MobileActionBar>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else}
		<div class="space-y-4 lg:space-y-6 pb-20 lg:pb-0">
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Server Currency</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Configure the default currency used across all features in your server. Individual features like auctions can override these settings.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label class="label">Currency Symbol</label>
						<p class="text-fluid-xs text-text-secondary mb-2">The symbol or name displayed with amounts (e.g., $, coins, credits, points).</p>
						<input
							type="text"
							bind:value={settings.currency}
							placeholder="$ or coins"
							class="input"
						/>
					</div>
					<div>
						<label class="label">Symbol Position</label>
						<p class="text-fluid-xs text-text-secondary mb-2">Choose whether the symbol appears before or after the amount.</p>
						<select bind:value={settings.side} class="select">
							<option value="">Default (Left)</option>
							<option value="left">Left ($100)</option>
							<option value="right">Right (100$)</option>
						</select>
					</div>
				</div>

				<!-- Preview -->
				<div class="mt-6 p-4 bg-surface-800 border border-surface-600 rounded-lg">
					<label class="label">Preview</label>
					<p class="text-fluid-xs text-text-secondary mb-2">This is how currency amounts will appear to users.</p>
					<p class="text-lg text-text-primary">
						{#if settings.side === 'right'}
							100{settings.currency || '$'}
						{:else}
							{settings.currency || '$'}100
						{/if}
					</p>
				</div>
			</div>
		</div>
	{/if}
</div>
