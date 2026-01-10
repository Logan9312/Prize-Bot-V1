<script lang="ts">
	interface Props {
		feature: 'auctions' | 'giveaways' | 'claims';
		activeTab: 'settings' | 'list';
		guildId: string;
		listCount?: number;
	}

	let { feature, activeTab, guildId, listCount }: Props = $props();

	const baseUrl = $derived(`/dashboard/${guildId}/${feature}`);
</script>

<div class="flex items-center gap-1 p-1 bg-surface-900 rounded-lg w-fit">
	<a
		href={baseUrl}
		class="px-4 py-2 text-sm font-medium rounded-md transition-all {activeTab === 'settings'
			? 'bg-surface-700 text-text-primary shadow-sm'
			: 'text-text-secondary hover:text-text-primary hover:bg-surface-800'}"
	>
		Settings
	</a>
	<a
		href="{baseUrl}/list"
		class="px-4 py-2 text-sm font-medium rounded-md transition-all flex items-center gap-2 {activeTab === 'list'
			? 'bg-surface-700 text-text-primary shadow-sm'
			: 'text-text-secondary hover:text-text-primary hover:bg-surface-800'}"
	>
		List
		{#if listCount !== undefined && listCount > 0}
			<span class="px-1.5 py-0.5 text-xs font-medium rounded-full bg-accent text-white min-w-[1.25rem] text-center">
				{listCount > 99 ? '99+' : listCount}
			</span>
		{/if}
	</a>
</div>
