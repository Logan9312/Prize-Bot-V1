<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type ClaimSettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import Toggle from '$lib/components/Toggle.svelte';

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
	<div class="flex items-center justify-between mb-6">
		<h1 class="text-xl font-semibold text-text-primary">Claim Settings</h1>
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
					<ChannelSelect bind:value={settings.category} label="Ticket Category" type="category" />
					<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
				</div>
				<div class="mt-4">
					<label class="label">Channel Prefix</label>
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
				<h2 class="text-sm font-medium text-text-primary mb-4">Staff</h2>
				<RoleSelect bind:value={settings.staff_role} label="Staff Role" />
			</div>

			<!-- Instructions Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Instructions</h2>
				<label class="label">Ticket Instructions</label>
				<textarea
					bind:value={settings.instructions}
					placeholder="Instructions shown to users when they open a ticket..."
					rows="4"
					class="input"
				></textarea>
			</div>

			<!-- Options Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Options</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
					<div>
						<label class="label">Ticket Expiration</label>
						<input
							type="text"
							bind:value={settings.expiration}
							placeholder="e.g., 7d, 24h"
							class="input"
						/>
					</div>
				</div>
				<Toggle bind:checked={settings.disable_claiming} label="Disable Claiming System" />
			</div>
		</div>
	{/if}
</div>
