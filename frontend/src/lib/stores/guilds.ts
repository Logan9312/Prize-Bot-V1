import { writable, derived } from 'svelte/store';
import type { Guild } from '$lib/api/client';

function createGuildsStore() {
	const { subscribe, set } = writable<Guild[]>([]);

	return {
		subscribe,
		set,
		clear() {
			set([]);
		}
	};
}

export const guilds = createGuildsStore();

export function getGuildById(guildId: string) {
	return derived(guilds, ($guilds) => $guilds.find((g) => g.id === guildId) || null);
}
