import { writable } from 'svelte/store';
import type { UserProfileResponse } from '$lib/types.js';

export const currentUser = writable<UserProfileResponse | null>(null);
export const isAuthenticated = writable<boolean>(false);
export const isLoading = writable<boolean>(true);