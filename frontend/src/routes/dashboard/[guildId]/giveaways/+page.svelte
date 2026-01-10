<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { settingsAPI, type GiveawaySettings } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';
	import ChannelSelect from '$lib/components/ChannelSelect.svelte';
	import RoleSelect from '$lib/components/RoleSelect.svelte';
	import MobileActionBar from '$lib/components/MobileActionBar.svelte';

	$: guildId = $page.params.guildId;

	let loading = true;
	let saving = false;
	let settings: GiveawaySettings = { guild_id: guildId };

	onMount(async () => {
		try {
			settings = await settingsAPI.getGiveaway(guildId);
		} catch {
			// Empty settings is fine
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		try {
			await settingsAPI.updateGiveaway(guildId, settings);
			toast.success('Giveaway settings saved');
		} catch (err) {
			toast.error('Failed to save settings');
		} finally {
			saving = false;
		}
	}

	async function reset() {
		if (!confirm('Are you sure you want to reset all giveaway settings?')) return;
		try {
			await settingsAPI.deleteGiveaway(guildId);
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
			<h1 slot="header" class="text-fluid-xl font-semibold text-text-primary">Giveaway Settings</h1>
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
				<h2 class="text-fluid-sm font-medium text-text-primary mb-2">Logging</h2>
				<p class="text-fluid-sm text-text-secondary mb-4">
					Configure where giveaway events are logged.
				</p>
				<div>
					<ChannelSelect bind:value={settings.log_channel} label="Log Channel" type="text" />
					<p class="text-fluid-xs text-text-secondary mt-1">Channel where giveaway creation, completion, and winner announcements are logged.</p>
				</div>
			</div>

			<!-- Roles Section -->
			<div class="card">
				<h2 class="text-fluid-sm font-medium text-text-primary mb-4">Roles</h2>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<RoleSelect bind:value={settings.alert_role} label="Alert Role" />
					<RoleSelect bind:value={settings.host_role} label="Host Role (Deprecated)" />
				</div>
			</div>
		</div>
	{/if}
</div>
