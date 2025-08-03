import { apiClient } from '$lib/apiClient.js';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for the entire app

export const load: LayoutLoad = async () => {
  try {
    const user = await apiClient.users.getProfile();
    return { user };
  } catch (error) {
    if (!(error instanceof Error && 'status' in error && error.status === 401)) {
      console.error('Failed to fetch user profile:', error);
    }
    return { user: null };
  }
};