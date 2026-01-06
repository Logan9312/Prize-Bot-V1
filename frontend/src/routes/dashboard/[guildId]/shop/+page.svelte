<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type ShopSettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';

	$: guildId = $page.params.guildId;

	let loading = true;
	let saving = false;
	let settings: ShopSettings = { guild_id: guildId };

	onMount(async () => {
		try {
			settings = await settingsAPI.getShop(guildId);
		} catch {
			// Empty settings is fine
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		try {
			await settingsAPI.updateShop(guildId, settings);
			toast.success('Shop settings saved');
		} catch (err) {
			toast.error('Failed to save settings');
		} finally {
			saving = false;
		}
	}

	async function reset() {
		if (!confirm('Are you sure you want to reset all shop settings?')) return;
		try {
			await settingsAPI.deleteShop(guildId);
			settings = { guild_id: guildId };
			toast.success('Settings reset');
		} catch {
			toast.error('Failed to reset settings');
		}
	}
</script>

<div>
	<div class="flex items-center justify-between mb-6">
		<h1 class="text-xl font-semibold text-text-primary">Shop Settings</h1>
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
				<h2 class="text-sm font-medium text-text-primary mb-4">Logging</h2>
				<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
			</div>

			<!-- Roles Section -->
			<div class="card">
				<h2 class="text-sm font-medium text-text-primary mb-4">Roles</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<RoleSelect bind:value={settings.host_role} label="Host Role" />
					<RoleSelect bind:value={settings.alert_role} label="Alert Role" />
				</div>
			</div>
		</div>
	{/if}
</div>
