<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { toast, type Toast } from '$lib/stores/toast';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let toasts: Toast[] = $state([]);
	toast.subscribe((value) => (toasts = value));

	onMount(() => {
		auth.init();
	});

	function getToastIcon(type: string) {
		switch (type) {
			case 'success':
				return 'M5 13l4 4L19 7';
			case 'error':
				return 'M6 18L18 6M6 6l12 12';
			default:
				return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z';
		}
	}

	function getToastClasses(type: string) {
		switch (type) {
			case 'success':
				return 'bg-status-success/10 border-status-success/20 text-status-success';
			case 'error':
				return 'bg-status-danger/10 border-status-danger/20 text-status-danger';
			default:
				return 'bg-accent/10 border-accent/20 text-accent';
		}
	}
</script>

<div class="min-h-screen bg-surface-900">
	{@render children()}
</div>

<!-- Toast notifications -->
<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
	{#each toasts as t (t.id)}
		<div
			class="pointer-events-auto flex items-center gap-3 px-4 py-3 rounded-lg border {getToastClasses(t.type)}"
		>
			<svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getToastIcon(t.type)} />
			</svg>
			<p class="text-sm text-text-primary">{t.message}</p>
		</div>
	{/each}
</div>
