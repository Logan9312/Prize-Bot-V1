<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';

	onMount(async () => {
		const error = $page.url.searchParams.get('error');
		if (error) {
			toast.error('Login failed: ' + error);
			goto('/');
			return;
		}

		await auth.init();
		goto('/dashboard');
	});
</script>

<div class="min-h-screen flex items-center justify-center">
	<div class="text-center">
		<div
			class="w-12 h-12 border-2 border-discord-blurple border-t-transparent rounded-full animate-spin mx-auto mb-4"
		></div>
		<p class="text-gray-400">Completing login...</p>
	</div>
</div>
