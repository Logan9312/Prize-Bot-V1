<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type ClaimSettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import Toggle from '$lib/components/Toggle.svelte';
	import MobileActionBar from '$lib/components/MobileActionBar.svelte';

	$: guildId = $page.params.guildId;

	let loading = true;
	let saving = false;
	let settings: ClaimSettings = { guild_id: guildId };

	onMount(async () => {
		try {
			settings = await settingsAPI.getClaim(guildId);
		} catch {
			// Empty settings is fine
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		try {
			await settingsAPI.updateClaim(guildId, settings);
			toast.success('Claim settings saved');
		} catch (err) {
			toast.error('Failed to save settings');
		} finally {
			saving = false;
		}
	}

	async function reset() {
		if (!confirm('Are you sure you want to reset all claim settings?')) return;
		try {
			await settingsAPI.deleteClaim(guildId);
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
			<h1 slot="header" class="text-fluid-xl font-semibold text-text-primary">Claim Settings</h1>
		</MobileActionBar>
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
					Configure where claim tickets are created and where activity is logged.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<ChannelSelect bind:value={settings.category} label="Ticket Category" type="category" />
						<p class="text-fluid-xs text-text-secondary mt-1">Category where new ticket channels will be created.</p>
					</div>
					<div>
						<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
						<p class="text-fluid-xs text-text-secondary mt-1">Channel where ticket events are logged.</p>
					</div>
				</div>
				<div class="mt-4">
					<label class="label">Channel Prefix</label>
					<p class="text-fluid-xs text-text-secondary mb-2">
						Text prepended to ticket channel names (e.g., "ticket-username").
					</p>
					<input
						type="text"
						bind:value={settings.channel_prefix}
						placeholder="ticket-"
						class="input max-w-xs"
					/>
				</div>
			</div>

			<!-- Role Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Staff</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Configure which role can manage claim tickets.
				</p>
				<div>
					<RoleSelect bind:value={settings.staff_role} label="Staff Role" />
					<p class="text-fluid-xs text-text-secondary mt-1">Members with this role can view and manage all claim tickets.</p>
				</div>
			</div>

			<!-- Instructions Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Instructions</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Customize the message users see when they open a claim ticket.
				</p>
				<label class="label">Ticket Instructions</label>
				<p class="text-fluid-xs text-text-secondary mb-2">
					This message is sent automatically when a ticket is created. Include any info users need to provide.
				</p>
				<textarea
					bind:value={settings.instructions}
					placeholder="Instructions shown to users when they open a ticket..."
					rows="4"
					class="input"
				></textarea>
			</div>

			<!-- Options Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Options</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Additional settings to control the claiming system behavior.
				</p>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
					<div>
						<label class="label">Ticket Expiration</label>
						<p class="text-fluid-xs text-text-secondary mb-2">
							Automatically close inactive tickets after this duration. Use formats like "7d" (7 days) or "24h" (24 hours). Leave empty to disable.
						</p>
						<input
							type="text"
							bind:value={settings.expiration}
							placeholder="e.g., 7d, 24h"
							class="input"
						/>
					</div>
				</div>
				<div>
					<Toggle bind:checked={settings.disable_claiming} label="Disable Claiming System" />
					<p class="text-fluid-xs text-text-secondary mt-1 ml-11">Temporarily disable the ability for users to create new claim tickets.</p>
				</div>
			</div>
		</div>
	{/if}
</div>
