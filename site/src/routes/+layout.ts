import { apiClient } from '$lib/apiClient.js';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for the entire app

export const load: LayoutLoad = () => {
  // Return promise without awaiting - layout renders immediately while user data loads
  const userPromise = apiClient.users.getProfile().catch((error) => {
    console.error('[Auth] Failed to fetch user profile:', error);
    return null;
  });

  return { userPromise };
};