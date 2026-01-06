<script lang="ts">
	export let checked: boolean = false;
	export let label: string = '';
</script>

<label class="flex items-center gap-3 cursor-pointer">
	<div class="relative">
		<input type="checkbox" bind:checked class="sr-only peer" />
		<div
			class="w-10 h-5 bg-surface-600 rounded-full peer-checked:bg-accent transition-colors"
		></div>
		<div
			class="absolute left-0.5 top-0.5 w-4 h-4 bg-white rounded-full transition-transform peer-checked:translate-x-5"
		></div>
	</div>
	{#if label}
		<span class="text-sm text-text-secondary">{label}</span>
	{/if}
</label>
