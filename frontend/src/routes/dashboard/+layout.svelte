<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let loading = $state(true);
	let showUserMenu = $state(false);

	onMount(() => {
		const unsubscribe = auth.subscribe((state) => {
			loading = state.loading;
			if (!state.loading && !state.user) {
				goto('/');
			}
		});
		return unsubscribe;
	});

	async function handleLogout() {
		showUserMenu = false;
		await auth.logout();
		goto('/');
	}

	function toggleUserMenu(event: MouseEvent) {
		event.stopPropagation();
		showUserMenu = !showUserMenu;
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

{#if loading}
	<div class="min-h-screen flex items-center justify-center">
		<div class="spinner spinner-lg"></div>
	</div>
{:else if $auth.user}
	<div class="min-h-screen bg-surface-900">
		<!-- Header -->
		<header class="sticky top-0 z-50 bg-surface-800 border-b border-surface-600">
			<div class="max-w-6xl mx-auto px-4 lg:px-6">
				<div class="flex items-center justify-between h-14 sm:h-16">
					<!-- Logo -->
					<a href="/dashboard" class="flex items-center gap-2.5">
						<img src="/icon.png" alt="Prize Bot" class="w-8 h-8 rounded-lg" />
						<span class="font-semibold text-text-primary hidden sm:block">Prize Bot</span>
					</a>

					<!-- Right side -->
					<div class="flex items-center gap-4">
						<!-- Support link -->
						<a
							href="https://discord.gg/RxP2z5NGtj"
							target="_blank"
							rel="noopener"
							class="text-sm text-text-secondary hover:text-text-primary transition-colors hidden md:block"
						>
							Support
						</a>

						<!-- User menu -->
						<div class="relative user-menu-container">
							<button
								onclick={toggleUserMenu}
								class="flex items-center gap-2.5 px-2 py-1.5 rounded-lg hover:bg-surface-700 transition-colors"
							>
								{#if $auth.user.avatar_url}
									<img
										src={$auth.user.avatar_url}
										alt={$auth.user.username}
										class="w-7 h-7 rounded-full"
									/>
								{:else}
									<div class="w-7 h-7 rounded-full bg-accent flex items-center justify-center text-xs font-medium text-white">
										{$auth.user.username[0].toUpperCase()}
									</div>
								{/if}
								<span class="text-sm text-text-primary hidden sm:block">{$auth.user.username}</span>
								<svg class="w-4 h-4 text-text-muted {showUserMenu ? 'rotate-180' : ''} transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
								</svg>
							</button>

							<!-- Dropdown -->
							{#if showUserMenu}
								<div class="absolute right-0 mt-1 w-48 bg-surface-700 border border-surface-600 rounded-lg shadow-lg overflow-hidden">
									<div class="px-3 py-2 border-b border-surface-600">
										<p class="text-sm font-medium text-text-primary">{$auth.user.username}</p>
										<p class="text-xs text-text-muted">Discord</p>
									</div>
									<div class="p-1">
										<a href="/dashboard" class="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-600 rounded-md transition-colors">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
											</svg>
											My Servers
										</a>
										<a href="/dashboard/subscription" class="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-600 rounded-md transition-colors">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
											</svg>
											Subscription
										</a>
										<a href="/dashboard/whitelabels" class="flex items-center gap-2 px-3 py-2 text-sm text-text-secondary hover:text-text-primary hover:bg-surface-600 rounded-md transition-colors">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
											</svg>
											Whitelabels
										</a>
										<button
											onclick={handleLogout}
											class="w-full flex items-center gap-2 px-3 py-2 text-sm text-status-danger hover:bg-status-danger/10 rounded-md transition-colors"
										>
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
											</svg>
											Sign out
										</button>
									</div>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</header>

		<!-- Main content -->
		<main class="max-w-6xl mx-auto px-4 lg:px-6 py-4 lg:py-6">
			{@render children()}
		</main>
	</div>
{/if}
