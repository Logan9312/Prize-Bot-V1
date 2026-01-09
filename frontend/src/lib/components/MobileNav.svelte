<script lang="ts">
	interface NavItem {
		href: string;
		label: string;
		icon: string;
	}

	let {
		navItems,
		currentPath,
		guildId
	}: {
		navItems: NavItem[];
		currentPath: string;
		guildId: string;
	} = $props();

	let isOpen = $state(false);

	function toggleDrawer() {
		isOpen = !isOpen;
	}

	function closeDrawer() {
		isOpen = false;
	}

	// Close drawer when route changes
	$effect(() => {
		currentPath;
		isOpen = false;
	});
</script>

<!-- Hamburger Button -->
<button
	onclick={toggleDrawer}
	class="fixed top-3.5 left-4 z-50 flex items-center justify-center w-12 h-12 rounded-lg bg-surface-800 border border-surface-600 hover:bg-surface-700 transition-colors duration-200"
	aria-label={isOpen ? 'Close menu' : 'Open menu'}
	aria-expanded={isOpen}
>
	<div class="relative w-6 h-6 flex items-center justify-center">
		<!-- Hamburger icon with animated lines -->
		<span class="absolute inset-0 flex flex-col items-center justify-center gap-1.5">
			<span
				class="block w-5 h-0.5 bg-text-primary rounded-full transition-all duration-300 ease-out origin-center"
				class:rotate-45={isOpen}
				class:translate-y-[5px]={isOpen}
			></span>
			<span
				class="block w-5 h-0.5 bg-text-primary rounded-full transition-all duration-300 ease-out"
				class:opacity-0={isOpen}
				class:scale-x-0={isOpen}
			></span>
			<span
				class="block w-5 h-0.5 bg-text-primary rounded-full transition-all duration-300 ease-out origin-center"
				class:-rotate-45={isOpen}
				class:-translate-y-[5px]={isOpen}
			></span>
		</span>
	</div>
</button>

<!-- Backdrop -->
{#if isOpen}
	<button
		onclick={closeDrawer}
		class="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm animate-fade-in"
		aria-label="Close menu"
		tabindex="-1"
	></button>
{/if}

<!-- Drawer -->
<nav
	class="fixed top-0 left-0 bottom-0 z-40 w-[280px] max-w-[85vw] bg-surface-800 border-r border-surface-600 shadow-2xl transition-transform duration-300 ease-out"
	class:-translate-x-full={!isOpen}
	class:translate-x-0={isOpen}
	aria-hidden={!isOpen}
>
	<div class="flex flex-col h-full pt-20 pb-6 px-4">
		<!-- Back to servers link -->
		<a
			href="/dashboard"
			class="flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary transition-colors mb-6 py-2"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
			Back to servers
		</a>

		<!-- Navigation items -->
		<div class="flex-1 space-y-1 overflow-y-auto">
			{#each navItems as item, index}
				<a
					href={item.href}
					class="flex items-center gap-3 px-4 py-3 rounded-lg font-medium text-sm transition-all duration-200 min-h-[48px]"
					class:bg-accent={currentPath === item.href}
					class:text-white={currentPath === item.href}
					class:hover:bg-accent={currentPath === item.href}
					class:text-text-secondary={currentPath !== item.href}
					class:hover:bg-surface-700={currentPath !== item.href}
					class:hover:text-text-primary={currentPath !== item.href}
					style="animation-delay: {isOpen ? index * 40 : 0}ms"
					class:slide-in={isOpen}
				>
					<svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
					</svg>
					<span>{item.label}</span>
				</a>
			{/each}
		</div>
	</div>
</nav>

<style>
	@keyframes slide-in {
		from {
			opacity: 0;
			transform: translateX(-12px);
		}
		to {
			opacity: 1;
			transform: translateX(0);
		}
	}

	.slide-in {
		animation: slide-in 0.3s ease-out forwards;
	}

	/* Custom scrollbar for nav items */
	nav > div {
		scrollbar-width: thin;
		scrollbar-color: #3f4147 transparent;
	}

	nav > div::-webkit-scrollbar {
		width: 6px;
	}

	nav > div::-webkit-scrollbar-track {
		background: transparent;
	}

	nav > div::-webkit-scrollbar-thumb {
		background-color: #3f4147;
		border-radius: 3px;
	}

	nav > div::-webkit-scrollbar-thumb:hover {
		background-color: #4e5058;
	}
</style>
