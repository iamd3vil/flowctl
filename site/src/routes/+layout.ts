import { apiClient } from '$lib/apiClient.js';
import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const ssr = false; // Disable SSR for the entire app

export const load: LayoutLoad = async ({ url }) => {
  // Allow access to login page and root
  if (url.pathname === '/login' || url.pathname === '/') {
    try {
      const user = await apiClient.users.getProfile();
      return { user };
    } catch (error) {
      return { user: null };
    }
  }

  // For all other routes, require authentication
  try {
    const user = await apiClient.users.getProfile();
    return { user };
  } catch (error) {
    if (error instanceof Error && 'status' in error && error.status === 401) {
      throw redirect(302, '/login');
    }
    console.error('Failed to fetch user profile:', error);
    throw redirect(302, '/login');
  }
};