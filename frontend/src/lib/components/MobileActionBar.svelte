<script lang="ts">
	let {
		onSave,
		onReset,
		saving = $bindable(false),
		disabled = false
	}: {
		onSave: () => void | Promise<void>;
		onReset: () => void;
		saving?: boolean;
		disabled?: boolean;
	} = $props();

	async function handleSave() {
		const result = onSave();
		if (result instanceof Promise) {
			await result;
		}
	}

	const isDisabled = $derived(saving || disabled);
</script>

<!-- Mobile: Sticky Bottom Bar -->
<div class="lg:hidden fixed bottom-0 left-0 right-0 z-40 bg-surface-800 border-t border-surface-600 shadow-2xl safe-bottom">
	<div class="p-4 pb-safe">
		<div class="flex gap-2">
			<button
				onclick={onReset}
				disabled={isDisabled}
				class="btn btn-secondary flex-1"
			>
				Reset
			</button>
			<button
				onclick={handleSave}
				disabled={isDisabled}
				class="btn btn-primary flex-1 relative"
			>
				{#if saving}
					<span class="flex items-center justify-center gap-2">
						<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						<span>Saving...</span>
					</span>
				{:else}
					Save Changes
				{/if}
			</button>
		</div>
	</div>
</div>

<!-- Desktop: Inline Bar -->
<div class="hidden lg:flex items-center justify-between">
	<div class="flex-1">
		<slot name="header" />
	</div>
	<div class="flex gap-2">
		<button
			onclick={onReset}
			disabled={isDisabled}
			class="btn btn-secondary"
		>
			Reset
		</button>
		<button
			onclick={handleSave}
			disabled={isDisabled}
			class="btn btn-primary relative"
		>
			{#if saving}
				<span class="flex items-center justify-center gap-2">
					<svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					<span>Saving...</span>
				</span>
			{:else}
				Save Changes
			{/if}
		</button>
	</div>
</div>

<style>
	/* Safe area padding for modern phones (notch/home indicator) */
	.safe-bottom {
		padding-bottom: env(safe-area-inset-bottom);
	}

	.pb-safe {
		padding-bottom: max(1rem, env(safe-area-inset-bottom));
	}

	/* Enhanced shadow for elevation */
	.shadow-2xl {
		box-shadow:
			0 -4px 6px -1px rgb(0 0 0 / 0.1),
			0 -2px 4px -2px rgb(0 0 0 / 0.1),
			0 -20px 25px -5px rgb(0 0 0 / 0.2);
	}
</style>
