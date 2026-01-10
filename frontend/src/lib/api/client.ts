const API_BASE = import.meta.env.VITE_API_URL || '/api';

interface FetchOptions extends RequestInit {
	body?: any;
}

async function fetchAPI<T>(endpoint: string, options: FetchOptions = {}): Promise<T> {
	const { body, ...rest } = options;

	const config: RequestInit = {
		...rest,
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...rest.headers
		}
	};

	if (body) {
		config.body = JSON.stringify(body);
	}

	const response = await fetch(`${API_BASE}${endpoint}`, config);

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: 'Request failed' }));
		throw new Error(error.error || 'Request failed');
	}

	return response.json();
}

export const api = {
	get: <T>(endpoint: string) => fetchAPI<T>(endpoint, { method: 'GET' }),
	post: <T>(endpoint: string, body?: any) => fetchAPI<T>(endpoint, { method: 'POST', body }),
	put: <T>(endpoint: string, body?: any) => fetchAPI<T>(endpoint, { method: 'PUT', body }),
	delete: <T>(endpoint: string) => fetchAPI<T>(endpoint, { method: 'DELETE' })
};

// Auth API
export const authAPI = {
	getCurrentUser: () => api.get<User>('/auth/me'),
	logout: () => api.post('/auth/logout')
};

// Guilds API
export const guildsAPI = {
	list: () => api.get<{ guilds: Guild[] }>('/guilds'),
	getChannels: (guildId: string) => api.get<{ channels: Channel[] }>(`/guilds/${guildId}/channels`),
	getRoles: (guildId: string) => api.get<{ roles: Role[] }>(`/guilds/${guildId}/roles`),
	getStats: (guildId: string) => api.get<GuildStats>(`/guilds/${guildId}/stats`),
	listAuctions: (guildId: string) => api.get<{ auctions: AuctionListItem[] }>(`/guilds/${guildId}/auctions/list`),
	listGiveaways: (guildId: string) => api.get<{ giveaways: GiveawayListItem[] }>(`/guilds/${guildId}/giveaways/list`),
	listClaims: (guildId: string) => api.get<{ claims: ClaimListItem[] }>(`/guilds/${guildId}/claims/list`)
};

// Claims API
export const claimsAPI = {
	update: (guildId: string, messageId: string, data: { item?: string; winner?: string; cost?: number }) =>
		api.put(`/guilds/${guildId}/claims/${messageId}`, data),
	resend: (guildId: string, messageId: string) =>
		api.post(`/guilds/${guildId}/claims/${messageId}/resend`),
	cancel: (guildId: string, messageId: string) =>
		api.delete(`/guilds/${guildId}/claims/${messageId}`)
};

// Premium API
export const premiumAPI = {
	getUserStatus: () => api.get<UserPremiumStatus>('/premium/status'),
	getGuildStatus: (guildId: string) => api.get<GuildPremiumStatus>(`/guilds/${guildId}/premium`),
	createPortalSession: () => api.post<BillingPortalResponse>('/premium/portal')
};

// Whitelabel API
export const whitelabelAPI = {
	list: () => api.get<WhitelabelListResponse>('/whitelabels'),
	create: (botToken: string) =>
		api.post<WhitelabelCreateResponse>('/whitelabels', { bot_token: botToken }),
	delete: (botId: string) => api.delete<{ message: string }>(`/whitelabels/${botId}`),
	validate: (botToken: string) =>
		api.post<ValidateTokenResponse>('/whitelabels/validate', { bot_token: botToken })
};

// Settings API
export const settingsAPI = {
	// Auction
	getAuction: (guildId: string) => api.get<AuctionSettings>(`/guilds/${guildId}/settings/auction`),
	updateAuction: (guildId: string, data: Partial<AuctionSettings>) =>
		api.put(`/guilds/${guildId}/settings/auction`, data),
	deleteAuction: (guildId: string) => api.delete(`/guilds/${guildId}/settings/auction`),

	// Claim
	getClaim: (guildId: string) => api.get<ClaimSettings>(`/guilds/${guildId}/settings/claim`),
	updateClaim: (guildId: string, data: Partial<ClaimSettings>) =>
		api.put(`/guilds/${guildId}/settings/claim`, data),
	deleteClaim: (guildId: string) => api.delete(`/guilds/${guildId}/settings/claim`),

	// Giveaway
	getGiveaway: (guildId: string) => api.get<GiveawaySettings>(`/guilds/${guildId}/settings/giveaway`),
	updateGiveaway: (guildId: string, data: Partial<GiveawaySettings>) =>
		api.put(`/guilds/${guildId}/settings/giveaway`, data),
	deleteGiveaway: (guildId: string) => api.delete(`/guilds/${guildId}/settings/giveaway`),

	// Currency
	getCurrency: (guildId: string) => api.get<CurrencySettings>(`/guilds/${guildId}/settings/currency`),
	updateCurrency: (guildId: string, data: Partial<CurrencySettings>) =>
		api.put(`/guilds/${guildId}/settings/currency`, data),
	deleteCurrency: (guildId: string) => api.delete(`/guilds/${guildId}/settings/currency`),

	// Shop
	getShop: (guildId: string) => api.get<ShopSettings>(`/guilds/${guildId}/settings/shop`),
	updateShop: (guildId: string, data: Partial<ShopSettings>) =>
		api.put(`/guilds/${guildId}/settings/shop`, data),
	deleteShop: (guildId: string) => api.delete(`/guilds/${guildId}/settings/shop`)
};

// Types
export interface User {
	id: string;
	username: string;
	avatar: string;
	avatar_url: string;
}

export interface Guild {
	id: string;
	name: string;
	icon: string;
	icon_url: string;
	owner: boolean;
	is_admin: boolean;
	bot_in: boolean;
}

export interface Channel {
	id: string;
	name: string;
	type: number;
	position: number;
	parent_id: string;
}

export interface Role {
	id: string;
	name: string;
	color: number;
	position: number;
	managed: boolean;
}

export interface AuctionSettings {
	guild_id: string;
	category?: string;
	alert_role?: string;
	currency?: string;
	log_channel?: string;
	host_role?: string;
	snipe_extension?: number;
	snipe_range?: number;
	snipe_limit?: number;
	snipe_cap?: number;
	currency_side?: string;
	integer_only?: boolean;
	channel_override?: string;
	channel_lock?: boolean;
	channel_prefix?: string;
	use_currency?: boolean;
}

export interface ClaimSettings {
	guild_id: string;
	category?: string;
	staff_role?: string;
	instructions?: string;
	log_channel?: string;
	expiration?: string;
	disable_claiming?: boolean;
	channel_prefix?: string;
}

export interface GiveawaySettings {
	guild_id: string;
	host_role?: string;
	alert_role?: string;
	log_channel?: string;
}

export interface CurrencySettings {
	guild_id: string;
	currency?: string;
	side?: string;
}

export interface ShopSettings {
	guild_id: string;
	host_role?: string;
	alert_role?: string;
	log_channel?: string;
}

export interface GuildStats {
	active_auctions: number;
	running_giveaways: number;
	open_claims: number;
	shop_items: number;
}

export interface AuctionListItem {
	channel_id: string;
	item: string;
	bid: number;
	winner: string;
	host: string;
	end_time: string;
	currency: string;
	currency_side: string;
}

export interface GiveawayListItem {
	message_id: string;
	channel_id: string;
	item: string;
	end_time: string;
	host: string;
	winners: number;
}

export interface ClaimListItem {
	message_id: string;
	channel_id: string;
	item: string;
	winner: string;
	cost: number;
	status: string;
	ticket_id: string;
}

// Premium types
export interface SubscriptionInfo {
	id: string;
	status: string;
	current_period_end: number;
	guild_id?: string;
	plan_name: string;
}

export interface UserPremiumStatus {
	is_premium: boolean;
	subscriptions: SubscriptionInfo[];
}

export interface GuildPremiumStatus {
	is_premium: boolean;
	guild_id: string;
}

export interface BillingPortalResponse {
	url: string;
}

// Whitelabel types
export interface Whitelabel {
	bot_id: string;
	bot_name: string;
	bot_avatar: string;
	user_id: string;
	created_at: string;
}

export interface WhitelabelListResponse {
	whitelabels: Whitelabel[];
	is_admin: boolean;
}

export interface WhitelabelCreateResponse {
	message: string;
	bot_id: string;
	bot_name: string;
	bot_avatar: string;
}

export interface ValidateTokenResponse {
	valid: boolean;
	bot_id?: string;
	bot_name?: string;
	bot_avatar?: string;
	error?: string;
}
