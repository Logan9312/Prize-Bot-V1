<script lang="ts">
	import { page } from '$app/stores';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let mobileMenuOpen = $state(false);

	const navItems = [
		{ href: '/docs', label: 'Overview', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
		{ href: '/docs/auctions', label: 'Auctions', icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
		{ href: '/docs/giveaways', label: 'Giveaways', icon: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7' },
		{ href: '/docs/claims', label: 'Claims', icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' },
		{ href: '/docs/currency', label: 'Currency', icon: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z' },
		{ href: '/docs/premium', label: 'Premium', icon: 'M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z' },
		{ href: '/docs/whitelabel', label: 'Whitelabel', icon: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01' }
	];

	function isActive(href: string, currentPath: string): boolean {
		if (href === '/docs') {
			return currentPath === '/docs';
		}
		return currentPath.startsWith(href);
	}
</script>

<div class="min-h-screen">
	<!-- Header -->
	<header class="sticky top-0 z-40 bg-surface-900/95 backdrop-blur border-b border-surface-600">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<!-- Logo -->
				<a href="/" class="flex items-center gap-3">
					<img src="/icon.png" alt="Prize Bot" class="w-8 h-8 rounded-lg" />
					<span class="font-semibold text-text-primary">Prize Bot</span>
				</a>

				<!-- Desktop Navigation -->
				<nav class="hidden md:flex items-center gap-6">
					<a href="/docs" class="text-sm text-text-secondary hover:text-text-primary transition-colors">
						Docs
					</a>
					<a href="/pricing" class="text-sm text-text-secondary hover:text-text-primary transition-colors">
						Pricing
					</a>
					<a href="/" class="btn btn-primary text-sm !py-2 !px-4 !min-h-0">
						Get Started
					</a>
				</nav>

				<!-- Mobile menu button -->
				<button
					onclick={() => mobileMenuOpen = !mobileMenuOpen}
					class="md:hidden p-2 text-text-secondary hover:text-text-primary"
					aria-label="Toggle menu"
				>
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						{#if mobileMenuOpen}
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						{:else}
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
						{/if}
					</svg>
				</button>
			</div>
		</div>

		<!-- Mobile menu -->
		{#if mobileMenuOpen}
			<div class="md:hidden border-t border-surface-600 bg-surface-800">
				<nav class="px-4 py-4 space-y-2">
					{#each navItems as item}
						<a
							href={item.href}
							onclick={() => mobileMenuOpen = false}
							class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm {isActive(item.href, $page.url.pathname) ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-700 hover:text-text-primary'}"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
							</svg>
							{item.label}
						</a>
					{/each}
					<div class="pt-2 border-t border-surface-600">
						<a href="/pricing" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-text-secondary hover:bg-surface-700 hover:text-text-primary">
							Pricing
						</a>
						<a href="/" class="flex items-center justify-center gap-2 mt-2 btn btn-primary w-full">
							Get Started
						</a>
					</div>
				</nav>
			</div>
		{/if}
	</header>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<div class="flex gap-8">
			<!-- Desktop Sidebar -->
			<aside class="hidden md:block w-64 flex-shrink-0">
				<nav class="sticky top-24 space-y-1">
					{#each navItems as item}
						<a
							href={item.href}
							class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors {isActive(item.href, $page.url.pathname) ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-700 hover:text-text-primary'}"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
							</svg>
							{item.label}
						</a>
					{/each}
				</nav>
			</aside>

			<!-- Main content -->
			<main class="flex-1 min-w-0">
				{@render children()}
			</main>
		</div>
	</div>

	<!-- Footer -->
	<footer class="border-t border-surface-600 mt-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="flex flex-col sm:flex-row items-center justify-between gap-4">
				<div class="flex items-center gap-3">
					<img src="/icon.png" alt="Prize Bot" class="w-6 h-6 rounded" />
					<span class="text-sm text-text-muted">Prize Bot</span>
				</div>
				<div class="flex items-center gap-6 text-sm text-text-muted">
					<a href="/docs" class="hover:text-text-primary transition-colors">Docs</a>
					<a href="/pricing" class="hover:text-text-primary transition-colors">Pricing</a>
					<a href="https://discord.gg/RxP2z5NGtj" target="_blank" rel="noopener" class="hover:text-text-primary transition-colors">Support</a>
				</div>
			</div>
		</div>
	</footer>
</div>
