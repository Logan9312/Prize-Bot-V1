<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type AuctionSettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import DurationInput from '$lib/components/DurationInput.svelte';
	import Toggle from '$lib/components/Toggle.svelte';
	import MobileActionBar from '$lib/components/MobileActionBar.svelte';

	const guildId = $derived($page.params.guildId!);

	let loading = $state(true);
	let saving = $state(false);
	let settings: AuctionSettings = $state({ guild_id: '' });

	onMount(async () => {
		settings = { guild_id: guildId };
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

{#snippet header()}
	<h1 class="text-fluid-xl font-semibold text-text-primary">Auction Settings</h1>
{/snippet}

<div>
	<div class="mb-4 lg:mb-6">
		<MobileActionBar onSave={save} onReset={reset} bind:saving {header} />
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="spinner spinner-lg"></div>
		</div>
	{:else}
		<div class="space-y-4 lg:space-y-6 pb-20 lg:pb-0">
			<!-- Channels Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Channels</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Configure where auction channels are created and where logs are sent.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<ChannelSelect bind:value={settings.category} label="Auction Category" type="category" />
					<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
				</div>
				<div class="mt-4">
					<label for="auction-channel-prefix" class="label">Channel Prefix</label>
					<p class="text-fluid-xs text-text-secondary mb-2">
						Text prepended to auction channel names (e.g., "auction-item-name").
					</p>
					<input
						id="auction-channel-prefix"
						type="text"
						bind:value={settings.channel_prefix}
						placeholder="auction-"
						class="input max-w-xs"
					/>
				</div>
			</div>

			<!-- Roles Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Roles</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Configure which roles receive notifications about auctions.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<RoleSelect bind:value={settings.alert_role} label="Alert Role" />
						<p class="text-fluid-xs text-text-secondary mt-1">Pinged when new auctions are created.</p>
					</div>
					<div>
						<RoleSelect bind:value={settings.host_role} label="Host Role (Deprecated)" />
						<p class="text-fluid-xs text-text-secondary mt-1">No longer used. Will be removed in a future update.</p>
					</div>
				</div>
			</div>

			<!-- Currency Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Currency</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Customize how bid amounts are displayed in auctions. These settings override the server-wide currency settings for auctions.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label for="auction-currency-symbol" class="label">Currency Symbol</label>
						<p class="text-fluid-xs text-text-secondary mb-2">Symbol shown next to bid amounts (e.g., $, coins, tokens).</p>
						<input id="auction-currency-symbol" type="text" bind:value={settings.currency} placeholder="$" class="input" />
					</div>
					<div>
						<label for="auction-currency-side" class="label">Currency Side</label>
						<p class="text-fluid-xs text-text-secondary mb-2">Where to display the symbol relative to the amount.</p>
						<select id="auction-currency-side" bind:value={settings.currency_side} class="select">
							<option value="">Default (Left)</option>
							<option value="left">Left ($100)</option>
							<option value="right">Right (100$)</option>
						</select>
					</div>
				</div>
				<div class="mt-4 space-y-3">
					<div>
						<Toggle bind:checked={settings.integer_only} label="Integer Only (no decimals)" />
						<p class="text-fluid-xs text-text-secondary mt-1 ml-11">Only allow whole number bids, no cents or fractions.</p>
					</div>
					<div>
						<Toggle bind:checked={settings.use_currency} label="Use Server Currency for Bids" />
						<p class="text-fluid-xs text-text-secondary mt-1 ml-11">Deduct bids from users' server currency balance instead of just tracking amounts.</p>
					</div>
				</div>
			</div>

			<!-- Anti-Snipe Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Anti-Snipe Settings</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
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
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Options</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Additional settings to control auction behavior.
				</p>
				<div class="space-y-3">
					<div>
						<Toggle
							bind:checked={settings.channel_lock}
							label="Lock Auction to Command Channel"
						/>
						<p class="text-fluid-xs text-text-secondary mt-1 ml-11">Restrict bidding commands to only work in the auction's dedicated channel.</p>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
