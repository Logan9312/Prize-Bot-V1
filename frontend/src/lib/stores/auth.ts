import { writable } from 'svelte/store';
import { authAPI, type User } from '$lib/api/client';

interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: true,
		error: null
	});

	return {
		subscribe,

		async init() {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const user = await authAPI.getCurrentUser();
				set({ user, loading: false, error: null });
			} catch {
				set({ user: null, loading: false, error: null });
			}
		},

		async logout() {
			try {
				await authAPI.logout();
			} catch {
				// Ignore logout errors
			}
			set({ user: null, loading: false, error: null });
		},

		clear() {
			set({ user: null, loading: false, error: null });
		}
	};
}

export const auth = createAuthStore();
