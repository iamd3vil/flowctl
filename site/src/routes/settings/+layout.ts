import { error, redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent, url }) => {
  const { userPromise } = await parent();
  const user = await userPromise;

  // Redirect to login if not authenticated
  if (!user) {
    const redirectUrl = url.pathname + url.search;
    throw redirect(302, `/login?redirect_url=${encodeURIComponent(redirectUrl)}`);
  }

  // Settings page requires superuser role
  if (user.role !== 'superuser') {
    throw error(403, 'Access denied. Superuser privileges required.');
  }

  return { user };
};
