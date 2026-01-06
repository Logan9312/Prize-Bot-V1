<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type AuctionSettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import DurationInput from '$lib/components/DurationInput.svelte';
	import Toggle from '$lib/components/Toggle.svelte';

	$: guildId = $page.params.guildId;

	let loading = true;
	let saving = false;
	let settings: AuctionSettings = { guild_id: guildId };

	onMount(async () => {
		try {
			settings = await settingsAPI.getAuction(guildId);
		} catch {
			// Empty settings is fine
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		try {
			await settingsAPI.updateAuction(guildId, settings);
			toast.success('Auction settings saved');
		} catch (err) {
			toast.error('Failed to save settings');
		} finally {
			saving = false;
		}
	}

	async function reset() {
		if (!confirm('Are you sure you want to reset all auction settings?')) return;
		try {
			await settingsAPI.deleteAuction(guildId);
			settings = { guild_id: guildId };
			toast.success('Settings reset');
		} catch {
			toast.error('Failed to reset settings');
		}
	}
</script>

<div>
	<div class="flex items-center justify-between mb-6">
		<h1 class="text-xl font-semibold text-text-primary">Auction Settings</h1>
		<div class="flex gap-2">
			<button on:click={reset} class="btn btn-secondary">Reset</button>
			<button on:click={save} disabled={saving} class="btn btn-primary">
				{saving ? 'Saving...' : 'Save Changes'}
			</button>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else}
		<div class="space-y-6">
			<!-- Channels Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Channels</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<ChannelSelect bind:value={settings.category} label="Auction Category" type="category" />
					<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
				</div>
				<div class="mt-4">
					<label class="label">Channel Prefix</label>
					<input
						type="text"
						bind:value={settings.channel_prefix}
						placeholder="auction-"
						class="input max-w-xs"
					/>
				</div>
			</div>

			<!-- Roles Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Roles</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<RoleSelect bind:value={settings.alert_role} label="Alert Role" />
					<RoleSelect bind:value={settings.host_role} label="Host Role (Deprecated)" />
				</div>
			</div>

			<!-- Currency Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Currency</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label class="label">Currency Symbol</label>
						<input type="text" bind:value={settings.currency} placeholder="$" class="input" />
					</div>
					<div>
						<label class="label">Currency Side</label>
						<select bind:value={settings.currency_side} class="select">
							<option value="">Default (Left)</option>
							<option value="left">Left ($100)</option>
							<option value="right">Right (100$)</option>
						</select>
					</div>
				</div>
				<div class="mt-4 space-y-3">
					<Toggle bind:checked={settings.integer_only} label="Integer Only (no decimals)" />
					<Toggle bind:checked={settings.use_currency} label="Use Server Currency for Bids" />
				</div>
			</div>

			<!-- Anti-Snipe Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-2">Anti-Snipe Settings</h2>
				<p class="text-sm text-text-secondary mb-4">
					Prevent last-second bids by extending auction time when bids are placed near the end.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<DurationInput
						bind:value={settings.snipe_extension}
						label="Snipe Extension (time added per snipe)"
					/>
					<DurationInput
						bind:value={settings.snipe_range}
						label="Snipe Range (trigger window before end)"
					/>
					<DurationInput
						bind:value={settings.snipe_limit}
						label="Snipe Limit (max total extensions)"
					/>
					<DurationInput
						bind:value={settings.snipe_cap}
						label="Snipe Cap (max auction duration)"
					/>
				</div>
			</div>

			<!-- Options Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Options</h2>
				<div class="space-y-3">
					<Toggle
						bind:checked={settings.channel_lock}
						label="Lock Auction to Command Channel"
					/>
				</div>
			</div>
		</div>
	{/if}
</div>
