<script lang="ts">
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { Role } from '$lib/api/client';

	export let value: string = '';
	export let label: string = '';

	const roles = getContext<Writable<Role[]>>('roles');
	let isOpen = false;
	let dropdownElement: HTMLDivElement;
	const dropdownId = `role-select-${Math.random().toString(36).substring(2, 11)}`;

	function intToHex(color: number): string {
		if (!color) return '#99aab5';
		return '#' + color.toString(16).padStart(6, '0');
	}

	function toggleDropdown() {
		isOpen = !isOpen;
	}

	function selectRole(roleId: string) {
		value = roleId;
		isOpen = false;
	}

	function handleClickOutside(event: MouseEvent) {
		if (dropdownElement && !dropdownElement.contains(event.target as Node)) {
			isOpen = false;
		}
	}

	$: selectedRole = $roles.find(r => r.id === value);
</script>

<svelte:window on:click={handleClickOutside} />

<div bind:this={dropdownElement}>
	{#if label}
		<label for={dropdownId} class="label block text-sm font-medium mb-1">{label}</label>
	{/if}
	<div class="relative">
		<button
			id={dropdownId}
			type="button"
			onclick={toggleDropdown}
			class="select w-full text-left"
		>
			{#if selectedRole}
				<span style="color: {intToHex(selectedRole.color)}">
					@ {selectedRole.name}
				</span>
			{:else}
				<span class="text-text-muted">None</span>
			{/if}
		</button>

		{#if isOpen}
			<div class="absolute z-10 w-full mt-1 bg-surface-700 border border-surface-600 rounded-lg shadow-lg max-h-60 overflow-auto">
				<button
					type="button"
					onclick={() => selectRole('')}
					class="w-full text-left px-4 py-3 sm:px-3 sm:py-2.5 hover:bg-surface-600 text-text-muted transition-colors min-h-[44px]"
				>
					None
				</button>
				{#each $roles as role}
					<button
						type="button"
						onclick={() => selectRole(role.id)}
						class="w-full text-left px-4 py-3 sm:px-3 sm:py-2.5 hover:bg-surface-600 transition-colors min-h-[44px]"
						style="color: {intToHex(role.color)}"
					>
						@ {role.name}
					</button>
				{/each}
			</div>
		{/if}
	</div>
</div>
