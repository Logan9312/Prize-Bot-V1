<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { Role } from '$lib/api/client';

	export let value: string = '';
	export let label: string = '';

	const roles = getContext<Writable<Role[]>>('roles');

	function intToHex(color: number): string {
		if (!color) return '#99aab5';
		return '#' + color.toString(16).padStart(6, '0');
	}
</script>

<div>
	{#if label}
		<label class="label">{label}</label>
	{/if}
	<select bind:value class="select">
		<option value="">None</option>
		{#each $roles as role}
			<option value={role.id} style="color: {intToHex(role.color)}">
				@ {role.name}
			</option>
		{/each}
	</select>
</div>
