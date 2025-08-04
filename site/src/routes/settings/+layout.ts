import { error } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
  const { user } = await parent();
  
  // Settings page requires superuser role
  if (user?.role != 'superuser') {
    error(403, 'Access denied. Superuser privileges required.');
  }
  
  return {};
};