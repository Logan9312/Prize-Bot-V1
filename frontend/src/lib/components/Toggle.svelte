<script lang="ts">
	let { checked = $bindable(false), label = '' }: {
		checked?: boolean;
		label?: string;
	} = $props();
</script>

<label class="flex items-center gap-3 cursor-pointer py-2 -mx-2 px-2 rounded-lg hover:bg-surface-700 transition-colors min-h-[44px]">
	<div class="relative">
		<input type="checkbox" bind:checked class="sr-only peer" />
		<div
			class="w-11 h-6 sm:w-10 sm:h-5 bg-surface-600 rounded-full peer-checked:bg-accent transition-colors"
		></div>
		<div
			class="absolute left-0.5 top-0.5 w-5 h-5 sm:w-4 sm:h-4 bg-white rounded-full transition-transform peer-checked:translate-x-5"
		></div>
	</div>
	{#if label}
		<span class="text-fluid-sm text-text-secondary">{label}</span>
	{/if}
</label>
