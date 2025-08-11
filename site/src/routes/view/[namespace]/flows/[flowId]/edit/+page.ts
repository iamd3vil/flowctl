import { apiClient } from '$lib/apiClient.js';
import { permissionChecker } from '$lib/utils/permissions';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent }) => {
  const { user, namespaceId } = await parent();
  
  // Check update permissions
  try {
    const permissions = await permissionChecker(user!, 'flow', namespaceId, ['update']);
    if (!permissions.canUpdate) {
      error(403, {
        message: 'You do not have permission to edit flows in this namespace',
        code: 'INSUFFICIENT_PERMISSIONS'
      });
    }
  } catch (err) {
    if (err && typeof err === 'object' && 'status' in err) {
      throw err; // Re-throw SvelteKit errors
    }
    error(500, {
      message: 'Failed to check permissions',
      code: 'PERMISSION_CHECK_FAILED'
    });
  }
  
  try {
    const data = await apiClient.executors.list();
    
    const availableExecutors = data.executors.map(name => ({
      name: name,
      display_name: name.charAt(0).toUpperCase() + name.slice(1)
    }));

    return {
      availableExecutors
    };
  } catch (loadError) {
    console.error('Error loading executors:', loadError);
    return {
      availableExecutors: []
    };
  }
};